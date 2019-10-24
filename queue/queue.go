package queue

import (
	"container/list"
	"fmt"
	"io"
	"log"
	"os"
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
	appendLog  *os.File

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
	log.Printf("creating queue %s", key)
	appendLog := mustOpenFile(key)
	q := &Queue{
		key:        key,
		windowSize: windowSize,
		events:     list.New(),
		appendLog:  appendLog,
	}
	q.readEvents(appendLog)
	return q
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

	ev := &event{curTs}

	b, _ := ev.MarshalBinary()
	_, err := q.appendLog.Write(b)
	if err != nil {
		log.Printf("error writing to append log: %s", err)
	}

	q.addEvent(ev)
}

func (q *Queue) addEvent(ev *event) {
	oldf := q.events.Front()
	if oldf == nil {
		q.prevOldestTs = ev.timestamp
	} else {
		q.deltaSum += ev.timestamp.Sub(oldf.Value.(*event).timestamp).Milliseconds()
	}
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

func mustOpenFile(key string) *os.File {
	pattern := key + "_%04d.pip"
	name := fmt.Sprintf(pattern, 0)
	f, e := os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if e != nil {
		log.Fatalf("io error %s", e)
	}
	log.Printf("opened file %s", f.Name())
	return f
}

func (q *Queue) readEvents(f *os.File) {
	count := 0
	for {
		buf := make([]byte, 15)
		_, err := f.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Printf("unable to read bytes: %s", err)
			}
			break
		}
		var ev event
		err = ev.UnmarshalBinary(buf)
		if err != nil {
			log.Printf("error reading binary info: %s", err)
			break
		}
		q.addEvent(&ev)
		count++
	}
	log.Printf("added %d events into %s", count, q)
}

func (q *Queue) Close() {
	err := q.appendLog.Close()
	if err != nil {
		log.Printf("error closing file: %s", err)
	}
}
