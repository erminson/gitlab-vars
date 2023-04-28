package types

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParams_Key_ProjectId_String(t *testing.T) {
	key := "var"
	var projectId int64 = 123

	params := Params{
		Key:       key,
		ProjectId: projectId,
	}

	expString := fmt.Sprintf("Params { ProjectId:%v, Key:%v }", projectId, key)

	require.Equal(t, expString, params.String())
}

func TestParams_ProjectId_String(t *testing.T) {
	var projectId int64 = 123

	params := Params{
		ProjectId: projectId,
	}

	expString := fmt.Sprintf("Params { ProjectId:%v }", projectId)

	require.Equal(t, expString, params.String())
}

func TestValidateParams(t *testing.T) {
	params := Params{
		Key:       "key",
		ProjectId: 1,
	}

	err := params.Validate()
	require.NoError(t, err)
}

func TestValidateParamsError(t *testing.T) {
	cases := []struct {
		name     string
		p        Params
		expError error
	}{
		{
			name: "Invalid Params Project Id. 0",
			p: Params{
				ProjectId: 0,
			},
			expError: ErrParamsInvalidProjectId,
		},
		{
			name: "Invalid Params Project Id. -100",
			p: Params{
				ProjectId: -100,
			},
			expError: ErrParamsInvalidProjectId,
		},
		{
			name: "Invalid Params Key. Empty",
			p: Params{
				ProjectId: 1,
				Key:       "",
			},
			expError: ErrParamsInvalidKey,
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			err := tCase.p.Validate()
			require.Error(t, err)
			require.EqualError(t, err, tCase.expError.Error())
		})
	}
}

func TestValidateParamsProjectId(t *testing.T) {
	params := Params{
		ProjectId: 1,
	}

	err := params.ValidateProjectId()
	require.NoError(t, err)
}

func TestValidateParamsProjectId_Error(t *testing.T) {
	cases := []struct {
		name   string
		p      Params
		expErr error
	}{
		{
			name:   "Invalid Project Id. 0",
			p:      Params{ProjectId: 0},
			expErr: ErrParamsInvalidProjectId,
		},
		{
			name:   "Invalid Project Id. -100",
			p:      Params{ProjectId: -100},
			expErr: ErrParamsInvalidProjectId,
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			err := tCase.p.ValidateProjectId()
			require.Error(t, err)
			require.EqualError(t, err, tCase.expErr.Error())
		})
	}
}
