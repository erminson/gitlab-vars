package types

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidateVarData(t *testing.T) {
	cases := []struct {
		name string
		v    VarData
	}{
		{
			name: "Valid VarData #1",
			v: VarData{
				"variable_type": string(VarTypeEnvVar),
				"key":           "var_key",
				"value":         "var_value",
				"protected":     "true",
				"masked":        "false",
				"raw":           "true",
			},
		},
		{
			name: "Valid VarData #2",
			v: VarData{
				"variable_type": string(VarTypeFile),
				"key":           "var_key",
				"value":         "var_value",
				"protected":     "false",
				"masked":        "true",
				"raw":           "false",
			},
		},
		{
			name: "Valid VarData #3",
			v: VarData{
				"variable_type":     string(VarTypeFile),
				"key":               "var_key123",
				"value":             "var_value",
				"environment_scope": "*",
			},
		},
		{
			name: "Valid VarData #4",
			v: VarData{
				"variable_type": string(VarTypeFile),
				"key":           "qi5S0evHkCYW7VwK3y3w6e4afsZGZkeURjxZmfGVt2EQhax6JkD_Kgub5CRnbiaAiY7_AbXSYg6qhlRPFRxMKXwJAAVw6vWcAk_FX1TcOTVXsaVOunhffTsNKiDd1T8nNOxnZAUj0GHnuefcmobCTai3AOx5ykwX2juNkgcJBHGI6Vm60lgekXuUwdKDDa9PBwqDO9hiopN0PXAiwKTFpgU8pq7bjlJA4_NpImyfiOM2zQ9rF7aIVF8LMce7n0o",
				"value":         "var_value",
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

func TestValidateVarDataError(t *testing.T) {
	cases := []struct {
		name   string
		v      VarData
		expErr error
	}{
		{
			name:   "Invalid Variable Key. Missing Key",
			v:      VarData{},
			expErr: ErrVarDataInvalidKey,
		},
		{
			name: "Invalid Variable Key. Empty Key",
			v: VarData{
				"key": "",
			},
			expErr: ErrVarDataInvalidKey,
		},
		{
			name: "Invalid Variable Key. More than 255 characters",
			v: VarData{
				"key": "GVbrPSGqeYyyLaGIM2ehFIWgHdGOHK62eNSyJ7nK6MgdgWJZaZhbbQbdk0C6YqeKInuh8axI8lodhqzGphXkubiWF2pNtiBt3gPRq7BatFi3OLJTVOlLnbegTkao3KCSYq9sYC9Oz9JLAh9kEaUWhmuYbhX1JrlsLMBoEhBxNKUfQHnVOimk4NXY7oWmV7kxnhpmRd2sYoMbWaH20WCONMCj0UdsPgS8SRsyJ5wNnwHR8dLSNubkc3jvLZDKXPwe",
			},
			expErr: ErrVarDataInvalidKey,
		},
		{
			name: "Invalid Variable Key. Use spaces",
			v: VarData{
				"key": " key ",
			},
			expErr: ErrVarDataInvalidKey,
		},
		{
			name: "Invalid Variable Key. Use %",
			v: VarData{
				"key": "key%",
			},
			expErr: ErrVarDataInvalidKey,
		},
		{
			name: "Invalid Variable Key. Use $",
			v: VarData{
				"key": "key$",
			},
			expErr: ErrVarDataInvalidKey,
		},
		{
			name: "Invalid Variable Value. Missing Value",
			v: VarData{
				"key": "key",
			},
			expErr: ErrVarDataInvalidValue,
		},
		{
			name: "Invalid Variable Value. Missing Value",
			v: VarData{
				"key":   "key",
				"value": "",
			},
			expErr: ErrVarDataInvalidValue,
		},
		{
			name: "Invalid Variable Type. Missing Type",
			v: VarData{
				"key":   "key",
				"value": "value",
			},
			expErr: ErrVarDataInvalidType,
		},
		{
			name: "Invalid Variable Type.",
			v: VarData{
				"key":           "key",
				"value":         "value",
				"variable_type": "not_correct_type",
			},
			expErr: ErrVarDataInvalidType,
		},
		{
			name: "Invalid Variable Protected. abc",
			v: VarData{
				"key":           "key",
				"value":         "value",
				"variable_type": string(VarTypeFile),
				"protected":     "abc",
			},
			expErr: ErrVarDataInvalidProtected,
		},
		{
			name: "Invalid Variable Masked. abc",
			v: VarData{
				"key":           "key",
				"value":         "value",
				"variable_type": string(VarTypeFile),
				"protected":     "true",
				"masked":        "abc",
			},
			expErr: ErrVarDataInvalidMasked,
		},
		{
			name: "Invalid Variable Raw. abc",
			v: VarData{
				"key":           "key",
				"value":         "value",
				"variable_type": string(VarTypeFile),
				"protected":     "true",
				"masked":        "false",
				"raw":           "abc",
			},
			expErr: ErrVarDataInvalidRaw,
		},
		{
			name: "Invalid Variable EnvScope. Empty",
			v: VarData{
				"key":               "key",
				"value":             "value",
				"variable_type":     string(VarTypeFile),
				"protected":         "true",
				"masked":            "false",
				"raw":               "true",
				"environment_scope": "",
			},
			expErr: ErrVarDataInvalidEnvScope,
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
