package service

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	pbStorage "karma/gen/storage"
)

type storageClient struct {
	id     string
	conn   grpc.ClientConnInterface
	client pbStorage.StorageClient
	cap    uint64
}

func (st *storageClient) saveSlice(ctx context.Context, path string, content []byte, n int) error {
	sliceLen := len(content) / 5

	if n < 4 {
		content = content[sliceLen*n : sliceLen*(n+1)]
	} else {
		content = content[sliceLen*n:]
	}

	req := &pbStorage.SaveRequest{
		Path:    path,
		Content: content,
	}

	reply, err := st.client.SaveSlice(ctx, req)
	if err != nil {
		return fmt.Errorf("error saving slice to %s: %s", st.id, err.Error())
	}

	if !reply.GetSuccess() {
		return fmt.Errorf("error saving slice to %s: %s", st.id, reply.GetMessage())
	}

	st.cap -= uint64(len(content))

	return nil
}

func (st *storageClient) loadSlice(ctx context.Context, path string) ([]byte, error) {
	req := &pbStorage.LoadRequest{Path: path}
	reply, err := st.client.LoadSlice(ctx, req)
	if err != nil {
		return nil, err
	}

	if !reply.GetSuccess() {
		return nil, fmt.Errorf("error loading slice from %s: %s", st.id, reply.GetMessage())
	}

	return reply.GetContent(), nil
}

type storageIterator struct {
	storages []*storageClient
	i, cnt   int
}

func (it *storageIterator) hasNext() bool {
	return it.cnt < len(it.storages)
}

func (it *storageIterator) next() *storageClient {
	st := it.storages[it.i]

	it.cnt++
	if it.i < len(it.storages)-1 {
		it.i++
	} else {
		it.i = 0
	}

	return st
}
