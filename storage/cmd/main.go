package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	pb "karma/gen/storage"
	"karma/storage/service"
	"log"
	"net"
)

var port = flag.Int("port", 37000, "The server port")

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterStorageServer(s, &service.Service{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
