package handlers

import (
	"context"
	"fmt"
	"log"
	// "os"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/gin-gonic/gin"
	"net/http"
)

type presignedUrlRequest struct {
	BucketName string `json:"bucketName" binding:"required"`
	FileName string `json:"fileName" binding:"required"`
	TokenDurationSeconds int64 `json:"tokenDurationSeconds" binding:"required"`
}


func GenerateUploadUrl(c *gin.Context) {
	endpoint := "minio:9000"
	// accessKey := os.Getenv("MINIO_ACCESS_KEY")
	// secretKey := os.Getenv("MINIO_SECRET_KEY")
	accessKey := "myminioadmin"
	secretKey := "password@123456"

	var reqData presignedUrlRequest
	if err := c.ShouldBindJSON(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 4000001,
			"message": fmt.Sprintf("request param invalid: %v", err),
		})
		return
	}
	log.Printf("req data: %+v", reqData)

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Printf("MinIO initial failed: %v", err)
        c.JSON(http.StatusServiceUnavailable, gin.H{
			"code": 5030001,
			"message": "MinIO initial failed",
        })
        return
	}

	expiry := time.Duration(reqData.TokenDurationSeconds) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()
	presignedUrl, err := client.PresignedPutObject(ctx, reqData.BucketName, reqData.FileName, expiry)
	if err != nil {
		log.Printf("generate presigned url failed: %v", err)
        c.JSON(http.StatusServiceUnavailable, gin.H{
			"code": 5030002,
			"message": "generate presigned url failed",
        })
        return
	}

	resp := gin.H{
		"presignedUrl": presignedUrl.String(),
	}
	c.PureJSON(http.StatusOK, resp)
}
