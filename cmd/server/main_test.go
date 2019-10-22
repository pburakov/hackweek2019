package main

import (
	"fmt"
	"spotify/pipe/queue"
	"testing"
)

func TestPipe_Push(t *testing.T) {
	key := "test_key"

	s := new(server)
	s.getOrInsert(key)

	a, ok := s.queues.Load(key)
	if !ok {
		t.Fatalf("Expected %q to be present in the map", key)
	}
	q := a.(*queue.Queue)
	fmt.Println(q.String())
}
