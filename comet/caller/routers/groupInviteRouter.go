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

// 群组邀请
type GroupInviteRouter struct {
	router.Router
	success bool
	msg []byte
}

func (gr *GroupInviteRouter) Handle(r api.IRequest) {
	// 当前用户信息
	user, err := r.GetConnection().GetProperty("user")
	if err != nil {
		log.Error.Println("conn get user property err:", err)
		return
	}
	u := user.(*model.User)

	var stp packet.ServerTransPack
	err = json.Unmarshal(r.GetData(), &stp)
	if err != nil {
		log.Error.Printf("unmarshal server trans pack err: %v\n", err)
		return
	}

	// 好友信息
	var friend model.User
	err = json.Unmarshal(stp.Data, &friend)
	if err != nil {
		log.Error.Printf("unmarshal friend err: %v\n", err)
		return
	}

	userDAO := dao.NewUserDAO(variable.GoDB)
	var f *model.User
	var fExist bool

	// 校验好友信息
	switch stp.Opt {
	case constants.GroupInviteByNameCommand:
		fExist, f = userDAO.QueryUserByName(friend)
		if !fExist {
			gr.msg = []byte("friend username not exist, check it again")
			return
		}
	case constants.GroupInviteByIdCommand:
		fExist, f = userDAO.QueryUserById(friend)
		if !fExist {
			gr.msg = []byte("friend user id not exist, check it again")
			return
		}
	}

	// 群组 id
	groupId, err := strconv.ParseInt(friend.Extra, 10, 64)
	group := model.Group{GroupId: groupId}
	// 查询群组是否存在
	groupDAO := dao.NewGroupDAO(variable.GoDB)
	gExist, g := groupDAO.QueryGroupByGroupId(group)
	if !gExist {
		gr.msg = []byte("group id not exist, check it again")
		return
	}

	// 查询当前用户是否有邀请权限
	groupUserDAO := dao.NewGroupUserDAO(variable.GoDB)
	exist, _ := groupUserDAO.QueryGroupUserById(*u, *g)
	if !exist {
		// 当前用户不是群组成员
		gr.msg = []byte("no permission to invite new users")
		return
	}

	// 查询好友是否已经在群组中
	exist, _ = groupUserDAO.QueryGroupUserById(*f, *g)
	if exist {
		// 当前好友已是群组成员
		gr.msg = []byte("friend are already members of this group")
		return
	}

	// 向好友 channel 发送邀请
	chanName := "UserChannel:" + strconv.FormatInt(f.UserId,10)
	inviteMsg := "user <" + u.Username + "> (" + strconv.FormatInt(u.UserId, 10) + ") invite you join group (" + strconv.FormatInt(g.GroupId, 10) + ")"

	serverTransPack := packet.NewServerTransPack(constants.SendGroupRequest, []byte(inviteMsg))
	ret, err := json.Marshal(serverTransPack)
	if err != nil {
		log.Info.Println("serialize server trans pack (send friend request) object err:", err)
		return
	}
	// 发送
	_, err = variable.RedisPool.Publish(chanName, string(ret))
	if err != nil {
		log.Error.Println("redis pool publish to user channel err:", err)
		return
	}

	// 存储群组请求
	groupInvite := model.GroupInvite{UserId: u.UserId, FriendId: f.UserId, GroupId: groupId}
	// 序列化群组请求对象
	giJson, err := json.Marshal(groupInvite)
	if err != nil {
		log.Info.Println("serialize group request object err:", err)
		return
	}
	// 包装
	dw := packet.NewDataPersistWrap(constants.GroupRequestPersistOpt, giJson)
	re, err := json.Marshal(dw)
	if err != nil {
		log.Info.Println("serialize group request data wrap object err:", err)
		return
	}
	// 将序列化后的信息发布到异步存储管道
	go func() {
		_, err = variable.RedisPool.Publish("AsyncPersistence", string(re))
		if err != nil {
			log.Error.Println("redis pool publish to async persistence err:", err)
			return
		}
	}()

	gr.msg = []byte("send group invite success")
}

// 回执消息
func (gr *GroupInviteRouter) PostHandle(r api.IRequest) {
	// 包装 ack
	ackPack := packet.NewServerAckPack(constants.GroupRequestAckOpt, true, gr.msg)
	ret, err := json.Marshal(ackPack)
	if err != nil {
		log.Info.Println("serialize create group ack pack object err:", err)
		return
	}

	// 发送回执
	err = r.GetConnection().SendMsg(constants.TCPNetwork, constants.AckOption, 101, ret)
	if err != nil {
		log.Error.Println("create group send ack message to client err:", err)
	}
}