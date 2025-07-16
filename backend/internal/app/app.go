package app

import (
	"github.com/pc/filemgr/internal/handler"
	"github.com/pc/filemgr/internal/infra/storage"
	"github.com/pc/filemgr/internal/service"
)

type Application struct {
	FileHandler *handler.FileHandler
}

func Init() (*Application, error) {
	// TODO: get config from env
	storageClient, err := storage.NewMinioClient(
		"minio:9000",
		"filemgr",
		"password@123456",
	)
	if err != nil {
		return nil, err
	}

	fileService := service.NewFileService(storageClient)
	fileHandler := handler.NewFileHandler(fileService)

	return &Application{
		FileHandler: fileHandler,
	}, nil
}
