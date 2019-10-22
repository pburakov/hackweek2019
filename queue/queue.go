package queue

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

type Stats struct {
	Count       int
	NewestEvent *time.Time
	OldestEvent *time.Time
}

func (s Stats) String() string {
	return fmt.Sprintf("count: %d, newest: %s, oldest: %s",
		s.Count, s.NewestEvent, s.OldestEvent)
}

type Queue struct {
	key        string
	windowSize time.Duration
	events     *list.List
	lock       sync.Mutex

	// internal stats (updated on queue change)
	count   int
	avgDist time.Duration
	first   *list.Element
	last    *list.Element
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

	stats := &Stats{Count: q.count}
	f := q.first
	if f != nil {
		stats.NewestEvent = &f.Value.(*event).timestamp
	} else {
		stats.NewestEvent = nil
	}
	b := q.last
	if b != nil {
		stats.OldestEvent = &b.Value.(*event).timestamp
	} else {
		stats.OldestEvent = nil
	}
	return stats
}

func (q *Queue) Tick() {
	q.lock.Lock()
	defer q.lock.Unlock()

	ev := &event{timestamp: time.Now()}
	q.events.PushFront(ev)
	q.updateStats()
}

func (q *Queue) trimUntil(cutoff time.Time) {
	q.lock.Lock()
	defer q.lock.Unlock()

	for elem := q.events.Back(); elem != nil; elem = q.events.Back() {
		ev := elem.Value.(*event)
		if ev.timestamp.After(cutoff) {
			break
		} else {
			q.events.Remove(elem)
		}
	}
	q.updateStats()
}

func (q *Queue) updateStats() {
	q.count = q.events.Len()
	q.first = q.events.Front()
	q.last = q.events.Back()
}
