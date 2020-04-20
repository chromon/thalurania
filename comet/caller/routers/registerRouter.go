package routers

import (
	"chalurania/api"
	"chalurania/comet/packet"
	"chalurania/comet/router"
	"chalurania/comet/variable"
	"chalurania/service/log"
	"encoding/json"
)

// 注册 router
type RegisterRouter struct {
	router.Router
}

// 注册处理
func (rr *RegisterRouter) Handle(r api.IRequest) {
	//log.Info.Println("received from client message id:", r.GetMsgID(), " data:", string(r.GetData()))

	// 将注册信息包装并序列化
	dw := packet.NewDataWrap(1, r.GetData())
	ret, err := json.Marshal(dw)
	if err != nil {
		log.Info.Println("serialize register data wrap object err:", err)
		return
	}

	// 将序列化后的信息发布到异步存储管道
	_, err = variable.RedisPool.Publish("AsyncPersistence", string(ret))
	if err != nil {
		log.Error.Println("redis pool publish to async persistence err:", err)
	}
}

// 回执消息
func (rr *RegisterRouter) PostHandle(r api.IRequest) {
	// 反向客户端发送 ack 数据
	err := r.GetConnection().SendMsg(1, 2, 101, []byte("register an account success"))
	if err != nil {
		log.Error.Println("register send ack message to client err:", err)
	}
}
