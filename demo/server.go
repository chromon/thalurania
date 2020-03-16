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

func (h *HiRouter) PreHandle(r api.IRequest) {
	log.Info.Println("Call router preHandler")
	_, err := r.GetConnection().GetTCPConnection().Write([]byte("Before hi"))
	if err != nil {
		log.Error.Println("Call router preHandler err:", err)
	}
}

func (h *HiRouter) Handle(r api.IRequest) {
	log.Info.Println("Call router handler")
	_, err := r.GetConnection().GetTCPConnection().Write([]byte("HiHiHi"))
	if err != nil {
		log.Error.Println("Call router handler err:", err)
	}
}

func (h *HiRouter) PostHandle(r api.IRequest) {
	log.Info.Println("Call router postHandler")
	_, err := r.GetConnection().GetTCPConnection().Write([]byte("After hi"))
	if err != nil {
		log.Error.Println("Call router postHandler err:", err)
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