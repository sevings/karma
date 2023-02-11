package service

import (
	"context"
	"github.com/spf13/afero"
	pb "karma/gen/storage"
	"os"
)

type Service struct {
	pb.UnimplementedStorageServer
	fs afero.Fs
}

func NewService(fs afero.Fs) *Service {
	return &Service{
		fs: fs,
	}
}

func (s *Service) SaveSlice(ctx context.Context, req *pb.SaveRequest) (*pb.SaveReply, error) {
	err := afero.WriteFile(s.fs, req.GetPath(), req.GetContent(), os.FileMode(0666))
	if err != nil {
		return &pb.SaveReply{Success: false, Message: err.Error()}, nil
	}

	return &pb.SaveReply{Success: true}, nil
}

func (s *Service) LoadSlice(ctx context.Context, req *pb.LoadRequest) (*pb.LoadReply, error) {
	content, err := afero.ReadFile(s.fs, req.GetPath())
	if err != nil {
		return &pb.LoadReply{Success: false, Message: err.Error()}, nil
	}

	return &pb.LoadReply{Success: true, Content: content}, nil
}
