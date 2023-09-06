package service

import "github.com/gin-gonic/gin"

//向客户端返回信息

func respOK(c *gin.Context, info interface{}) {
	c.JSON(200, gin.H{
		"code": 200,
		"info": info,
	})
}

func respError(c *gin.Context, code int, msg error) {
	c.JSON(200, gin.H{
		"code": code,
		"msg":  msg.Error(),
	})
}
