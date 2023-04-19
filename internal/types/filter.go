package types

import (
	"bytes"
	"fmt"
)

const FilterEnvScope = "environment_scope"

type Filter map[string]string

func (f Filter) String() string {
	var buffer bytes.Buffer
	for k, v := range f {
		buffer.WriteString(fmt.Sprintf("%s:%s ", k, v))
	}

	return fmt.Sprintf("Filter { %s}", buffer.String())
}
