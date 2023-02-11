package service

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "karma/gen/storage"
)

type Service struct {
	pb.UnimplementedStorageServer
}

func (s *Service) SaveSlice(context.Context, *pb.SaveRequest) (*pb.SaveReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SaveSlice not implemented")
}

func (s *Service) LoadSlice(context.Context, *pb.LoadRequest) (*pb.LoadReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoadSlice not implemented")
}
