package service

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "karma/gen/server"
	pbStorage "karma/gen/storage"
)

type Storages interface {
	AddStorage(st *storageClient)
	SaveFile(ctx context.Context, path string, content []byte) error
	LoadFile(ctx context.Context, path string) ([]byte, error)
}

type Service struct {
	pb.UnimplementedServerServer
	storages Storages
}

func NewService(storages Storages) *Service {
	return &Service{
		storages: storages,
	}
}

func (s *Service) AddStorage(_ context.Context, req *pb.AddRequest) (*pb.AddReply, error) {
	conn, err := grpc.Dial(req.GetAddress(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return &pb.AddReply{}, err
	}
	c := pbStorage.NewStorageClient(conn)

	st := &storageClient{
		id:     req.GetAddress(),
		conn:   conn,
		client: c,
		cap:    req.GetCapacity(),
	}

	s.storages.AddStorage(st)

	return &pb.AddReply{}, nil
}
