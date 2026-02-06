package types

import "net/http"

type APIResponse struct {
	Result  []byte
	Headers http.Header
	Status  int
}
