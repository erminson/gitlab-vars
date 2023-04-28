package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

type VarType string

const (
	VarTypeEnvVar VarType = "env_var"
	VarTypeFile   VarType = "file"
)

var (
	ErrVariableInvalidKey   = errors.New("invalid variable key")
	ErrVariableInvalidValue = errors.New("invalid variable value")
	ErrVariableInvalidType  = errors.New("invalid variable type")
)

type Variable struct {
	Type             string `json:"variable_type"`
	Key              string `json:"key"`
	Value            string `json:"value"`
	Protected        bool   `json:"protected"`
	Masked           bool   `json:"masked"`
	Raw              bool   `json:"raw"`
	EnvironmentScope string `json:"environment_scope"`
}

func (v *Variable) String() string {
	return fmt.Sprintf("Variable { Type:%v, Key:%v, Value:%v, Protected:%v, Masked:%v, Raw:%v, EnvironmentScope:%v}",
		v.Type,
		v.Key,
		v.Value,
		v.Protected,
		v.Masked,
		v.Raw,
		v.EnvironmentScope,
	)
}

func (v *Variable) VariableToData() VarData {
	vd := VarData{
		"variable_type":     v.Type,
		"key":               v.Key,
		"value":             v.Value,
		"protected":         strconv.FormatBool(v.Protected),
		"masked":            strconv.FormatBool(v.Masked),
		"raw":               strconv.FormatBool(v.Raw),
		"environment_scope": v.EnvironmentScope,
	}

	return vd
}

func (v *Variable) Validate() error {
	re := regexp.MustCompile("^[a-zA-Z0-9_]*$")
	if v.Key == "" || len(v.Key) > 255 || !re.MatchString(v.Key) {
		return ErrVariableInvalidKey
	}

	if v.Value == "" {
		return ErrVariableInvalidValue
	}

	if v.Type != string(VarTypeEnvVar) && v.Type != string(VarTypeFile) {
		return ErrVariableInvalidType
	}

	return nil
}

func (v *Variable) UnmarshalJSON(b []byte) error {
	type innerVar Variable
	out := &innerVar{
		EnvironmentScope: "*",
	}

	if err := json.Unmarshal(b, out); err != nil {
		return err
	}

	*v = Variable(*out)

	return nil
}
