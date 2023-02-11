package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/spf13/afero"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pbMain "karma/gen/server"
	pb "karma/gen/storage"
	"karma/storage/service"
	"log"
	"net"
)

var host = flag.String("host", "127.0.0.1", "The grpc host")
var port = flag.Int("port", 37000, "The grpc port")
var server = flag.String("server", "127.0.0.1:37700", "The server address")
var capacity = flag.Uint64("capacity", 1024*1024, "The storage capacity")

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterStorageServer(s, service.NewService(afero.NewMemMapFs()))
	log.Printf("server listening at %v", lis.Addr())

	go addStorage()

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func addStorage() {
	conn, err := grpc.Dial(*server, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()
	c := pbMain.NewServerClient(conn)

	req := &pbMain.AddRequest{
		Address:  fmt.Sprintf("%s:%d", *host, *port),
		Capacity: *capacity,
	}

	_, err = c.AddStorage(context.Background(), req)
	if err != nil {
		log.Fatalf("could not add storage: %v", err)
	}
}
