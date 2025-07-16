package main

import (
	"log"
	"github.com/pc/filemgr/internal/app"
	"github.com/gin-gonic/gin"
)

func main() {
	app, err := app.Init()
	if err != nil {
		log.Fatal("App init failed: ", err)
	}

	r := gin.Default()
	r.POST("/file/uploadUrl", app.FileHandler.GenerateUploadURL)
	r.POST("/file/downloadUrl", app.FileHandler.GenerateDownloadURL)
	log.Fatal(r.Run(":8080"))
}
