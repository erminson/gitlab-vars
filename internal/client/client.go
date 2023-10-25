package varsapi

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/erminson/gitlab-vars/internal/types"
	"io"
	"net/http"
	"net/url"
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

func (v *VarsAPI) MakeRequestWithContext(ctx context.Context, method string, endpoint string, filter types.Filter, varData types.VarData) (*types.APIResponse, error) {
	uri := fmt.Sprintf(v.apiEndpoint, endpoint)

	values := buildBody(varData)

	req, err := http.NewRequestWithContext(ctx, method, uri, strings.NewReader(values.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("PRIVATE-TOKEN", v.Token)

	q := req.URL.Query()
	queryValues := buildRawQueryValues(q, filter)
	req.URL.RawQuery = queryValues

	resp, err := v.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp types.APIResponse
	_, err = decodeAPIResponse(resp.Body, &apiResp)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 && resp.StatusCode <= 499 {
		var apiErr types.APIError
		_, err := parseAPIError(apiResp.Result, &apiErr)
		if err != nil {
			return nil, err
		}
		apiErr.Code = resp.StatusCode

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

func decodeAPIResponse(responseBody io.Reader, resp *types.APIResponse) ([]byte, error) {
	data, err := io.ReadAll(responseBody)
	if err != nil {
		return nil, err
	}

	resp.Result = data

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

	err := params.ValidateProjectId()
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf(APIEndpointVars+"?per_page=100", params.ProjectId, params.Key)
	resp, err := v.MakeRequestWithContext(ctx, "GET", endpoint, types.Filter{}, types.VarData{})
	if err != nil {
		return nil, err
	}

	variables := make([]types.Variable, 0)
	err = json.Unmarshal(resp.Result, &variables)
	if err != nil {
		return nil, err
	}

	return variables, nil
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
