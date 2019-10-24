package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	pb "spotify/pipe/schema"
	"sync"
	"time"
)

const (
	minDelayMs   = 1
	maxDelayMs   = 20
	batchDelayMs = 100
)

type demo struct {
	keys      int
	batch     bool
	counts    sync.Map
	deltas    sync.Map
	printLock sync.Mutex
	client    pb.PipeClient
}

func main() {
	addr := flag.String("h", "localhost:32232", "server addr")
	keys := flag.Int("k", 10, "keys")
	batch := flag.Bool("b", false, "batch mode")
	flag.Parse()

	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("error dialing %s - %s", *addr, err)
	}
	defer conn.Close()
	client := pb.NewPipeClient(conn)
	d := &demo{keys: *keys, batch: *batch, client: client}

	run(d)
}

func run(d *demo) {
	wg := &sync.WaitGroup{}

	if !d.batch {
		fmt.Println("pushing random events...")
		fmt.Print("\n")
		for {
			rnd := rand.Intn(d.keys)
			time.Sleep(time.Duration(rand.Intn(maxDelayMs-minDelayMs)+minDelayMs) * time.Millisecond)
			wg.Add(1)
			key := fmt.Sprintf("A%02d", rnd)
			go pushAndPrint(d.client, []string{key}, d, wg)
		}
	} else {
		fmt.Println("pushing batches...")
		fmt.Print("\n")
		keys := make([]string, d.keys)
		for k := 0; k < d.keys; k++ {
			keys[k] = fmt.Sprintf("A%02d", k)
		}
		for {
			time.Sleep(batchDelayMs * time.Millisecond)
			wg.Add(1)
			pushAndPrint(d.client, keys, d, wg)
		}
	}
	wg.Wait()
}

func pushAndPrint(client pb.PipeClient, keys []string, d *demo, wg *sync.WaitGroup) {
	defer wg.Done()
	r, err := client.Push(context.Background(), &pb.Events{Keys: keys})
	if err == nil && r != nil {
		for k, q := range r.Queues {
			d.counts.Store(k, q.Stats.Count)
			d.deltas.Store(k, q.Stats.AvgDeltaMs)
		}
		printout(d)
	}
}

func printout(d *demo) {
	d.printLock.Lock()
	defer d.printLock.Unlock()

	fmt.Print("\033[0;0H")
	for i := 0; i < d.keys; i++ {
		fmt.Print("\n")
		k := fmt.Sprintf("A%02d", i)
		c, _ := d.counts.LoadOrStore(k, uint64(0))
		f, _ := d.deltas.LoadOrStore(k, 0.0)
		fmt.Printf("%s[count: %03d, avg.delta: %7.2fms]: ", k, c, f.(float64))
		bar := int(float64(c.(uint64)) / 3)
		for k := 0; k < bar; k++ {
			fmt.Print("#")
		}
		fmt.Print("                     ")
	}
}
