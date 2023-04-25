package types

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidate(t *testing.T) {
	cases := []struct {
		name string
		v    *Variable
	}{
		{
			name: "Valid Variable #1",
			v: &Variable{
				Type:  string(VarTypeEnvVar),
				Key:   "var_key",
				Value: "var_value",
			},
		},
		{
			name: "Valid Variable #2",
			v: &Variable{
				Type:  string(VarTypeFile),
				Key:   "var_key",
				Value: "var_value",
			},
		},
		{
			name: "Valid Variable #3",
			v: &Variable{
				Type:  string(VarTypeFile),
				Key:   "var_key123",
				Value: "var_value",
			},
		},
		{
			name: "Valid Variable #4",
			v: &Variable{
				Type:  string(VarTypeFile),
				Key:   "qi5S0evHkCYW7VwK3y3w6e4afsZGZkeURjxZmfGVt2EQhax6JkD_Kgub5CRnbiaAiY7_AbXSYg6qhlRPFRxMKXwJAAVw6vWcAk_FX1TcOTVXsaVOunhffTsNKiDd1T8nNOxnZAUj0GHnuefcmobCTai3AOx5ykwX2juNkgcJBHGI6Vm60lgekXuUwdKDDa9PBwqDO9hiopN0PXAiwKTFpgU8pq7bjlJA4_NpImyfiOM2zQ9rF7aIVF8LMce7n0o",
				Value: "var_value",
			},
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			err := tCase.v.Validate()
			require.NoError(t, err)
		})
	}
}

func TestValidateError(t *testing.T) {
	cases := []struct {
		name   string
		v      *Variable
		expErr error
	}{
		{
			name: "Invalid Variable Key. Empty",
			v: &Variable{
				Key: "",
			},
			expErr: ErrVarInvalidKey,
		},
		{
			name: "Invalid Variable Key. More than 255 characters",
			v: &Variable{
				Key: "GVbrPSGqeYyyLaGIM2ehFIWgHdGOHK62eNSyJ7nK6MgdgWJZaZhbbQbdk0C6YqeKInuh8axI8lodhqzGphXkubiWF2pNtiBt3gPRq7BatFi3OLJTVOlLnbegTkao3KCSYq9sYC9Oz9JLAh9kEaUWhmuYbhX1JrlsLMBoEhBxNKUfQHnVOimk4NXY7oWmV7kxnhpmRd2sYoMbWaH20WCONMCj0UdsPgS8SRsyJ5wNnwHR8dLSNubkc3jvLZDKXPwe",
			},
			expErr: ErrVarInvalidKey,
		},
		{
			name: "Invalid Variable Key. Use spaces",
			v: &Variable{
				Key: " key ",
			},
			expErr: ErrVarInvalidKey,
		},
		{
			name: "Invalid Variable Key. Use %",
			v: &Variable{
				Key: "key%",
			},
			expErr: ErrVarInvalidKey,
		},
		{
			name: "Invalid Variable Key. Use $",
			v: &Variable{
				Key: "key$",
			},
			expErr: ErrVarInvalidKey,
		},
		{
			name: "Invalid Variable Value.",
			v: &Variable{
				Key:   "key",
				Value: "",
			},
			expErr: ErrVarInvalidValue,
		},
		{
			name: "Invalid Variable Type.",
			v: &Variable{
				Key:   "key",
				Value: "value",
				Type:  "not_correct_type",
			},
			expErr: ErrVarInvalidType,
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			err := tCase.v.Validate()
			require.Error(t, err)
			require.EqualError(t, err, tCase.expErr.Error())
		})
	}
}

func TestVariable_String(t *testing.T) {
	v := Variable{
		Type:             string(VarTypeEnvVar),
		Key:              "key",
		Value:            "value",
		Protected:        true,
		Masked:           false,
		Raw:              false,
		EnvironmentScope: "prod",
	}

	expected := "Variable { Type:env_var, Key:key, Value:value, Protected:true, Masked:false, Raw:false, EnvironmentScope:prod}"
	require.Equal(t, expected, v.String())
}

func TestVariable_VariableToData(t *testing.T) {
	vd := VarData{
		"variable_type":     "env_var",
		"key":               "key",
		"value":             "value",
		"protected":         "true",
		"masked":            "false",
		"raw":               "false",
		"environment_scope": "prod",
	}

	v := Variable{
		Type:             string(VarTypeEnvVar),
		Key:              "key",
		Value:            "value",
		Protected:        true,
		Masked:           false,
		Raw:              false,
		EnvironmentScope: "prod",
	}

	require.Equal(t, vd, v.VariableToData())
}
