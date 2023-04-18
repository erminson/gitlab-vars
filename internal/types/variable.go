package types

import "strconv"

type Variable struct {
	Type             string `json:"variable_type"`
	Key              string `json:"key"`
	Value            string `json:"value"`
	Protected        bool   `json:"protected"`
	Masked           bool   `json:"masked"`
	Raw              bool   `json:"raw"`
	EnvironmentScope string `json:"environment_scope"`
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
