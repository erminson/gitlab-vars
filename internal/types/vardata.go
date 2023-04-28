package types

import (
	"errors"
	"regexp"
)

type VarData map[string]string

var (
	ErrVarDataInvalidKey   = errors.New("invalid variable key")
	ErrVarDataInvalidValue = errors.New("invalid variable value")
	ErrVarDataInvalidType  = errors.New("invalid variable type")

	ErrVarDataInvalidProtected = errors.New("invalid vardata protected")
	ErrVarDataInvalidMasked    = errors.New("invalid vardata masked")
	ErrVarDataInvalidRaw       = errors.New("invalid vardata raw")
	ErrVarDataInvalidEnvScope  = errors.New("invalid vardata env scope")
)

func (vd VarData) Validate() error {
	key, ok := vd["key"]
	if !ok {
		return ErrVarDataInvalidKey
	}

	re := regexp.MustCompile("^[a-zA-Z0-9_]*$")
	if key == "" || len(key) > 255 || !re.MatchString(key) {
		return ErrVarDataInvalidKey
	}

	value, ok := vd["value"]
	if !ok {
		return ErrVarDataInvalidValue
	}
	if value == "" {
		return ErrVarDataInvalidValue
	}

	vartype, ok := vd["variable_type"]
	if !ok {
		return ErrVarDataInvalidType
	}
	if vartype != string(VarTypeEnvVar) && vartype != string(VarTypeFile) {
		return ErrVarDataInvalidType
	}

	protected, ok := vd["protected"]
	if ok {
		if protected != "false" && protected != "true" {
			return ErrVarDataInvalidProtected
		}
	}

	masked, ok := vd["masked"]
	if ok {
		if masked != "false" && masked != "true" {
			return ErrVarDataInvalidMasked
		}
	}

	raw, ok := vd["raw"]
	if ok {
		if raw != "false" && raw != "true" {
			return ErrVarDataInvalidRaw
		}
	}

	envScope, ok := vd["environment_scope"]
	if ok {
		if envScope == "" {
			return ErrVarDataInvalidEnvScope
		}
	}

	return nil
}
