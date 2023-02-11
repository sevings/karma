package service

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "karma/gen/server"
)

type Service struct {
	pb.UnimplementedServerServer
	conns   []grpc.ClientConnInterface
	clients map[string]pb.ServerClient
}

func NewService() *Service {
	return &Service{
		conns:   make([]grpc.ClientConnInterface, 0, 5),
		clients: make(map[string]pb.ServerClient),
	}
}

func (s *Service) AddStorage(ctx context.Context, req *pb.AddRequest) (*pb.AddReply, error) {
	conn, err := grpc.Dial(req.GetAddress(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return &pb.AddReply{}, err
	}
	c := pb.NewServerClient(conn)

	s.conns = append(s.conns, conn)
	s.clients[req.GetAddress()] = c

	return &pb.AddReply{}, nil
}
