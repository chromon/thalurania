package routers

import (
	"chalurania/api"
	"chalurania/comet/constants"
	"chalurania/comet/packet"
	"chalurania/comet/router"
	"chalurania/comet/variable"
	"chalurania/service/log"
	"chalurania/service/model"
	"encoding/json"
	"strconv"
)

// 登出路由
type LogoutRouter struct {
	router.Router
	success bool
}

func (lr *LogoutRouter) Handle(r api.IRequest) {
	user, err := r.GetConnection().GetProperty("user")
	if err != nil {
		log.Error.Println("conn get user property err:", err)
		lr.success = false
		return
	}

	// 获取 redis 连接
	redisConn := variable.RedisPool.Pool.Get()
	defer func() {
		if err := redisConn.Close(); err != nil {
			log.Error.Println("redis conn close err:", err)
			return
		}
	}()

	u := user.(*model.User)

	// 删除掉 redis 中的用户登录信息
	_, err = redisConn.Do("hdel", "user:" + strconv.FormatInt(u.UserId, 10), "channel")
	if err != nil {
		log.Error.Println("redis hash del user info err:", err)
		lr.success = false
		return
	}
	lr.success = true
}

// 回执消息
func (lr *LogoutRouter) PostHandle(r api.IRequest) {
	// 反向客户端发送 ack 数据
	var logoutMsg []byte
	if lr.success {
		logoutMsg = []byte("logout success, bye")
	} else {
		logoutMsg = []byte("logout failed, system error")
	}

	// 包装 ack
	ackPack := packet.NewServerAckPack(constants.LogoutAckOpt, lr.success, logoutMsg)
	ret, err := json.Marshal(ackPack)
	if err != nil {
		log.Info.Println("serialize logout ack pack object err:", err)
		return
	}

	// 发送回执
	err = r.GetConnection().SendMsg(constants.TCPNetwork, constants.AckOption, 101, ret)
	if err != nil {
		log.Error.Println("logout send ack message to client err:", err)
	}
}