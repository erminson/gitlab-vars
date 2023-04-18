package types

import "encoding/json"

type APIResponse struct {
	Result json.RawMessage
}
