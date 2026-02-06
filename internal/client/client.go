package varsapi

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/erminson/gitlab-vars/internal/types"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const APIHost = "https://gitlab.com"
const APIEndpoint = "/api/v4/%s"
const APIEndpointVars = "projects/%d/variables/%s"
const APIEndpointPersonalTokens = "personal_access_tokens/self"

const timeout = time.Second * 3

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type VarsAPI struct {
	Token       string
	Client      HTTPClient
	apiEndpoint string
	tokenInfo   types.Token
}

func NewVars(token string) (*VarsAPI, error) {
	return NewVarsAPIWithClient(token, APIHost, &http.Client{})
}

func NewVarsWithHost(token, host string) (*VarsAPI, error) {
	return NewVarsAPIWithClient(token, host, &http.Client{})
}

func NewVarsAPIWithClient(token, host string, client HTTPClient) (*VarsAPI, error) {
	apiEndpoint := fmt.Sprintf("%s%s", host, APIEndpoint)
	varsAPI := &VarsAPI{
		Token:       token,
		Client:      client,
		apiEndpoint: apiEndpoint,
	}

	// TODO: Validate token
	tokenInfo, err := varsAPI.GetSelfToken()
	if err != nil {
		return nil, err
	}

	varsAPI.tokenInfo = tokenInfo

	return varsAPI, nil
}

func (v *VarsAPI) MakeRequest(method string, endpoint string, filter types.Filter, varData types.VarData) (*types.APIResponse, error) {
	ctx := context.Background()
	return v.MakeRequestWithContext(ctx, method, endpoint, filter, varData)
}

func (v *VarsAPI) MakeRequestWithContext(
	ctx context.Context,
	method string,
	endpoint string,
	filter types.Filter,
	varData types.VarData,
) (*types.APIResponse, error) {

	uri := fmt.Sprintf(v.apiEndpoint, endpoint)

	values := buildBody(varData)

	req, err := http.NewRequestWithContext(
		ctx,
		method,
		uri,
		strings.NewReader(values.Encode()),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("PRIVATE-TOKEN", v.Token)

	q := req.URL.Query()
	req.URL.RawQuery = buildRawQueryValues(q, filter)

	httpResp, err := v.Client.Do(req)
	if err != nil {
		return nil, err
	}

	var apiResp types.APIResponse
	_, err = decodeAPIResponse(httpResp, &apiResp)
	if err != nil {
		return nil, err
	}

	// Handle client-side errors (GitLab-style)
	if apiResp.Status >= 400 && apiResp.Status <= 499 {
		var apiErr types.APIError

		if _, err := parseAPIError(apiResp.Result, &apiErr); err != nil {
			return nil, err
		}

		apiErr.Code = apiResp.Status
		return nil, apiErr
	}

	return &apiResp, nil
}

func buildBody(in types.VarData) url.Values {
	if in == nil {
		return url.Values{}
	}

	out := url.Values{}
	for k, v := range in {
		out.Set(k, v)
	}

	return out
}

func buildRawQueryValues(in url.Values, filter types.Filter) string {
	if in == nil {
		return url.Values{}.Encode()
	}

	out := in
	for k, v := range filter {
		if k != "" && v != "" {
			in.Set(fmt.Sprintf("filter[%s]", k), v)
		}
	}

	return out.Encode()
}

func decodeAPIResponse(httpResp *http.Response, resp *types.APIResponse) ([]byte, error) {
	defer httpResp.Body.Close()

	data, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}

	resp.Result = data
	resp.Headers = httpResp.Header
	resp.Status = httpResp.StatusCode

	return data, nil
}

func parseAPIError(data []byte, errResp *types.APIError) ([]byte, error) {
	err := json.Unmarshal(data, errResp)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (v *VarsAPI) GetSelfToken() (types.Token, error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	resp, err := v.MakeRequestWithContext(ctx, "GET", APIEndpointPersonalTokens, types.Filter{}, types.VarData{})
	if err != nil {
		return types.Token{}, err
	}

	var token types.Token
	err = json.Unmarshal(resp.Result, &token)

	return token, err
}

func (v *VarsAPI) GetVariables(params types.Params) ([]types.Variable, error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err := params.ValidateProjectId(); err != nil {
		return nil, err
	}

	allVariables := make([]types.Variable, 0)
	page := 1

	for {
		endpoint := fmt.Sprintf(
			APIEndpointVars+"?per_page=100&page=%s",
			params.ProjectId,
			params.Key,
			strconv.Itoa(page),
		)

		resp, err := v.MakeRequestWithContext(
			ctx,
			"GET",
			endpoint,
			types.Filter{},
			types.VarData{},
		)
		if err != nil {
			return nil, err
		}

		var varsPage []types.Variable
		if err := json.Unmarshal(resp.Result, &varsPage); err != nil {
			return nil, err
		}

		allVariables = append(allVariables, varsPage...)

		// GitLab pagination logic
		nextPage := resp.Headers.Get("X-Next-Page")
		if nextPage == "" {
			break
		}

		page, err = strconv.Atoi(nextPage)
		if err != nil {
			return nil, fmt.Errorf("invalid X-Next-Page header: %w", err)
		}
	}

	return allVariables, nil
}

func (v *VarsAPI) GetVariable(params types.Params, filter types.Filter) (types.Variable, error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	err := params.Validate()
	if err != nil {
		return types.Variable{}, err
	}
	endpoint := fmt.Sprintf(APIEndpointVars, params.ProjectId, params.Key)
	resp, err := v.MakeRequestWithContext(ctx, "GET", endpoint, filter, types.VarData{})
	if err != nil {
		return types.Variable{}, nil
	}

	var variable types.Variable
	err = json.Unmarshal(resp.Result, &variable)
	if err != nil {
		return types.Variable{}, err
	}

	return variable, nil
}

func (v *VarsAPI) CreateVariableFromVarData(params types.Params, data types.VarData) (types.Variable, error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	err := params.ValidateProjectId()
	if err != nil {
		return types.Variable{}, err
	}

	err = data.Validate()
	if err != nil {
		return types.Variable{}, err
	}

	endpoint := fmt.Sprintf(APIEndpointVars, params.ProjectId, "")
	resp, err := v.MakeRequestWithContext(ctx, "POST", endpoint, types.Filter{}, data)
	if err != nil {
		return types.Variable{}, err
	}

	var variable types.Variable
	err = json.Unmarshal(resp.Result, &variable)
	if err != nil {
		return types.Variable{}, err
	}

	return variable, nil
}

func (v *VarsAPI) CreateVariable(params types.Params, variable types.Variable) (types.Variable, error) {
	return v.CreateVariableFromVarData(params, variable.VariableToData())
}

func (v *VarsAPI) UpdateVariableFromVarData(params types.Params, data types.VarData, filter types.Filter) (types.Variable, error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	err := params.Validate()
	if err != nil {
		return types.Variable{}, err
	}

	err = data.Validate()
	if err != nil {
		return types.Variable{}, err
	}

	endpoint := fmt.Sprintf(APIEndpointVars, params.ProjectId, params.Key)
	resp, err := v.MakeRequestWithContext(ctx, "PUT", endpoint, filter, data)
	if err != nil {
		return types.Variable{}, err
	}

	var variable types.Variable
	err = json.Unmarshal(resp.Result, &variable)
	if err != nil {
		return types.Variable{}, err
	}

	return variable, nil
}

func (v *VarsAPI) UpdateVariable(params types.Params, variable types.Variable, filter types.Filter) (types.Variable, error) {
	return v.UpdateVariableFromVarData(params, variable.VariableToData(), filter)
}

func (v *VarsAPI) DeleteVariable(params types.Params, filter types.Filter) error {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	err := params.Validate()
	if err != nil {
		return err
	}

	endpoint := fmt.Sprintf(APIEndpointVars, params.ProjectId, params.Key)
	_, err = v.MakeRequestWithContext(ctx, "DELETE", endpoint, filter, types.VarData{})

	return err
}
