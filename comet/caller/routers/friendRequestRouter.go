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
)

// 添加好友请求
type FriendRequestRouter struct {
	router.Router
	success bool
}

func (fr *FriendRequestRouter) Handle(r api.IRequest) {
	// 当前用户信息
	user, err := r.GetConnection().GetProperty("user")
	if err != nil {
		log.Error.Println("conn get user property err:", err)
		fr.success = false
		return
	}
	u := user.(*model.User)

	var stp packet.ServerTransPack
	err = json.Unmarshal(r.GetData(), &stp)
	if err != nil {
		log.Error.Printf("unmarshal server trans pack err: %v\n", err)
		fr.success = false
		return
	}

	// 好友信息
	var friend model.User
	err = json.Unmarshal(stp.Data, &friend)
	if err != nil {
		log.Error.Printf("unmarshal friend err: %v\n", err)
		fr.success = false
		return
	}

	switch stp.Opt {
	case constants.AddUserByNameCommand:
		// 通过用户名添加好友
		// 校验好友用户信息
		userDAO := dao.NewUserDAO(variable.GoDB)
		exist, f := userDAO.QueryUserByName(friend)
		if exist {
			// 添加好友请求信息到持久化 channel
			friendReq := model.FriendRequest{UserId: u.UserId, FriendId: f.UserId}
			// 序列化好友请求对象
			frJson, err := json.Marshal(friendReq)
			if err != nil {
				log.Info.Println("serialize friend request object err:", err)
			}

			// 将好友请求信息包装并序列化
			dw := packet.NewDataPersistWrap(constants.FriendRequestPersistOpt, frJson)
			ret, err := json.Marshal(dw)
			if err != nil {
				log.Info.Println("serialize friend request data wrap object err:", err)
				return
			}

			// 将序列化后的信息发布到异步存储管道
			_, err = variable.RedisPool.Publish("AsyncPersistence", string(ret))
			if err != nil {
				log.Error.Println("redis pool publish to async persistence err:", err)
				return
			}
			fr.success = true
		} else {
			fr.success = false
		}
	case constants.AddUserByIdCommand:
		// 通过用户 id 添加好友
		// 校验好友用户信息
		userDAO := dao.NewUserDAO(variable.GoDB)
		exist, f := userDAO.QueryUserById(friend)
		if exist {
			// 添加好友请求信息到持久化 channel
			friendReq := model.FriendRequest{UserId: u.UserId, FriendId: f.UserId}
			// 序列化好友请求对象
			frJson, err := json.Marshal(friendReq)
			if err != nil {
				log.Info.Println("serialize friend request object err:", err)
			}

			// 将好友请求信息包装并序列化
			dw := packet.NewDataPersistWrap(constants.FriendRequestPersistOpt, frJson)
			ret, err := json.Marshal(dw)
			if err != nil {
				log.Info.Println("serialize friend request data wrap object err:", err)
				return
			}

			// 将序列化后的信息发布到异步存储管道
			_, err = variable.RedisPool.Publish("AsyncPersistence", string(ret))
			if err != nil {
				log.Error.Println("redis pool publish to async persistence err:", err)
				return
			}
			fr.success = true
		} else {
			fr.success = false
		}
	}
}

// 回执消息
func (fr *FriendRequestRouter) PostHandle(r api.IRequest) {
	// 反向客户端发送 ack 数据
	var friendRequestMsg []byte
	if fr.success {
		friendRequestMsg = []byte("send request succeeded, wait for confirmation")
	} else {
		friendRequestMsg = []byte("oops, send friend request failed, check the input information")
	}

	// 包装 ack
	ackPack := packet.NewServerAckPack(constants.FriendRequestAckOpt, fr.success, friendRequestMsg)
	ret, err := json.Marshal(ackPack)
	if err != nil {
		log.Info.Println("serialize friend request ack pack object err:", err)
		return
	}

	// 发送回执
	err = r.GetConnection().SendMsg(constants.TCPNetwork, constants.AckOption, 101, ret)
	if err != nil {
		log.Error.Println("friend request send ack message to client err:", err)
	}
}