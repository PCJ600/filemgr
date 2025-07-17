package handler

import (
	"net/http"
	"context"
	"log"
	"strings"
	"time"
	"github.com/pc/filemgr/internal/service"
	"github.com/gin-gonic/gin"
)

type FileHandler struct {
	fileService *service.FileService
}

func NewFileHandler(fileService *service.FileService) *FileHandler {
	return &FileHandler{
		fileService: fileService,
	}
}

func (h *FileHandler) GenerateUploadURL(c *gin.Context) {
	var req struct {
		BucketName string `json:"bucketName" binding:"required,min=3,max=63"`
		ObjectKey string `json:"objectKey" binding:"required,min=1,max=1024"`
		ExpireSeconds int64 `json:"expireSeconds" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("req data: %+v", req)
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 4000001,
			"message": "req body invalid: " + err.Error(),
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10 * time.Second)
    defer cancel()
	url, err := h.fileService.GenerateUploadURL(ctx, req.BucketName, req.ObjectKey, req.ExpireSeconds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 5000001,
			"message": "generate upload url failed: " + err.Error(),
		})
		return
	}

	c.PureJSON(http.StatusOK, gin.H {"code": 0, "url": url})
}

func (h *FileHandler) GenerateDownloadURL(c *gin.Context) {
	var req struct {
		BucketName string `json:"bucketName" binding:"required,min=3,max=63"`
		ObjectKey string `json:"objectKey" binding:"required,min=1,max=1024"`
		ExpireSeconds int64 `json:"expireSeconds" binding:"required"`
	}

    if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("req data: %+v", req)
        c.JSON(http.StatusBadRequest, gin.H{
            "code": 4000001,
            "msg":  "req body invalid: " + err.Error(),
        })
        return
    }

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10 * time.Second)
    defer cancel()
    url, err := h.fileService.GenerateDownloadURL(ctx, req.BucketName, req.ObjectKey, req.ExpireSeconds)
    if err != nil {
		if strings.Contains(err.Error(), "The specified key does not exist") {
			c.JSON(http.StatusNotFound, gin.H{
				"code": 4040001,
				"msg":  "file not found",
			})
		} else {
			c.JSON(http.StatusNotFound, gin.H{
				"code": 5000001,
				"msg":  "internel server err: " + err.Error(),
			})
		}
        return
    }

	c.PureJSON(http.StatusOK, gin.H {"code": 0, "url": url})
}

func (h *FileHandler) DeleteObject(c *gin.Context) {
	var req struct {
		BucketName string `json:"bucketName" binding:"required,min=3,max=63"`
		ObjectKey string `json:"objectKey" binding:"required,min=1,max=1024"`
	}

    if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("req data: %+v", req)
        c.JSON(http.StatusBadRequest, gin.H{
            "code": 4000001,
            "msg":  "req body invalid: " + err.Error(),
        })
        return
    }

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10 * time.Second)
    defer cancel()

	if err := h.fileService.DeleteObject(ctx, req.BucketName, req.ObjectKey); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 5000001,
			"msg":  "delete object failed: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0})
}
