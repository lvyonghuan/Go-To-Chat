package service

import (
	"log"
	"strconv"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	Id       int
	Username string
	Password string
}

// Friend 好友模型
type Friend struct {
	Id        string `json:"room_id"`
	UidA      string `json:"uid_a"`
	UidB      string `json:"uid_b"`
	UsernameA string `json:"username_a"`
	UsernameB string `json:"username_b"`
	Agree     int    //对方是否同意添加好友。0表示状态待确认，1表示同意，2表示拒绝。
	Counter   int    //服务器接收到一条消息后，先于counter对比，如果小于counter则先返回一条要求同步的消息。在对比完成后，计数器加1。
}

//数据库，与mysql连通

var db *gorm.DB

// InitDataBase 使用gorm初始化与mysql的对接
func InitDataBase() {
	dsn := "root:42424242@tcp(127.0.0.1:3306)/chat"
	db, _ = gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

// 注册，返回id和错误
func registerInDB(username, password string) (string, error) {
	var u = User{
		Username: username,
		Password: password,
	}
	err := db.Create(&u).Error
	return strconv.Itoa(u.Id), err
}

// 登陆，返回布尔（数据是否和库内数据对得上）和错误
func loginInDB(userID, password string) (bool, error) {
	//userid转化为int
	id, err := strconv.Atoi(userID)
	if err != nil {
		return false, err
	}

	var u = User{
		Id:       id,
		Password: password,
	}
	err = db.Where("id = ? AND Password = ?", id, password).First(&u).Error
	return err == nil, err
}

// 添加好友，返回好友和错误
func addFriendInDB(friend Friend) (Friend, error) {
	err := db.Create(&friend).Error
	return friend, err
}

// 删除好友
func deleteFriendInDB(friend Friend) error {
	err := db.Delete(&friend).Error
	return err
}

// 获取用户名
func getUsername(uid string) (string, error) {
	id, err := strconv.Atoi(uid)
	if err != nil {
		return "", err
	}

	var u = User{
		Id: id,
	}
	if err := db.Where("id = ?", id).First(&u).Error; err != nil {
		return "", err
	}
	return u.Username, nil
}

// 获取好友信息
func getFriendInfo(id string) (Friend, error) {
	var f = Friend{
		Id: id,
	}
	if err := db.Where("id = ?", id).First(&f).Error; err != nil {
		return f, err
	}
	return f, nil
}

// 存储对话数据
func (r *room) saveMessageCount() {
	var f = Friend{
		Counter: r.counter,
	}
	//更新数据
	err := db.Model(&f).Where("id = ?", r.id).Update("counter", r.counter).Error
	if err != nil {
		log.Println(err)
		return
	}
}
