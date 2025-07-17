package storage

import (
	"context"
	"fmt"
	"time"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Client struct {
	cli *minio.Client
}

func NewStorageClient(endpoint string , accessKey string, secretKey string) (*Client, error) {
	cli, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}
	return &Client{cli: cli}, nil
}

func (c *Client) StatObject(ctx context.Context, bucketName string, objectKey string) (minio.ObjectInfo, error) {
	return c.cli.StatObject(ctx, bucketName, objectKey, minio.StatObjectOptions{})
}

func (c *Client) PresignedPutObject(ctx context.Context, bucketName string, objectKey string,
									expiry time.Duration) (string, error) {
	url, err := c.cli.PresignedPutObject(ctx, bucketName, objectKey, expiry)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}

func (c *Client) PresignedGetObject(ctx context.Context, bucketName string, objectKey string,
									expiry time.Duration) (string, error) {
	url, err := c.cli.PresignedGetObject(ctx, bucketName, objectKey, expiry, nil)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}

func (c *Client) DeleteObject(ctx context.Context, bucketName string, objectKey string) error {
	return c.cli.RemoveObject(ctx, bucketName, objectKey, minio.RemoveObjectOptions{})
}

func (c *Client) DeleteObjectsWithPrefix(ctx context.Context, bucketName string, prefix string) error {
	objectsCh := make(chan minio.ObjectInfo)

	go func() {
		defer close(objectsCh)
		for obj := range c.cli.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
			Prefix:    prefix,
			Recursive: true,
		}) {
			objectsCh <- obj
		}
	}()

	for err := range c.cli.RemoveObjects(ctx, bucketName, objectsCh, minio.RemoveObjectsOptions{}) {
		return fmt.Errorf("delete failed: %v", err)
	}
	return nil
}

