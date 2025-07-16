package storage

import (
	"context"
	"time"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Client struct {
	cli *minio.Client
}

func NewMinioClient(endpoint string , accessKey string, secretKey string) (*Client, error) {
	cli, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}
	return &Client{cli: cli}, nil
}

func (c *Client) PresignedPutObject(ctx context.Context, bucketName string, objectKey string, expiry time.Duration) (string, error) {
	url, err := c.cli.PresignedPutObject(ctx, bucketName, objectKey, expiry)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}
