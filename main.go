package main

import (
	"simple-demo/data"
	"simple-demo/service"
	"github.com/gin-gonic/gin"
)

func main() {
	// 启动消息服务，使用go关键字创建一个新的goroutine运行service.RunMessageServer()函数，
	// 这样可以在后台并发运行消息服务，不会阻塞主程序的执行。
	go service.RunMessageServer()

	//初始化数据库
	data.InitDb()

	// 创建一个默认的gin.Engine实例，用于构建HTTP路由。
	r := gin.Default()

	// 初始化路由，调用initRouter函数，用于设置API路由及其处理函数。
	initRouter(r)

	// 启动HTTP服务器并监听端口，使用r.Run()方法，默认监听0.0.0.0:8080地址。
	r.Run()
}
