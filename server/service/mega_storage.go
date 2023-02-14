package service

import (
	"context"
	"fmt"
	"log"
)

type PathStore interface {
	SetPath(path string, size int, ids []string)
	GetPath(path string) (int, []string)
}

type MegaStorage struct {
	paths      PathStore
	idStorages map[string]*storageClient
	storages   []*storageClient
	nextI      int
}

func NewMegaStorage(paths PathStore) *MegaStorage {
	return &MegaStorage{
		paths:      paths,
		idStorages: make(map[string]*storageClient),
	}
}

func (ms *MegaStorage) AddStorage(st *storageClient) {
	ms.idStorages[st.id] = st
	ms.storages = append(ms.storages, st)
}

func (ms *MegaStorage) SaveFile(ctx context.Context, path string, content []byte) error {
	log.Printf("Save %d %s", len(content), path)

	storages := make([]*storageClient, 0, 5)
	lens := make([]int, 0, 5)

	sliceLen := len(content) / 5
	for i := 0; i < 4; i++ {
		lens = append(lens, sliceLen)
	}
	lens = append(lens, len(content)-sliceLen*4)

	it := storageIterator{storages: ms.storages, i: ms.nextI}
	for it.hasNext() {
		if len(storages) == 5 {
			break
		}

		st := it.next()
		if st.cap < uint64(lens[len(storages)]) {
			continue
		}

		storages = append(storages, st)
	}

	if len(storages) < 5 {
		return fmt.Errorf("not enough space")
	}

	var ids []string

	for i := 0; i < 5; i++ {
		err := storages[i].saveSlice(ctx, path, content, i)
		if err != nil {
			return err
		}

		ids = append(ids, storages[i].id)
	}

	ms.paths.SetPath(path, len(content), ids)

	return nil
}

func (ms *MegaStorage) LoadFile(ctx context.Context, path string) ([]byte, error) {
	log.Printf("Load %s", path)

	size, ids := ms.paths.GetPath(path)
	if len(ids) == 0 {
		return nil, fmt.Errorf("path not found")
	}

	var slices [][]byte

	for _, id := range ids {
		st := ms.idStorages[id]
		if st == nil {
			return nil, fmt.Errorf("internal error")
		}

		fileSlice, err := st.loadSlice(ctx, path)
		if err != nil {
			return nil, err
		}

		slices = append(slices, fileSlice)
	}

	content := make([]byte, size)

	var i int
	for _, s := range slices {
		i += copy(content[i:], s)
	}

	return content, nil
}
