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
	log.Info.Println("Call router handler")
	log.Info.Println("Receive from client message id:", r.GetMsgID(), " data:", string(r.GetData()))

	// 反向客户端发送数据
	err := r.GetConnection().SendMsg(1, []byte("Hi data pack"))
	if err != nil {
		log.Error.Println("Send message to client err:", err)
	}
}

func main() {
	// 创建 server
	s := comet.NewServer()

	// 添加自定义路由
	s.AddRouter(&HiRouter{})

	// 开启服务
	s.Serve()
}