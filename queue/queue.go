package queue

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

type Queue struct {
	key        string
	windowSize time.Duration
	events     *list.List
	lock       sync.Mutex
}

type event struct {
	timestamp time.Time
}

func New(key string, windowSize time.Duration) *Queue {
	return &Queue{
		key:        key,
		windowSize: windowSize,
		events:     list.New(),
	}
}

func (q Queue) GetCount() int {
	q.trimUntil(time.Now().Add(-q.windowSize))
	return q.events.Len()
}

func (q *Queue) Tick() {
	event := &event{timestamp: time.Now()}
	q.push(event)
	q.Trim()
}

func (q *Queue) push(e *event) {
	q.lock.Lock()
	q.events.PushFront(e)
	q.lock.Unlock()
}

func (q Queue) String() string {
	return fmt.Sprintf("Queue[%s](%d)", q.key, q.events.Len())
}

func (q *Queue) Trim() {
	q.trimUntil(time.Now().Add(-q.windowSize))
}

func (q *Queue) trimUntil(cutoff time.Time) {
	q.lock.Lock()
	defer q.lock.Unlock()

	for elem := q.events.Back(); elem != nil; elem = q.events.Back() {
		event := elem.Value.(*event)
		if event.timestamp.After(cutoff) {
			break
		} else {
			q.events.Remove(elem)
		}
	}
}
