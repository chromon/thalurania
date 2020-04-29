package main

import (
	"chalurania/api"
	"chalurania/comet"
	"chalurania/comet/caller/routers"
	"chalurania/comet/constants"
	"chalurania/comet/variable"
	"chalurania/service/log"
	"chalurania/service/model"
	"strconv"
)

// 创建连接时执行
func OnConnectionStart(conn api.IConnection) {
	log.Info.Println("on connection start called...")
}

// 断开连接时执行
func OnConnectionLost(conn api.IConnection) {
	//log.Info.Println("on connection lost called...")

	// 获取属性
	user, err := conn.GetProperty("user")
	if err != nil {
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
		return
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
	s.AddRouter(constants.SearchOption, &routers.SearchRouter{})
	s.AddRouter(constants.FriendRequestOption, &routers.FriendRequestRouter{})
	s.AddRouter(constants.FriendReqListOption, &routers.FriendReqListRouter{})

	// 开启服务
	s.Serve()
}