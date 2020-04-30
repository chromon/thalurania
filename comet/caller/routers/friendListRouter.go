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

// 好友列表
type FriendListRouter struct {
	router.Router
	success bool
	friendMap map[string][]byte
}

func (fr *FriendListRouter) Handle(r api.IRequest) {
	user, err := r.GetConnection().GetProperty("user")
	if err != nil {
		log.Error.Println("conn get user property err:", err)
		fr.success = false
		return
	}
	u := user.(*model.User)

	// 好友列表 map
	fr.friendMap = make(map[string][]byte)
	friendDAO := dao.NewFriendDAO(variable.GoDB)
	userDAO := dao.NewUserDAO(variable.GoDB)
	// 查询好友
	rows, err := friendDAO.QueryFriend(*u)
	defer func() {
		if err = rows.Close(); err != nil {
			panic(err)
		}
	}()
	if err != nil {
		fr.success = false
		log.Error.Println("query friend err:", err)
		return
	}

	// 遍历好友 id
	for rows.Next() {
		var friend model.User
		err = rows.Scan(&friend.UserId)
		if err != nil {
			fr.success = false
			log.Error.Println("scan friend id err:", err)
			return
		}

		// 由好友 id 查询好友信息并序列化
		_, f := userDAO.QueryUserById(friend)
		fJson, err := json.Marshal(f)
		if err != nil {
			log.Info.Println("serialize friend object err:", err)
		}
		// 将好友信息添加到好友 map 中
		fr.friendMap[strconv.FormatInt(f.UserId,10)] = fJson
	}

	fr.success = true
}

// 回执消息
func (fr *FriendListRouter) PostHandle(r api.IRequest) {
	// 反向客户端发送 ack 数据
	var friendMsg []byte
	if fr.success {
		mapJson, err := json.Marshal(fr.friendMap)
		if err != nil {
			log.Info.Println("serialize friend map object err:", err)
		}
		friendMsg = mapJson
	} else {
		friendMsg = []byte("query friend list failed, try again later")
	}

	// 包装 ack
	ackPack := packet.NewServerAckPack(constants.FriendListAckOpt, fr.success, friendMsg)
	ret, err := json.Marshal(ackPack)
	if err != nil {
		log.Info.Println("serialize friend list ack pack object err:", err)
		return
	}

	// 发送回执
	err = r.GetConnection().SendMsg(constants.TCPNetwork, constants.AckOption, 101, ret)
	if err != nil {
		log.Error.Println("friend list send ack message to client err:", err)
	}
}