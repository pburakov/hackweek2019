package queue

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

const maxDur = 1 * time.Hour

func TestQueue_Tick(t *testing.T) {
	defer cleanup()

	q := New("test", maxDur)
	defer q.Close()

	q.Tick()

	if q.events.Len() != 1 {
		t.Errorf("Expected 1 element in the queue but got %d", q.events.Len())
	}
}

func TestQueue_TrimNow(t *testing.T) {
	defer cleanup()

	wait := 500 * time.Millisecond
	q := New("test", wait)
	defer q.Close()

	q.Tick()

	if q.events.Len() != 1 {
		t.Errorf("Expected 1 element in the queue but got %d", q.events.Len())
	}

	time.Sleep(wait)
	q.TrimNow()

	if q.events.Len() != 0 {
		t.Errorf("Expected 0 elements in the queue but got %d", q.events.Len())
	}
}

func TestQueue_AddTrimCount(t *testing.T) {
	defer cleanup()

	q := New("test", maxDur)
	defer q.Close()

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
	defer cleanup()

	q := New("test", maxDur)
	defer q.Close()

	q.trimUntil(time.Now())
}

func TestQueue_Stats(t *testing.T) {
	defer cleanup()

	q := New("test", maxDur)
	defer q.Close()

	count := 5
	for i := 0; i < count; i++ {
		q.Tick()
	}

	if q.Stats().Count != count {
		t.Errorf("Expected %d elements in the queue but got %d", count, q.Stats().Count)
	}
}

func TestQueue_IO(t *testing.T) {
	defer cleanup()

	q1 := New("test", maxDur)
	count := 5
	for i := 0; i < count; i++ {
		q1.Tick()
	}
	q1.Close()

	q2 := New("test", maxDur)
	if q2.events.Len() != count {
		t.Errorf("Expected %d elements in the queue but got %d", count, q2.events.Len())
	}
	q2.Close()

	// verify events read in the same order
	elem1 := q1.events.Back()
	elem2 := q2.events.Back()
	for elem1 != nil && elem2 != nil {
		ev1 := elem1.Value.(*event)
		ev2 := elem2.Value.(*event)
		if !ev1.timestamp.Equal(ev2.timestamp) {
			t.Errorf("Expected events in both queues stored the same order")
		}
		elem1 = elem1.Prev()
		elem2 = elem2.Prev()
	}
}

func cleanup() {
	files, err := filepath.Glob("*.pip")
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		if err := os.Remove(f); err != nil {
			panic(err)
		}
	}
}
