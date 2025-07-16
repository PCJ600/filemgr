package service

import (
	"context"
	"time"
	"github.com/pc/filemgr/internal/infra/storage"
)

type FileService struct {
	storage *storage.Client
}

func NewFileService(storage *storage.Client) *FileService {
	return &FileService{storage: storage}
}

func (s *FileService) GenerateUploadURL(ctx context.Context, bucketName string, 
										objectKey string, expirySec int64) (string, error) {
	expiry := time.Duration(expirySec) * time.Second
	return s.storage.PresignedPutObject(ctx, bucketName, objectKey, expiry)
}


