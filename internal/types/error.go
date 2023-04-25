package types

import (
	"fmt"
	"strings"
)

type APIError struct {
	Message ErrorMessage `json:"message"`
	Code    int
}

func (e APIError) Error() string {
	return fmt.Sprintf("Code: %d Message: %s", e.Code, e.Message)
}

type ErrorMessage string

func (e *ErrorMessage) UnmarshalJSON(b []byte) (err error) {
	msg := strings.Trim(string(b), "\"")
	*e = ErrorMessage(msg)

	return
}
