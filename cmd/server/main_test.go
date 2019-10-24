package main

import (
	"fmt"
	"os"
	"path/filepath"
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

	cleanup()
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
