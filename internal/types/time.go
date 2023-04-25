package types

import (
	"fmt"
	"strings"
	"time"
)

var nilTime = (time.Time{}).UnixNano()

type Time time.Time

func (t *Time) UnmarshalJSON(data []byte) (err error) {
	s := strings.Trim(string(data), "\"")
	if s == "null" {
		*t = Time(time.Time{})
		return
	}

	tm, err := time.Parse(time.RFC3339, s)
	*t = Time(tm)

	return
}

func (t Time) MarshalJSON() ([]byte, error) {
	if time.Time(t).UnixNano() == nilTime {
		return []byte("null"), nil
	}

	return []byte(fmt.Sprintf("\"%s\"", time.Time(t).Format(time.RFC3339Nano))), nil
}

//
//func (t *Time) IsSet() bool {
//	return t.UnixNano() != nilTime
//}
