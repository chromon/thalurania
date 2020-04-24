package consumers

import (
	"chalurania/api"
	"chalurania/comet/constants"
	"chalurania/comet/variable"
	"chalurania/service/log"
	"chalurania/service/model"
	"context"
	"github.com/gomodule/redigo/redis"
	"strconv"
)

// 处理用户信息(缓存登录状态，订阅频道)
func UserConsume(user *model.User, r api.IRequest) {
	// 读取订阅频道发送的数据，并使用 tcp 连接发送到客户端

	// 获取 redis 连接
	redisConn := variable.RedisPool.Pool.Get()
	defer func() {
		if err := redisConn.Close(); err != nil {
			log.Error.Println("redis conn close err:", err)
			return
		}
	}()

	// 订阅自己的频道
	ctx, _ := context.WithCancel(context.Background())
	// 用户频道名称
	chanName := "UserChannel:" + strconv.FormatInt(user.UserId,10)
	go func() {
		//log.Info.Println(chanName)

		// 从订阅频道接收消息处理回调函数
		consume := func(msg redis.Message) error {
			log.Info.Printf("recv msg: %s", msg.Data)
			err := r.GetConnection().SendMsg(1, constants.AckOption, 101, msg.Data)
			if err != nil {
				log.Error.Println("user consumer message to client err:", err)
			}
			return nil
		}

		// 订阅频道
		err := variable.RedisPool.Subscribe(ctx, consume, chanName)
		if err != nil {
			log.Error.Println("subscribe UserConsume channel err:", err)
		}
	}()


}