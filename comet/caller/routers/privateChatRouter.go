package routers

import (
	"chalurania/api"
	"chalurania/comet/constants"
	"chalurania/comet/packet"
	"chalurania/comet/router"
	"chalurania/comet/variable"
	"chalurania/service/dao"
	"chalurania/service/log"
	"chalurania/service/model"
	"encoding/json"
	"strconv"
	"time"
)

// 接受好友请求
type PrivateChatRouter struct {
	router.Router
	msg []byte
}

func (pc *PrivateChatRouter) Handle(r api.IRequest) {
	// 当前用户信息
	user, err := r.GetConnection().GetProperty("user")
	if err != nil {
		log.Error.Println("conn get user property err:", err)
		pc.msg = []byte("server error, could not get user info")
		return
	}
	u := user.(*model.User)

	var stp packet.ServerTransPack
	err = json.Unmarshal(r.GetData(), &stp)
	if err != nil {
		log.Error.Printf("unmarshal server trans pack err: %v\n", err)
		return
	}

	// 好友用户信息
	var friend model.User
	err = json.Unmarshal(stp.Data, &friend)
	if err != nil {
		log.Error.Printf("unmarshal friend err: %v\n", err)
		return
	}

	userDAO := dao.NewUserDAO(variable.GoDB)

	// 校验好友用户信息
	var f *model.User
	var exist bool
	switch stp.Opt {
	case constants.ChatWithFriendByNameCommand:
		// 通过用户名与好友聊天
		// 查询好友信息（根据用户名）
		exist, f = userDAO.QueryUserByName(friend)
		if !exist {
			pc.msg = []byte("wrong information entered, friend name (" + friend.Username + ") not found")
			return
		}
	case constants.ChatWithFriendByIdCommand:
		// 通过用户 id 与好友聊天
		exist, f = userDAO.QueryUserById(friend)
		if !exist {
			pc.msg = []byte("wrong information entered, friend id (" + strconv.FormatInt(friend.UserId,10) + ") not found")
			return
		}
	}

	// 校验是否是真实好友
	friendDAO := dao.NewFriendDAO(variable.GoDB)
	isFriend := friendDAO.QueryFriendById(*u, *f)
	if !isFriend {
		pc.msg = []byte("no friendship exists")
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
	// 读取 redis 判断好友是否在线
	// 用户频道名定义：key - "user:用户id"，field - channel， value - "UserChannel：用户id"
	res, err := redisConn.Do("hget", "user:" + strconv.FormatInt(f.UserId, 10), "channel")
	if err != nil {
		log.Error.Println("redis hash get user info err:", err)
		return
	}

	// 消息内容，暂时由 extra 字段保存消息，便于传输
	f.Extra = friend.Extra

	if res != nil {
		// 好友在线，发送消息
		// 存储消息
		message := model.Message{Seq: r.GetMsgID(), Content: f.Extra, MessageTypeId: 1, SenderType: 1, SenderId: u.UserId,
			ReceiverType: 1, ReceiverId: f.UserId, SendTime: time.Now(), Status: 2, CreateTime: time.Now(), UpdateTime: time.Now()}
		// 序列化消息
		msgJson, err := json.Marshal(message)
		if err != nil {
			log.Info.Println("serialize message object err:", err)
			return
		}
		// 将信息包装并序列化
		dw := packet.NewDataPersistWrap(constants.MessagePersistenceOpt, msgJson)
		dwJson, err := json.Marshal(dw)
		if err != nil {
			log.Info.Println("serialize message data wrap object err:", err)
			return
		}

		// 将序列化后的信息发布到异步存储管道
		go func(){
			_, err = variable.RedisPool.Publish("AsyncPersistence", string(dwJson))
			if err != nil {
				log.Error.Println("redis pool publish to async persistence err:", err)
				return
			}
		}()

		pc.msg = []byte("[new] message from " + u.Username +" (" + strconv.FormatInt(u.UserId,10) +"): \n" + message.Content)

		// 打包 ack 消息
		serverTransPack := packet.NewServerTransPack(constants.SendMessage, pc.msg)
		ret, err := json.Marshal(serverTransPack)
		if err != nil {
			log.Info.Println("serialize server trans pack (send message) object err:", err)
			return
		}

		// publish 消息(pack), consumer 将消息发送给好友
		chanName := "UserChannel:" + strconv.FormatInt(f.UserId,10)
		_, err = variable.RedisPool.Publish(chanName, string(ret))
		if err != nil {
			log.Error.Println("redis pool publish to user channel err:", err)
			return
		}
	} else {
		// 好友离线，存储离线消息
		// 存储消息（将消息置为未读状态）
		message := model.Message{Seq: r.GetMsgID(), Content: f.Extra, MessageTypeId: 1, SenderType: 1, SenderId: u.UserId,
			ReceiverType: 1, ReceiverId: f.UserId, SendTime: time.Now(), Status: 1, CreateTime: time.Now(), UpdateTime: time.Now()}
		// 序列化消息
		msgJson, err := json.Marshal(message)
		if err != nil {
			log.Info.Println("serialize message object err:", err)
			return
		}
		// 将信息包装并序列化
		dw := packet.NewDataPersistWrap(constants.MessagePersistenceOpt, msgJson)
		dwJson, err := json.Marshal(dw)
		if err != nil {
			log.Info.Println("serialize message data wrap object err:", err)
			return
		}

		// 将序列化后的信息发布到异步存储管道
		go func(){
			_, err = variable.RedisPool.Publish("AsyncPersistence", string(dwJson))
			if err != nil {
				log.Error.Println("redis pool publish to async persistence err:", err)
				return
			}
		}()
	}

	pc.msg = []byte("send message success")
}

// 回执消息
func (pc *PrivateChatRouter) PostHandle(r api.IRequest) {
	// 包装 ack
	ackPack := packet.NewServerAckPack(constants.SendMessageAckOpt, true, pc.msg)
	ret, err := json.Marshal(ackPack)
	if err != nil {
		log.Info.Println("serialize login ack pack object err:", err)
		return
	}

	// 发送回执
	err = r.GetConnection().SendMsg(constants.TCPNetwork, constants.AckOption, 101, ret)
	if err != nil {
		log.Error.Println("login send ack message to client err:", err)
	}
}