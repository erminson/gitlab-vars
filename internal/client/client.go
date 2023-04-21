package varsapi

import (
	"encoding/json"
	"fmt"
	"github.com/erminson/gitlab-vars/internal/types"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const APIEndpoint = "https://gitlab.com/api/v4/%s"
const APIEndpointVars = "projects/%d/variables/%s"
const APIEndpointPersonalTokens = "personal_access_tokens/self"

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type VarsAPI struct {
	Token           string
	Client          HTTPClient
	apiEndpoint     string
	apiEndpointVars string
	tokenInfo       types.Token
}

func NewVars(token string) (*VarsAPI, error) {
	return NewVarsAPIWithClient(token, APIEndpoint, &http.Client{})
}

func NewVarsAPIWithClient(token, apiEndpoint string, client HTTPClient) (*VarsAPI, error) {
	varsAPI := &VarsAPI{
		Token:           token,
		Client:          client,
		apiEndpoint:     apiEndpoint,
		apiEndpointVars: APIEndpointVars,
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
	uri := fmt.Sprintf(v.apiEndpoint, endpoint)

	values := buildBody(varData)

	req, err := http.NewRequest(method, uri, strings.NewReader(values.Encode()))
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
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	var apiResp types.APIResponse
	_, err = decodeAPIResponse(resp.Body, &apiResp)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 && resp.StatusCode <= 499 {
		var apiErr types.APIError
		_, err = parseAPIError(apiResp.Result, &apiErr)
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
		in.Set(fmt.Sprintf("filter[%s]", k), v)
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
	resp, err := v.MakeRequest("GET", APIEndpointPersonalTokens, types.Filter{}, types.VarData{})
	if err != nil {
		return types.Token{}, err
	}

	var token types.Token
	err = json.Unmarshal(resp.Result, &token)

	return token, err
}

func (v *VarsAPI) GetVariables(params types.Params) ([]types.Variable, error) {
	err := params.ValidateProjectId()
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf(APIEndpointVars, params.ProjectId, params.Key)
	resp, err := v.MakeRequest("GET", endpoint, types.Filter{}, types.VarData{})
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
	err := params.Validate()
	if err != nil {
		return types.Variable{}, err
	}
	endpoint := fmt.Sprintf(APIEndpointVars, params.ProjectId, params.Key)
	resp, err := v.MakeRequest("GET", endpoint, filter, types.VarData{})
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
	err := params.ValidateProjectId()
	if err != nil {
		return types.Variable{}, err
	}

	err = data.Validate()
	if err != nil {
		return types.Variable{}, err
	}

	endpoint := fmt.Sprintf(APIEndpointVars, params.ProjectId, "")
	resp, err := v.MakeRequest("POST", endpoint, types.Filter{}, data)
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

func (v *VarsAPI) UpdateVariable(params types.Params, data types.VarData, filter types.Filter) (types.Variable, error) {
	err := params.Validate()
	if err != nil {
		return types.Variable{}, err
	}

	err = data.Validate()
	if err != nil {
		return types.Variable{}, err
	}

	endpoint := fmt.Sprintf(APIEndpointVars, params.ProjectId, params.Key)
	resp, err := v.MakeRequest("PUT", endpoint, filter, data)
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

func (v *VarsAPI) DeleteVariable(params types.Params, filter types.Filter) error {
	err := params.Validate()
	if err != nil {
		return err
	}

	endpoint := fmt.Sprintf(APIEndpointVars, params.ProjectId, params.Key)
	_, err = v.MakeRequest("DELETE", endpoint, filter, types.VarData{})

	return err
}
