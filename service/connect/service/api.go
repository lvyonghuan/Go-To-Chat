package service

import "github.com/gin-gonic/gin"

func InitRouter() {
	r := gin.Default()

	//登录注册
	r.POST("/register", register)
	r.GET("/login", login)

	r.POST("/addFriend", checkToken(), addFriend)

	//建立ws连接
	r.GET("/ws/:id", checkToken(), connect) //懒得想了，一个房间的id为两个id的拼接字符串。低位在前高位拼后。

	//资源获取
	//听说telegram的数据存储是四级降级的，内存一块，固态一块，机械硬盘一块，磁带机一块。只是顺带一提。
	//资源获取移步资源服务器。这里只处理客户端连接的逻辑。

	r.Run(":8080")
}
