package service

import (
	"context"
	"strings"
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

func (s *FileService) GenerateDownloadURL(ctx context.Context, bucketName string, objectKey string,
										  expireSeconds int64) (string, error) {
	_, err := s.storage.StatObject(ctx, bucketName, objectKey)
	if err != nil {
		return "", err
	}

	expiry := time.Duration(expireSeconds) * time.Second
	return s.storage.PresignedGetObject(ctx, bucketName, objectKey, expiry)
}

func (s *FileService) DeleteObject(ctx context.Context, bucketName string, objectKey string) error {
	if strings.HasSuffix(objectKey, "/") {
		return s.storage.DeleteObjectsWithPrefix(ctx, bucketName, objectKey)
	}
	return s.storage.DeleteObject(ctx, bucketName, objectKey)
}
