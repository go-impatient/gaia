package time

import (
	"fmt"
	"time"
)

// Time is based on time.Time.
//
// FYI: OF's domain type in this case is literally "time.Time".
type Time struct {
	time time.Time
}

// NewTime returns a new instance of Time in UTC timezone.
func NewTime(t time.Time) Time {
	return Time{time: t.UTC()}
}

// MarshalJSON implements the Marshaler interface in encoding/json.
// It leverages the String method of Time.
func (t Time) MarshalJSON() ([]byte, error) {
	// https://stackoverflow.com/a/23695774
	return []byte(fmt.Sprintf("\"%s\"", t)), nil
}

// String implements the fmt.Stringer interface.
//
// String returns Time as a string in UTC timezone and RFC3339 format.
func (t Time) String() string {
	return time.Time(t.time).UTC().Format(time.RFC3339)
}
