package client

import "github.com/gin-gonic/gin"

//？
//这不是前端，这是客户端的后端。（雾
//用户点击行为触发api

func initRouter() {
	r := gin.Default()

	r.POST("/register", register)
	r.GET("/login")
}
