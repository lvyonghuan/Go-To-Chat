package service

import "Go-To-Chat/service/saver/service"

func main() {
	service.InitRedis()
	service.InitRouter()
}
