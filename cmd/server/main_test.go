package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"spotify/pipe/queue"
	pb "spotify/pipe/schema"
	"testing"
)

func TestServer_Push(t *testing.T) {
	s := new(server)
	keys := []string{"key1", "key2"}

	for i := uint64(1); i <= 5; i++ {
		resp, _ := s.Push(context.Background(), &pb.Events{Keys: keys})
		count1 := resp.Queues["key1"].Stats.Count
		count2 := resp.Queues["key2"].Stats.Count
		if count1 != i || count2 != i {
			t.Fatalf("Expected count for both keys to be %d, but got %d and %d", i, count1, count2)
		}
	}
}

func TestServer_GetOrInsert(t *testing.T) {
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
