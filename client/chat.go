package client

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

//两个ws通信系统

var upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }, //允许跨域
}

// 客户端前端连接客户端后端
func clientFrontToBack(c *gin.Context) {

}

// 客户端后端连接服务器
func clientToService() {

}
