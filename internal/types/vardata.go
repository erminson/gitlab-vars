package types

import "regexp"

type VarData map[string]string

func (vd VarData) Validate() error {
	key, ok := vd["key"]
	if !ok {
		return ErrVarInvalidKey
	}

	re := regexp.MustCompile("^[a-zA-Z0-9_]*$")
	if key == "" || len(key) > 255 || !re.MatchString(key) {
		return ErrVarInvalidKey
	}

	value, ok := vd["value"]
	if !ok {
		return ErrVarInvalidValue
	}
	if value == "" {
		return ErrVarInvalidValue
	}

	vartype, ok := vd["variable_type"]
	if !ok {
		return ErrVarInvalidType
	}
	if vartype != string(VarTypeEnvVar) && vartype != string(VarTypeFile) {
		return ErrVarInvalidType
	}

	return nil
}
