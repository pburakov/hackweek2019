package queue

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

type Stats struct {
	Count      int
	AvgDeltaMs float64
}

func (s Stats) String() string {
	return fmt.Sprintf("count: %d, avg.delta: %fms", s.Count, s.AvgDeltaMs)
}

type Queue struct {
	key        string
	windowSize time.Duration
	events     *list.List
	lock       sync.Mutex

	// internal stats (updated on queue change)
	count        int
	deltaSum     int64
	prevOldestTs time.Time
	avg          float64
}

func (q Queue) String() string {
	return fmt.Sprintf("Queue[%s]", q.key)
}

func New(key string, windowSize time.Duration) *Queue {
	return &Queue{
		key:        key,
		windowSize: windowSize,
		events:     list.New(),
	}
}

func (q *Queue) Stats() *Stats {
	q.trimUntil(time.Now().Add(-q.windowSize))

	stats := &Stats{Count: q.count, AvgDeltaMs: q.avg}
	return stats
}

func (q *Queue) Tick() {
	q.lock.Lock()
	defer q.lock.Unlock()

	curTs := time.Now()

	oldf := q.events.Front()
	if oldf == nil {
		q.prevOldestTs = curTs
	} else {
		q.deltaSum += curTs.Sub(oldf.Value.(*event).timestamp).Milliseconds()
	}
	ev := &event{curTs}
	q.events.PushFront(ev)
}

func (q *Queue) trimUntil(cutoff time.Time) {
	q.lock.Lock()
	defer q.lock.Unlock()

	for elem := q.events.Back(); elem != nil; elem = q.events.Back() {
		ev := elem.Value.(*event)
		ts := ev.timestamp

		if ts.After(cutoff) {
			break
		} else {
			q.events.Remove(elem)
			q.deltaSum = q.deltaSum - ts.Sub(q.prevOldestTs).Milliseconds()
			q.prevOldestTs = ts
		}
	}
	q.count = q.events.Len()
	if q.count != 0 {
		q.avg = float64(q.deltaSum) / float64(q.count)
	}
}
