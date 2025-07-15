package main

import (
  "github.com/gin-gonic/gin"
  "github.com/pc/filemgr/handlers"
)

func main() {
  r := gin.Default()

  r.POST("/fileUpload/uploadUrl", handlers.GenerateUploadUrl)

  r.Run(":8080")
}
