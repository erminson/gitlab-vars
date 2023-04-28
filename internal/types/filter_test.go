package types

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFilter_String(t *testing.T) {
	expString := "Filter { k1:v1 }"

	filter := Filter{"k1": "v1"}

	require.Equal(t, expString, filter.String())
}
