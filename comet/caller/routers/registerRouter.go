package routers

import (
	"chalurania/api"
	"chalurania/comet/constants"
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
	//log.Info.Println("received from client register message id:", r.GetMsgID(), " data:", string(r.GetData()))

	// 将注册信息包装并序列化
	dw := packet.NewDataPersistWrap(constants.SignUpPersistenceOpt, r.GetData())
	ret, err := json.Marshal(dw)
	if err != nil {
		log.Info.Println("serialize register data wrap object err:", err)
		return
	}

	// 将序列化后的信息发布到异步存储管道
	_, err = variable.RedisPool.Publish("AsyncPersistence", string(ret))
	if err != nil {
		log.Error.Println("redis pool publish to async persistence err:", err)
		return
	}
}

// 回执消息
func (rr *RegisterRouter) PostHandle(r api.IRequest) {

	// 包装 ack
	ackPack := packet.NewServerAckPack(constants.SignUpAckOpt, true, []byte("register an account success, please logic again"))
	ret, err := json.Marshal(ackPack)
	if err != nil {
		log.Info.Println("serialize register ack pack object err:", err)
		return
	}

	// 反向客户端发送 ack 数据
	err = r.GetConnection().SendMsg(constants.TCPNetwork, constants.AckOption, 101, ret)
	if err != nil {
		log.Error.Println("register send ack message to client err:", err)
	}
}
