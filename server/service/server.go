package service

import (
	"context"
	pb "karma/gen/server"
)

type Service struct {
	pb.UnimplementedServerServer
	addr []string
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) AddStorage(ctx context.Context, req *pb.AddRequest) (*pb.AddReply, error) {
	s.addr = append(s.addr, req.GetAddress())
	return &pb.AddReply{}, nil
}
