package main

import (
	"bytes"
	"context"
	"karma/client/client"
	"karma/client/client/operations"
	"log"
)

func main() {
	karma := client.Default

	content := make([]byte, 4096)
	for i := range content {
		content[i] = byte(i % 255)
	}

	err := put(karma, content)
	if err != nil {
		return
	}

	loaded, err := get(karma)
	if err != nil {
		return
	}

	if len(loaded) != len(content) {
		log.Println("Loaded length: ", len(loaded))
		return
	}

	for i := range loaded {
		if loaded[i] != content[i] {
			log.Printf("Loaded %d, but content was %d at %d", loaded[i], content[i], i)
			break
		}
	}

	log.Println("done")
}

type testFile struct {
	reader *bytes.Reader
}

func (tf testFile) Read(b []byte) (n int, err error) {
	return tf.reader.Read(b)
}

func (testFile) Close() error {
	return nil
}

func (testFile) Name() string {
	return "test_name"
}

func put(karma *client.Karma, content []byte) error {
	file := testFile{reader: bytes.NewReader(content)}
	params := &operations.PutFileParams{
		File:    file,
		Path:    "/test/path",
		Context: context.Background(),
	}

	_, err := karma.Operations.PutFile(params)
	if err != nil {
		log.Println("put:", err)
		return err
	}

	log.Println("put success")

	return nil
}

func get(karma *client.Karma) ([]byte, error) {
	content := &bytes.Buffer{}
	params := &operations.GetFileParams{
		Path:    "/test/path",
		Context: context.Background(),
	}

	_, err := karma.Operations.GetFile(params, content)
	if err != nil {
		log.Println("get:", err)
		return nil, err
	}

	log.Println("get success")

	return content.Bytes(), nil
}
