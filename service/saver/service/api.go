package service

import "github.com/gin-gonic/gin"

func InitRouter() {
	r := gin.Default()

	r.POST("/storeFile", storeFile)

	r.GET("/getFile")

	r.Run(":8081")
}
