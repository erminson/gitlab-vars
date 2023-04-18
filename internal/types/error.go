package types

import "fmt"

type APIError struct {
	Message ErrorMessage `json:"message"`
	Code    int
}

func (e APIError) Error() string {
	return fmt.Sprintf("Code: %d Message: %s", e.Code, e.Message)
}

type ErrorMessage struct {
	string
}

func (e *ErrorMessage) UnmarshalJSON(b []byte) (err error) {
	e.string = string(b)

	return
}
