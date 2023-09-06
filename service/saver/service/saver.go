package service

import (
	"hash/crc32"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

//在服务器上存储文字聊天记录，并且负责维护存储的聊天记录。同时负责图片等的存储。
//分类别进行存储。考虑到采用ws进行cs通信，那么文字存储在文本领域，视频图片等分别存储。ws时文字直接发送，图片等发送url连接进行http的request。

func storeFile(c *gin.Context) {
	file := c.PostForm("file")
	fileName := c.PostForm("fileName")

	//获取后缀
	suffix := getSuffix(fileName)

	//根据文件数据计算文件哈希值
	hash := hashValue([]byte(file))

	//存储文件到files目录下
	filePath := "./files/" + hash + suffix
	err := os.WriteFile(filePath, []byte(file), 0666)
	if err != nil {
		respError(c, 500, err)
		return
	}

	//将文件名-hash键值对存储在redis中
	err = storeFileNameAndHash(fileName, hash)
	if err != nil {
		respError(c, 500, err)
		return
	}

	respOK(c, "图片存储成功")
}

// 获取后缀，返回后缀内容。如果没有后缀则返回空字符串。
func getSuffix(fileName string) string {
	for i := len(fileName) - 1; i >= 0; i-- {
		if fileName[i] == '.' {
			return fileName[i:]
		}
	}
	return ""
}

// 计算文件哈希值
func hashValue(file []byte) string {
	return strconv.Itoa(int(crc32.ChecksumIEEE(file)))
}
