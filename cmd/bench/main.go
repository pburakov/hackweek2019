package main

import (
	"log"
	"os"
	"spotify/pipe/queue"
	"time"
)

func main() {
	log.SetOutput(os.Stdout)

	q := queue.New("test_queue", 1*time.Second)

	for i := 0; i < 10000; i++ {
		q.Tick()
		log.Println(q)
		time.Sleep(10 * time.Millisecond)
	}
}
