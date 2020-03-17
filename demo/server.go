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
	err := r.GetConnection().SendMsg(1, []byte("Hi data pack1"))
	if err != nil {
		log.Error.Println("Send message to client err:", err)
	}
}

// 测试多路由
type HelloRouter struct {
	router.Router
}

func (h *HelloRouter) Handle(r api.IRequest) {
	log.Info.Println("Receive from client message id:", r.GetMsgID(), " data:", string(r.GetData()))

	// 反向客户端发送数据
	err := r.GetConnection().SendMsg(2, []byte("Hello data pack2"))
	if err != nil {
		log.Error.Println("Send message to client err:", err)
	}
}

func main() {
	// 创建 server
	s := comet.NewServer()

	// 添加自定义路由
	s.AddRouter(1, &HiRouter{})
	s.AddRouter(2, &HelloRouter{})

	// 开启服务
	s.Serve()
}