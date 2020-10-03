package time

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Confirm that wrap.Time implements the fmt.Stringer interface.
func TestTime_StringerInterface(t *testing.T) {
	wrapTime := NewTime(time.Time{})
	var _ fmt.Stringer = wrapTime
	assert.Nil(t, nil) // If we get this far, the test passed.
}

// Simple test of Time's implementation of the fmt.Stringer interface.
func TestTime_StringerSimple(t *testing.T) {
	wrapTime := NewTime(time.Date(2000, time.February, 29, 1, 2, 3, 4, time.UTC))
	assert.Equal(t, "2000-02-29T01:02:03Z", fmt.Sprintf("%s", wrapTime))
}

// Simple test of Time's implementation of the Marshaler interface in
// encoding/json.
func TestTime_JSONMarshalerSimple(t *testing.T) {
	assert := assert.New(t) // Prepare to assert multiple times
	input := time.Date(2000, time.February, 29, 1, 2, 3, 4, time.UTC)
	wrapTime := NewTime(input)
	result, err := wrapTime.MarshalJSON()
	assert.Nil(err, "Time.MarshalJSON() returned non-nil error")
	expect := []byte("\"2000-02-29T01:02:03Z\"")
	assert.Equal(expect, result, "Did not obtain expected result.")
}

// Confirm that wrap.Time implements the Marshaler interface in encoding/json.
func TestTime_JSONMarshalerInterface(t *testing.T) {
	wrapTime := NewTime(time.Time{})
	var _ json.Marshaler = wrapTime
	assert.Nil(t, nil) // If we get this far, the test passed.
}

