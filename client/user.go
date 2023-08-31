package client

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

var token string //懒了

//管理部分用户行为，例如：登陆、注册、添加好友
//聊天由chat.go进行处理
//服务器必须部署在公网的固定IP上
//数据存放于data目录下

// 注册
func register(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	//注册
	err, uid := registerService(username, password)
	if err != nil {
		respError(c, 500, err)
		return
	}

	respOK(c, "注册成功，账号id为："+uid)
}

// 登陆
func login(c *gin.Context) {
	uid := c.PostForm("uid")
	password := c.PostForm("password")

	err := loginService(uid, password)
	if err != nil {
		respError(c, 500, err)
		return
	}

	respOK(c, "登陆成功")
}

// 添加好友
func addFriend(c *gin.Context) {
	friendUID := c.PostForm("friend_id")

	err, info := addFriendService(friendUID)
	if err != nil {
		respError(c, 500, err)
		return
	}

	respOK(c, info)
}

//service服务

// 注册
func registerService(username, password string) (error, string) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	writer.WriteField("username", username)
	writer.WriteField("password", password)
	writer.Close()

	//发送注册请求
	client, _ := http.NewRequest("POST", url+"register", body)
	client.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(client)
	if err != nil {
		return err, ""
	}
	defer resp.Body.Close()

	//获取返回值
	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return err, ""
	}
	var (
		r    responseMessage
		json = jsoniter.ConfigCompatibleWithStandardLibrary
	)
	err = json.Unmarshal(response, &r)
	if err != nil {
		return err, ""
	}

	if r.Code != 200 {
		return errors.New(r.Message), ""
	}

	return nil, r.Info.(string)
}

func loginService(uid, password string) error {
	body := new(bytes.Buffer)

	writer := multipart.NewWriter(body)
	writer.WriteField("uid", uid)
	writer.WriteField("password", password)
	writer.Close()

	client, _ := http.NewRequest("POST", url+"login", body)
	client.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(client)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var (
		r    responseMessage
		json = jsoniter.ConfigCompatibleWithStandardLibrary
	)
	err = json.Unmarshal(response, &r)
	if err != nil {
		return err
	}
	if r.Code != 200 || r.Info == nil {
		return errors.New(r.Message)
	}

	token = r.Info.(string)

	return nil
}

// 添加好友
func addFriendService(friendID string) (error, any) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("friend_id", friendID)
	writer.Close()

	client, _ := http.NewRequest("POST", url+"add_friend", body)
	client.Header.Set("Content-Type", writer.FormDataContentType())
	client.Header.Set("Authorization", token)

	resp, err := http.DefaultClient.Do(client)
	if err != nil {
		return err, nil
	}

	defer resp.Body.Close()

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return err, nil
	}

	var (
		r    responseMessage
		json = jsoniter.ConfigCompatibleWithStandardLibrary
	)
	err = json.Unmarshal(response, &r)
	if err != nil {
		return err, nil
	}
	if r.Code != 200 || r.Info == nil {
		return errors.New(r.Message), nil
	}
	//TODO：存储好友数据

	return nil, r.Info
}
