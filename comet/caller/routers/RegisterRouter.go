package routers

import (
	"chalurania/api"
	"chalurania/comet/router"
	"chalurania/service/log"
)

// 注册 router
type RegisterRouter struct {
	router.Router
}

// 注册处理
func (rr *RegisterRouter) Handle(r api.IRequest) {
	log.Info.Println("received from client message id:", r.GetMsgID(), " data:", string(r.GetData()))





}

// 回执消息
func (rr *RegisterRouter) PostHandle(r api.IRequest) {
	// 反向客户端发送 ack 数据
	err := r.GetConnection().SendMsg(1, 2, 101, []byte("register an account success"))
	if err != nil {
		log.Error.Println("register send ack message to client err:", err)
	}
}