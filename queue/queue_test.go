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

	baseTS := time.Now()
	recent := baseTS.Add(-1 * time.Hour)
	cutoff := baseTS.Add(-2 * time.Hour) // events older than this timestamp should be evicted (inclusive)
	older1 := baseTS.Add(-3 * time.Hour)
	older2 := baseTS.Add(-4 * time.Hour)
	older3 := baseTS.Add(-5 * time.Hour)

	q.push(&event{timestamp: older3})
	q.push(&event{timestamp: older2})
	q.push(&event{timestamp: older1})
	q.push(&event{timestamp: cutoff})
	q.push(&event{timestamp: recent})

	if q.events.Len() != 5 {
		t.Errorf("Expected 5 elements in the queue but got %d", q.events.Len())
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
