package service

import (
	"bytes"
	"context"
	"github.com/go-openapi/runtime/middleware"
	"io"
	"karma/server/restapi/operations"
	"log"
	"strings"
)

func NewDownloadHandler(st Storages) func(operations.GetFileParams) middleware.Responder {
	return func(params operations.GetFileParams) middleware.Responder {
		content, err := st.LoadFile(context.Background(), params.Path)
		if err == nil {
			payload := io.NopCloser(bytes.NewReader(content))
			return operations.NewGetFileOK().WithPayload(payload)
		}

		if strings.Contains(err.Error(), "not found") {
			return operations.NewGetFileNotFound()
		}

		return operations.NewGetFileInternalServerError()
	}
}

func NewUploadHandler(st Storages) func(operations.PutFileParams) middleware.Responder {
	return func(params operations.PutFileParams) middleware.Responder {
		content, err := io.ReadAll(params.File)
		if err != nil {
			log.Println(err)
			return operations.NewPutFileInternalServerError()
		}

		err = st.SaveFile(context.Background(), params.Path, content)
		if err != nil {
			log.Println(err)
			return operations.NewPutFileInternalServerError()
		}

		return operations.NewPutFileOK()
	}
}
