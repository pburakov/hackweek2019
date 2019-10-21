package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	pb "spotify/pipe/schema"
)

type server struct {
	pb.UnimplementedPipeServer
}

func main() {
	port := flag.Int("port", 32232, "server port")
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	grpcServer := grpc.NewServer()
	pb.RegisterPipeServer(grpcServer, &server{})

	log.Println("Starting server")
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("error starting server: %s", err)
	}
}
