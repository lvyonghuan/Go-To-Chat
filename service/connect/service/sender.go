package service

import (
	"log"

	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
)

// 发送消息结构体
type send struct {
	RoomID       string `json:"room_id"`
	SenderID     string `json:"sender_id"`
	SerialNumber int    `json:"serial_number"`
	MessageType  int    `json:"message_type"`
	Message      string `json:"message"`
}

func (m *receive) sendMessage() {
	//获取房间
	room := rooms[m.RoomID]

	//将消息转化为发送格式
	var message = send{
		RoomID:       m.RoomID,
		SenderID:     m.SenderID,
		SerialNumber: m.SerialNumber,
		MessageType:  m.MessageType,
		Message:      m.Message,
	}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	messageByte, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
		return
	}

	//发送消息
	for _, v := range room.clients {
		err := v.WriteMessage(websocket.TextMessage, messageByte)
		if err != nil {
			log.Println(err)
			return
		}
	}
}
