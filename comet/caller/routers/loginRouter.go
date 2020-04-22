package routers

import (
	"chalurania/api"
	"chalurania/comet/router"
	"chalurania/comet/variable"
	"chalurania/service/dao"
	"chalurania/service/log"
	"chalurania/service/model"
	"chalurania/service/scrypt"
	"encoding/json"
)

// 登录路由
type LoginRouter struct {
	router.Router
	success bool
}

// 登录处理
func (lr *LoginRouter) Handle(r api.IRequest) {
	log.Info.Println("received from client message id:", r.GetMsgID(), " data:", string(r.GetData()))

	var u model.User
	err := json.Unmarshal(r.GetData(), &u)
	if err != nil {
		log.Error.Printf("unmarshal user err=%v\n", err)
	}
	// 加密密码
	u.Password = scrypt.Crypto(u.Password)

	// 校验用户信息
	userDAO := dao.NewUserDAO(variable.GoDB)
	exist, user := userDAO.QueryUserByNamePass(u)
	if exist {
		// 存储用户信息和连接信息
		lr.success = true
		log.Info.Println(user.Nickname)
		log.Info.Println(user.Password)
	} else {
		// 登录失败
		lr.success = false
		return
	}




	// 订阅用户频道


	//// 将注册信息包装并序列化
	//dw := packet.NewDataWrap(1, r.GetData())
	//ret, err := json.Marshal(dw)
	//if err != nil {
	//	log.Info.Println("serialize register data wrap object err:", err)
	//	return
	//}
	//
	//// 将序列化后的信息发布到异步存储管道
	//_, err = variable.RedisPool.Publish("AsyncPersistence", string(ret))
	//if err != nil {
	//	log.Error.Println("redis pool publish to async persistence err:", err)
	//}
}

// 回执消息
func (lr *LoginRouter) PostHandle(r api.IRequest) {
	// 反向客户端发送 ack 数据
	var loginMsg []byte
	if lr.success {
		loginMsg = []byte("login success")
	} else {
		loginMsg = []byte("login failed, please check your username or password")
	}
	err := r.GetConnection().SendMsg(1, 0, 101, loginMsg)
	if err != nil {
		log.Error.Println("login send ack message to client err:", err)
	}
}