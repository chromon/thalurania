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
)

// 登出路由
type FriendReqListRouter struct {
	router.Router
	success bool
	frArray [2][]byte
}

func (lr *FriendReqListRouter) Handle(r api.IRequest) {
	user, err := r.GetConnection().GetProperty("user")
	if err != nil {
		log.Error.Println("conn get user property err:", err)
		lr.success = false
		return
	}
	u := user.(*model.User)

	// 发送的好友请求列表
	sentFrMap := make(map[string][]byte)
	friendReqDAO := dao.NewFriendRequestDAO(variable.GoDB)
	// 查询已发送的好友请求列表
	rows, err := friendReqDAO.QuerySentFriendReq(*u)
	defer func() {
		if err = rows.Close(); err != nil {
			panic(err)
		}
	}()
	if err != nil {
		log.Error.Println("query sent friend req err:", err)
	}

	userDAO := dao.NewUserDAO(variable.GoDB)

	// 遍历发送好友请求列表
	for rows.Next() {
		var fr model.FriendRequest
		err = rows.Scan(&fr.Id, &fr.UserId, &fr.FriendId, &fr.Del)
		if err != nil {
			lr.success = false
			log.Error.Println("rows scan friend request err:", err)
			return
		}

		// 查询好友信息
		var f = model.User{UserId: fr.FriendId}
		_, friend := userDAO.QueryUserById(f)

		fJson, err := json.Marshal(friend)
		if err != nil {
			log.Info.Println("serialize friend object err:", err)
		}

		sentFrMap[strconv.FormatInt(friend.UserId,10)] = fJson
	}

	// 接收的好友请求列表
	receivedFrMap := make(map[string][]byte)
	rows2, err := friendReqDAO.QueryReceiveFriendReq(*u)
	defer func() {
		if err = rows2.Close(); err != nil {
			panic(err)
		}
	}()
	if err != nil {
		log.Error.Println(err)
	}
	for rows2.Next() {
		var fr model.FriendRequest
		err = rows2.Scan(&fr.Id, &fr.UserId, &fr.FriendId, &fr.Del)
		if err != nil {
			lr.success = false
			return
		}

		// 查询好友信息
		var f = model.User{UserId: fr.UserId}
		_, friend := userDAO.QueryUserById(f)

		fJson, err := json.Marshal(friend)
		if err != nil {
			log.Info.Println("serialize friend object err:", err)
		}

		receivedFrMap[strconv.FormatInt(friend.UserId,10)] = fJson
	}

	// 序列化 list
	sentFrMapJson, err := json.Marshal(sentFrMap)
	if err != nil {
		log.Info.Println("serialize sent friend request list err:", err)
		return
	}
	receivedFrMapJson, err := json.Marshal(receivedFrMap)
	if err != nil {
		log.Info.Println("serialize received friend request list err:", err)
		return
	}

	lr.success = true
	lr.frArray = [2][]byte{sentFrMapJson, receivedFrMapJson}
}

// 回执消息
func (lr *FriendReqListRouter) PostHandle(r api.IRequest) {

	// 反向客户端发送 ack 数据
	var frMsg []byte
	if lr.success {
		frArrayJson, err := json.Marshal(lr.frArray)
		if err != nil {
			log.Info.Println("serialize friend request list object err:", err)
		}
		frMsg = frArrayJson
	} else {
		frMsg = []byte("query friend request failed, try again later")
	}

	// 包装 ack
	ackPack := packet.NewServerAckPack(constants.FriendReqListAckOpt, lr.success, frMsg)
	ret, err := json.Marshal(ackPack)
	if err != nil {
		log.Info.Println("serialize friend request list ack pack object err:", err)
		return
	}

	// 发送回执
	err = r.GetConnection().SendMsg(constants.TCPNetwork, constants.AckOption, 101, ret)
	if err != nil {
		log.Error.Println("friend request list send ack message to client err:", err)
	}
}