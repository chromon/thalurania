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

// 接受好友请求
type AcceptFriendRouter struct {
	router.Router
	success bool
}

func (af *AcceptFriendRouter) Handle(r api.IRequest) {
	// 当前用户信息
	user, err := r.GetConnection().GetProperty("user")
	if err != nil {
		log.Error.Println("conn get user property err:", err)
		af.success = false
		return
	}
	u := user.(*model.User)

	var stp packet.ServerTransPack
	err = json.Unmarshal(r.GetData(), &stp)
	if err != nil {
		log.Error.Printf("unmarshal server trans pack err: %v\n", err)
		af.success = false
		return
	}

	// 好友用户信息
	var friend model.User
	err = json.Unmarshal(stp.Data, &friend)
	if err != nil {
		log.Error.Printf("unmarshal friend err: %v\n", err)
		af.success = false
		return
	}

	switch stp.Opt {
	case constants.AcceptFriendByNameCommand:
		// 通过用户名接受好友请求

		// 查询好友信息（根据用户名）
		userDAO := dao.NewUserDAO(variable.GoDB)
		exist, f := userDAO.QueryUserByName(friend)
		if !exist {
			af.success = false
			return
		}

		// 查询当前请求是否存在
		friendRequestDAO := dao.NewFriendRequestDAO(variable.GoDB)
		exist, fr := friendRequestDAO.QueryFriendReq(*f, *u)
		if !exist {
			af.success = false
			return
		}

		// 更新 friendRequest 删除信息
		friendRequestDAO.UpdateFriendReq(*fr)

		// 查询是否有反向好友请求，有则置为删除
		exist2, fr2 := friendRequestDAO.QueryFriendReq(*u, *f)
		if exist2 {
			// 存在反向的好友请求信息，更新 friendRequest 删除信息
			friendRequestDAO.UpdateFriendReq(*fr2)
		}

		// 添加好友（互相）
		friendDAO := dao.NewFriendDAO(variable.GoDB)
		var friendInfo = model.Friend{UserId: u.UserId, FriendId: f.UserId}
		friendDAO.AddFriend(friendInfo)
		friendInfo = model.Friend{UserId: f.UserId, FriendId: u.UserId}
		friendDAO.AddFriend(friendInfo)

		af.success = true

	case constants.AcceptFriendByIdCommand:
		// 通过用户 id 接受好友请求
		// 查询好友信息（根据用户名）
		userDAO := dao.NewUserDAO(variable.GoDB)
		exist, f := userDAO.QueryUserById(friend)
		if !exist {
			af.success = false
			return
		}

		// 查询当前请求是否存在
		friendRequestDAO := dao.NewFriendRequestDAO(variable.GoDB)
		exist, fr := friendRequestDAO.QueryFriendReq(*f, *u)
		if !exist {
			af.success = false
			return
		}

		// 更新 friendRequest 删除信息
		friendRequestDAO.UpdateFriendReq(*fr)

		// 查询是否有反向好友请求，有则置为删除
		exist2, fr2 := friendRequestDAO.QueryFriendReq(*u, *f)
		if exist2 {
			// 存在反向的好友请求信息，更新 friendRequest 删除信息
			friendRequestDAO.UpdateFriendReq(*fr2)
		}

		// 添加好友（互相）
		friendDAO := dao.NewFriendDAO(variable.GoDB)
		var friendInfo = model.Friend{UserId: u.UserId, FriendId: f.UserId}
		friendDAO.AddFriend(friendInfo)
		friendInfo = model.Friend{UserId: f.UserId, FriendId: u.UserId}
		friendDAO.AddFriend(friendInfo)

		af.success = true
	}
}

// 回执消息
func (af *AcceptFriendRouter) PostHandle(r api.IRequest) {
	// 反向客户端发送 ack 数据
	var acceptFriendMsg []byte
	if af.success {
		acceptFriendMsg = []byte("add friend success, now you can chat with friend")
	} else {
		acceptFriendMsg = []byte("oops, add friend fail")
	}

	// 包装 ack
	ackPack := packet.NewServerAckPack(constants.SearchAckOpt, af.success, acceptFriendMsg)
	ret, err := json.Marshal(ackPack)
	if err != nil {
		log.Info.Println("serialize search ack pack object err:", err)
		return
	}

	// 发送回执
	err = r.GetConnection().SendMsg(constants.TCPNetwork, constants.AcceptFriendRepAckOpt, 101, ret)
	if err != nil {
		log.Error.Println("search send ack message to client err:", err)
	}
}