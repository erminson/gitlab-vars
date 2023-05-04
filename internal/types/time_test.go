package types

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestUnmarshalTime(t *testing.T) {
	timeStr := "2023-04-15T14:25:18.085Z"
	expTime, err := time.Parse(time.RFC3339, timeStr)
	require.NoError(t, err)

	in := fmt.Sprintf(`{"created_at": "%s"}`, timeStr)
	var out struct {
		CreatedAt Time `json:"created_at"`
	}

	err = json.Unmarshal([]byte(in), &out)
	require.Equal(t, time.Time(out.CreatedAt), expTime)
	require.NoError(t, err)
}

func TestUnmarshalTime_Null(t *testing.T) {
	in := `{"created_at": "null"}`

	var out struct {
		CreatedAt Time `json:"created_at"`
	}

	err := json.Unmarshal([]byte(in), &out)
	require.Equal(t, Time{}, out.CreatedAt)
	require.NoError(t, err)
}

func TestUnmarshalTimeError(t *testing.T) {
	in := `{"created_at": ""}`

	var out struct {
		CreatedAt Time `json:"created_at"`
	}

	err := json.Unmarshal([]byte(in), &out)
	require.Error(t, err)
}

func TestMarshalTime(t *testing.T) {
	timeStr := `2023-04-15T14:25:18.085Z`
	expTimeData := []byte(fmt.Sprintf("\"%s\"", timeStr))

	tm, _ := time.Parse(time.RFC3339Nano, timeStr)
	in := Time(tm)

	actualTimeData, err := json.Marshal(in)
	require.NoError(t, err)
	require.Equal(t, expTimeData, actualTimeData)
}

func TestMarshalTime_Null(t *testing.T) {
	timeStr := "null"
	expTimeData := []byte(timeStr)

	in := Time(time.Time{})

	actualTimeData, err := json.Marshal(in)
	require.NoError(t, err)
	require.Equal(t, expTimeData, actualTimeData)
}
