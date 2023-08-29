package main

import "Go-To-Chat/service/connect/service"

func main() {
	service.InitDataBase()
	service.InitRouter()
}
