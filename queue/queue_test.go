package queue

import (
	"testing"
	"time"
)

func TestQueue_Tick(t *testing.T) {
	q := New("test", 1*time.Hour)
	q.Tick()

	if q.events.Len() != 1 {
		t.Errorf("Expected 1 element in the queue but got %d", q.events.Len())
	}
}

func TestQueue_AddTrimCount(t *testing.T) {
	q := New("test", 1*time.Hour)

	count := 5
	for i := 0; i < count-1; i++ {
		q.Tick()
	}

	cutoff := time.Now() // All events added earlier than this timestamp should be trimmed
	q.Tick()

	if q.events.Len() != count {
		t.Errorf("Expected %d elements in the queue but got %d", count, q.events.Len())
	}

	q.trimUntil(cutoff)

	if q.events.Len() != 1 {
		t.Errorf("Expected 1 element in the queue but got %d", q.events.Len())
	}
}

func TestQueue_TrimEmpty(t *testing.T) {
	q := New("test", 1*time.Hour)
	q.trimUntil(time.Now())
}

func TestQueue_Stats(t *testing.T) {
	q := New("test", 1*time.Hour)

	count := 5
	for i := 0; i < count; i++ {
		q.Tick()
	}

	if q.Stats().Count != count {
		t.Errorf("Expected %d elements in the queue but got %d", count, q.Stats().Count)
	}
}
