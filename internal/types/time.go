package types

import (
	"fmt"
	"strings"
	"time"
)

var nilTime = (time.Time{}).UnixNano()

type Time struct {
	time.Time
}

func (t *Time) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		t.Time = time.Time{}
		return
	}

	t.Time, err = time.Parse(time.RFC3339, s)

	return
}

func (t *Time) MarshalJSON() ([]byte, error) {
	if t.Time.UnixNano() == nilTime {
		return []byte("null"), nil
	}

	return []byte(fmt.Sprintf("\"%s\"", t.Time.Format(time.RFC3339))), nil
}

//
//func (t *Time) IsSet() bool {
//	return t.UnixNano() != nilTime
//}
