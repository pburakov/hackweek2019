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
	minDelayMs = 1
	maxDelayMs = 20
)

type demo struct {
	keys      int
	counts    sync.Map
	deltas    sync.Map
	printLock sync.Mutex
}

func main() {
	addr := flag.String("h", "localhost:32232", "server addr")
	keys := flag.Int("k", 10, "keys")
	flag.Parse()

	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("error dialing %s - %s", *addr, err)
	}
	defer conn.Close()
	client := pb.NewPipeClient(conn)
	wg := &sync.WaitGroup{}
	d := &demo{keys: *keys}

	fmt.Println("pushing events...")
	fmt.Print("\n")
	for {
		rnd := rand.Intn(d.keys)
		time.Sleep(time.Duration(rand.Intn(maxDelayMs-minDelayMs)+minDelayMs) * time.Millisecond)
		wg.Add(1)
		go pushAndPrint(client, fmt.Sprintf("A%02d", rnd), d, wg)
	}
	wg.Wait()
}

func pushAndPrint(client pb.PipeClient, key string, d *demo, wg *sync.WaitGroup) {
	defer wg.Done()
	q, err := client.Push(context.Background(), &pb.Event{Key: key})
	if err != nil || q == nil {
		fmt.Printf("error or nil response")
	} else {
		d.counts.Store(key, q.Stats.Count)
		d.deltas.Store(key, q.Stats.AvgDeltaMs)
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
