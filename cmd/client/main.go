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

var numKeys = 15

type demo struct {
	results   sync.Map
	printLock sync.Mutex
}

func main() {
	addr := flag.String("host", "localhost:32232", "server addr")
	flag.Parse()

	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("error dialing %s - %s", *addr, err)
	}
	defer conn.Close()
	client := pb.NewPipeClient(conn)
	wg := &sync.WaitGroup{}
	d := new(demo)

	log.Println("pushing events...")
	fmt.Print("\n")
	for {
		rnd := rand.Intn(numKeys)
		time.Sleep(5 * time.Millisecond)
		wg.Add(1)
		go pushAndPrint(client, fmt.Sprintf("key%d", rnd), d, wg)
	}
	wg.Wait()
}

func pushAndPrint(client pb.PipeClient, key string, d *demo, wg *sync.WaitGroup) {
	defer wg.Done()
	q, err := client.Push(context.Background(), &pb.Event{Key: key})
	if err != nil || q == nil {
		panic("error or nil response")
	} else {
		d.results.Store(key, q.Stats.Count)
		printout(d)
	}
}

func printout(d *demo) {
	d.printLock.Lock()
	defer d.printLock.Unlock()

	fmt.Print("\033[0;0H")
	for i := 0; i < numKeys; i++ {
		fmt.Print("\n")
		k := fmt.Sprintf("key%d", i)
		v, _ := d.results.LoadOrStore(k, uint64(0))
		fmt.Printf("%s[%00d]", k, v)
		x := int(v.(uint64) / 5)
		for k := 0; k < x; k++ {
			fmt.Print("#")
		}
		fmt.Print("              ")
	}
}
