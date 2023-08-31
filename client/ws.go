package client

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func connect(id string) {
	wsURL := "ws://127.0.0.1:8080/" + id
	header := http.Header{}
	//TODO:鉴权
	//header.Add("Authorization",)

	conn, _, err := websocket.DefaultDialer.Dial(wsURL, header)
	if err != nil {
		log.Println("客户端连接错误：", err)
		return
	}

	go send(conn)
	go receive(conn)
}

// 发送消息
func send(conn *websocket.Conn) {
	defer conn.Close()

	for {
		select {
		//TODO：接收消息
		}
	}
}

// 接收消息
func receive(conn *websocket.Conn) {
	for {
		typ, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("接收消息错误：", err)
			return
		}

		switch typ {
		case websocket.TextMessage:
			//TODO：处理消息
			log.Println(string(message))
		}
	}
}
