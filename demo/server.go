package main

import (
	"chalurania/api"
	"chalurania/comet"
	"chalurania/comet/router"
	"chalurania/service/log"
)

// 测试自定义路由
type HiRouter struct {
	router.Router
}

func (h *HiRouter) Handle(r api.IRequest) {
	//log.Info.Println("Call router handler")
	log.Info.Println("Receive from client message id:", r.GetMsgID(), " data:", string(r.GetData()))

	// 反向客户端发送数据
	err := r.GetConnection().SendMsg(1, 2, 101, []byte("Hi data pack1"))
	if err != nil {
		log.Error.Println("Send message to client err:", err)
	}
}

// 测试多路由
type HelloRouter struct {
	router.Router
}

func (h *HelloRouter) Handle(r api.IRequest) {
	log.Info.Println("receive from client message id:", r.GetMsgID(), " data:", string(r.GetData()))

	// 反向客户端发送数据
	err := r.GetConnection().SendMsg(1, 2, 201, []byte("Hello data pack2"))
	if err != nil {
		log.Error.Println("send message to client err:", err)
	}
}

// 创建连接时执行
func OnConnectionStart(conn api.IConnection) {
	log.Info.Println("on connection start called...")

	// 设置属性
	conn.SetProperty("name", "ellery")

	err := conn.SendMsg(1, 2, 301, []byte("Connect success"))
	if err != nil {
		log.Error.Println("on conn start err:", err)
	}
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
	s.AddRouter(1, &HiRouter{})
	s.AddRouter(2, &HelloRouter{})

	// 开启服务
	s.Serve()
}