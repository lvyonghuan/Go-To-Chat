package service

import (
	"bytes"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
)

const (
	textMessage = iota
	fileMessage
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
	//文件名
	FileName string `json:"file_name"`

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
		switch receive.MessageType {
		case textMessage:

		case fileMessage:
			err := storeFile(receive.Message, receive.FileName)
			if err != nil {
				//TODO:反馈错误
			}
		}

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

// 向存储服务器发送文件
func storeFile(file, fileName string) error {
	//将file与filename加入到form-data中
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	writer.WriteField("file", file)
	writer.WriteField("fileName", fileName)
	writer.Close()

	req, err := http.NewRequest("POST", "http://localhost:8081/storeFile", body)
	if err != nil {
		log.Println("创建请求失败:", err)
		return err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("发送请求失败:", err)
		return err
	}

	defer resp.Body.Close()

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return err
	}
	var response responseMessage
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	err = json.Unmarshal(respData, &response)
	if err != nil {
		log.Println(err)
		return err
	}

	if response.Code != 200 {
		log.Println(response.Message)
		return errors.New(response.Message)
	}

	return nil
}
