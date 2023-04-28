package types

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAPIError_Error(t *testing.T) {
	message := "Error Message"
	code := 501
	expString := fmt.Sprintf("Code: %d Message: %s", code, message)

	err := APIError{
		Message: ErrorMessage(message),
		Code:    code,
	}

	require.Equal(t, expString, err.Error())
}

func TestErrorMessage_UnmarshalJSON(t *testing.T) {
	errMsg := "error message"
	expErrStr := fmt.Sprintf(`{"message" : "%s"}`, errMsg)

	var expErr APIError

	err := json.Unmarshal([]byte(expErrStr), &expErr)
	require.NoError(t, err)
	require.Equal(t, errMsg, string(expErr.Message))
}
