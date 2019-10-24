package queue

import (
	"testing"
	"time"
)

func TestEvent_MarshalBinary(t *testing.T) {
	curTs := time.Now()

	e1 := event{curTs}
	b, _ := e1.MarshalBinary()

	e2 := event{}
	_ = e2.UnmarshalBinary(b)

	if !e1.timestamp.Equal(e2.timestamp) {
		t.Fatalf("Expected timestamps to be equal but got %s and %s", e1.timestamp, e2.timestamp)
	}
}
