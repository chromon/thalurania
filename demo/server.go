package main

import (
	"chalurania/api"
	"chalurania/comet"
	"chalurania/comet/caller/routers"
	"chalurania/comet/constants"
	"chalurania/service/log"
)

// 创建连接时执行
func OnConnectionStart(conn api.IConnection) {
	log.Info.Println("on connection start called...")

	// 设置属性
	conn.SetProperty("name", "ellery")
}

// 断开连接时执行
func OnConnectionLost(conn api.IConnection) {
	log.Info.Println("on connection lost called...")

	// 获取属性
	if name, err := conn.GetProperty("name"); err == nil {
		log.Info.Println("conn property name:", name)
	}
}

func main() {
	// 创建 server
	s := comet.NewServer()

	// 注册连接回调函数
	s.SetOnConnStart(OnConnectionStart)
	s.SetOnConnStop(OnConnectionLost)

	// 添加自定义路由
	s.AddRouter(constants.SignUpOption, &routers.RegisterRouter{})
	s.AddRouter(constants.LoginOption, &routers.LoginRouter{})
	s.AddRouter(constants.LogoutOption, &routers.LogoutRouter{})

	// 开启服务
	s.Serve()
}