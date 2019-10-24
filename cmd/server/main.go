package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"spotify/pipe/queue"
	pb "spotify/pipe/schema"
	"sync"
	"time"
)

type server struct {
	pb.UnimplementedPipeServer

	queues sync.Map
}

func (s *server) Push(ctx context.Context, req *pb.Event) (*pb.Queue, error) {
	q := s.getOrInsert(req.Key)
	q.Tick()
	return convert(q), nil
}

func (s *server) Batch(ctx context.Context, req *pb.Batch) (*pb.Multi, error) {
	queues := make([]*pb.Queue, 0)
	for _, ev := range req.Events {
		v, ok := s.queues.Load(ev.Key)
		if ok {
			queues = append(queues, convert(v.(*queue.Queue)))
		}
	}
	return &pb.Multi{Queues: queues}, nil
}

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	log.Printf("pid: %d", os.Getpid())
	port := flag.Int("port", 32232, "server port")
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	grpcServer := grpc.NewServer()

	s := &server{}
	pb.RegisterPipeServer(grpcServer, s)

	log.Printf("Starting server on port %d", *port)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("error starting server: %s", err)
	}
}

func (s *server) getOrInsert(key string) *queue.Queue {
	q, found := s.queues.Load(key)
	if !found {
		n := queue.New(key, 10*time.Second)
		q, _ = s.queues.LoadOrStore(key, n)
	}
	return q.(*queue.Queue)
}

func convert(q *queue.Queue) *pb.Queue {
	stats := q.Stats()
	return &pb.Queue{
		Stats: &pb.Stats{
			Count:      uint64(stats.Count),
			AvgDeltaMs: stats.AvgDeltaMs,
		},
	}
}
