package service

import (
	"log"

	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
)

// 接收消息模型
type receive struct {
	//消息所在房间房间号
	RoomID string `json:"room_id"`
	//消息发送者ID
	SenderID string `json:"sender_id"`
	//消息序列号
	SerialNumber int `json:"serial_number"`
	//消息类型
	MessageType int `json:"message_type"`
	//消息内容
	Message string `json:"message"`

	//内部消息
	//消息发送时间
	time int64
}

// 接收websocket消息
func (r *room) read(conn *websocket.Conn) {
	defer func() {
		r.mu.Lock()
		r.connectCount--
		r.mu.Unlock()
		conn.Close()
		if r.connectCount == 0 {
			//将房间的对话量更新到数据库
			r.saveMessageCount()
			//房间无人连接，删除房间
			delete(rooms, r.id)
		}
	}()

	//获取消息
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			break
		}

		//处理消息
		go handelMessage(messageType, message)
	}
}

// 处理消息
func handelMessage(messageType int, message []byte) {
	switch messageType {
	case websocket.TextMessage:
		//json格式化
		var json = jsoniter.ConfigCompatibleWithStandardLibrary
		var receive receive
		if err := json.Unmarshal(message, &receive); err != nil {
			log.Println(receive)
			log.Println(err)
			return
		}
		//核对序列号
		receive.checkNumber()
		//TODO：云端存储（并发执行？）

		//TODO：频道广播
		go receive.sendMessage()
	}
}

// 序列号核对
func (m *receive) checkNumber() {
	roomID := m.RoomID
	r, ok := rooms[roomID]
	if !ok {
		log.Println("？？？")
		return
	}
	if r.counter > m.SerialNumber {
		m.SerialNumber = r.counter
		//TODO：发送同步消息
	}

	//序列号加1
	r.mu.Lock()
	r.counter++
	r.mu.Unlock()
}
