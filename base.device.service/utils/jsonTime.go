package utils

import (
	"fmt"
	"time"
)

//JSONTime struct to desribe json formated timestamp
type JSONTime struct {
	time.Time
}

//MarshalJSON serialize timestamp
func (t *JSONTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", t.Format("2006-01-02T15:04:05Z"))
	return []byte(stamp), nil
}

//UnmarshalJSON deserialize timestamp
func (t *JSONTime) UnmarshalJSON(b []byte) error {
	sd := string(b[1 : len(b)-1])
	datetime, terr := time.Parse("2006-01-02T15:04:05Z", sd)
	t.Time = datetime
	return terr
}
