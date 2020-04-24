package routers

import (
	"chalurania/api"
	"chalurania/comet/caller/consumers"
	"chalurania/comet/constants"
	"chalurania/comet/packet"
	"chalurania/comet/router"
	"chalurania/comet/variable"
	"chalurania/service/dao"
	"chalurania/service/log"
	"chalurania/service/model"
	"chalurania/service/scrypt"
	"encoding/json"
	"strconv"
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

		// 登录成功，校验用户的频道是否存在，存在则 publish 消息告诉另一个连接下线， 然后当前连接再 subscribe 订阅用户频道
		// 获取 redis 连接
		redisConn := variable.RedisPool.Pool.Get()
		defer func() {
			if err := redisConn.Close(); err != nil {
				log.Error.Println("redis conn close err:", err)
				return
			}
		}()

		// 用户频道名称
		chanName := "UserChannel:" + strconv.FormatInt(user.UserId,10)
		// 读取 redis 判断用户是否登录
		// 用户频道名定义：key - "user:用户id"，field - channel， value - "UserChannel：用户id"
		res, err := redisConn.Do("hget", "user:" + strconv.FormatInt(user.UserId, 10), "channel")
		if err != nil {
			log.Error.Println("redis hash set user info err:", err)
		}

		if res != nil {
			// 用户已登录
			// publish 消息告诉另一个连接下线
			log.Info.Println("already logged in")
			_, err = variable.RedisPool.Publish(chanName, "account has been logged in on other devices")
			if err != nil {
				log.Error.Println("redis pool publish to async persistence err:", err)
			}
		}

		// 存储用户信息和连接信息
		lr.success = true
		// 处理用户信息(缓存登录状态，订阅频道)
		consumers.UserConsume(user, r)

		// 将频道存储在 redis hash 中
		// 用户频道名定义：key - "user:用户id"，field - channel， value - "UserChannel：用户id"
		_, err = redisConn.Do("hset", "user:" + strconv.FormatInt(user.UserId, 10), "channel", chanName)
		if err != nil {
			log.Error.Println("redis hash set user info err:", err)
		}
	} else {
		// 登录失败
		lr.success = false
	}
}

// 回执消息
func (lr *LoginRouter) PostHandle(r api.IRequest) {
	// 反向客户端发送 ack 数据
	var loginMsg []byte
	if lr.success {
		loginMsg = []byte("login success, now you can get your message")
	} else {
		loginMsg = []byte("login failed, please check your username or password")
	}

	// 包装 ack
	ackPack := packet.NewAckPack(constants.LoginAckOpt, lr.success, loginMsg)
	ret, err := json.Marshal(ackPack)
	if err != nil {
		log.Info.Println("serialize login ack pack object err:", err)
		return
	}

	// 发送回执
	err = r.GetConnection().SendMsg(1, constants.AckOption, 101, ret)
	if err != nil {
		log.Error.Println("login send ack message to client err:", err)
	}
}