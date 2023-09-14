package service

import (
	"errors"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

//管理用户账号登陆注册，好友添加等操作

// 用户注册。获取用户名、密码，返回用户id
func register(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	//注册
	uid, err := registerInDB(username, password)
	if err != nil {
		respError(c, 500, err)
		return
	}

	respOK(c, uid)
}

func login(c *gin.Context) {
	uid := c.PostForm("uid")
	password := c.PostForm("password")

	isSuccess, err := loginInDB(uid, password)
	if err != nil {
		respError(c, 500, err)
		return
	}

	if !isSuccess {
		respError(c, 401, errors.New("id或密码错误"))
		return
	}

	//生成token
	token, err := generateToken(uid)
	if err != nil {
		respError(c, 500, err)
		return
	}

	respOK(c, token)
}

func addFriend(c *gin.Context) {
	//鉴权
	uid, isExist := c.Get("user_id")
	if !isExist {
		respError(c, 401, errors.New("token无效"))
		return
	}

	//获取好友id
	friendUID := c.PostForm("friend_id")
	if friendUID == uid {
		respError(c, 401, errors.New("不能添加自己为好友"))
		return
	}

	//比较id大小
	idA, idB, err := compareID(uid.(string), friendUID)
	if err != nil {
		respError(c, 500, err)
		return
	}

	//获取username
	usernameA, err := getUsername(idA)
	if err != nil {
		respError(c, 500, err)
		return
	}
	usernameB, err := getUsername(idB)
	if err != nil {
		respError(c, 500, err)
		return
	}

	friend := Friend{
		Id:        idA + idB,
		UidA:      uid.(string),
		UidB:      friendUID,
		UsernameA: usernameA,
		UsernameB: usernameB,
		Agree:     0,
		Counter:   0,
	}

	//添加好友
	friend, err = addFriendInDB(friend)
	if err != nil {
		respError(c, 500, err)
		return
	}

	respOK(c, friend)
}

func deleteFriend(c *gin.Context) {
	//鉴权
	uid, isExist := c.Get("user_id")
	if !isExist {
		respError(c, 401, errors.New("token无效"))
		return
	}

	//获取好友id
	friendUID := c.PostForm("friend_id")
	if friendUID == uid {
		respError(c, 401, errors.New("不能添加自己为好友"))
		return
	}

	//比较id大小
	idA, idB, err := compareID(uid.(string), friendUID)
	if err != nil {
		respError(c, 500, err)
		return
	}

	//获取username
	usernameA, err := getUsername(idA)
	if err != nil {
		respError(c, 500, err)
		return
	}
	usernameB, err := getUsername(idB)
	if err != nil {
		respError(c, 500, err)
		return
	}

	friend := Friend{
		Id:        idA + idB,
		UidA:      uid.(string),
		UidB:      friendUID,
		UsernameA: usernameA,
		UsernameB: usernameB,
		Agree:     0,
		Counter:   0,
	}
	err = deleteFriendInDB(friend)
	if err != nil {
		respError(c, 500, err)
		return
	}

	respOK(c, nil)
}

//生成用户token
//规定：token不专门设置过期时间，生命周期与客户端生命周期相同

var secret = []byte("114514")

func generateToken(userID string) (string, error) {
	//生成token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
	})
	tokenString, err := token.SignedString(secret)
	return tokenString, err
}

// 比较id大小
func compareID(idA, idB string) (string, string, error) {
	//将id转化为int
	idAInt, err := strconv.Atoi(idA)
	if err != nil {
		return "", "", err
	}

	idBInt, err := strconv.Atoi(idB)
	if err != nil {
		return "", "", err
	}
	if idAInt > idBInt {
		return idB, idA, nil
	}
	return idA, idB, nil
}
