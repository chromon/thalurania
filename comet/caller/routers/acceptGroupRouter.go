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

// 接受群组邀请
type AcceptGroupRouter struct {
	router.Router
	msg []byte
}

func (ar *AcceptGroupRouter) Handle(r api.IRequest) {
	// 当前用户信息
	user, err := r.GetConnection().GetProperty("user")
	if err != nil {
		log.Error.Println("conn get user property err:", err)
		return
	}
	u := user.(*model.User)

	// 群组信息
	var group model.Group
	err = json.Unmarshal(r.GetData(), &group)
	if err != nil {
		log.Error.Printf("unmarshal group err: %v\n", err)
		return
	}

	// 判断群组是否存在
	groupDAO := dao.NewGroupDAO(variable.GoDB)
	gExist, g := groupDAO.QueryGroupByGroupId(group)
	if !gExist {
		ar.msg = []byte("group id not found, check it again")
		return
	}

	// 是否存在邀请
	groupInviteDAO := dao.NewGroupInviteDAO(variable.GoDB)
	count := groupInviteDAO.QueryGroupInviteByGroupId(*u, *g)
	if count < 1 {
		// 邀请不存在
		ar.msg = []byte("invitation not found, check it again")
		return
	}

	// 删除已有邀请
	groupInvite := model.GroupInvite{FriendId: u.UserId, GroupId: g.GroupId}
	groupInviteDAO.Update(groupInvite)

	// 将用户添加到组
	groupUserDAO := dao.NewGroupUserDAO(variable.GoDB)
	_, err = groupUserDAO.AddGroupUser(*g, *u)
	if err != nil {
		log.Error.Println("add user to group err:", err)
		return
	}

	ar.msg = []byte("join group successfully")
}

// 回执消息
func (ar *AcceptGroupRouter) PostHandle(r api.IRequest) {
	// 包装 ack
	ackPack := packet.NewServerAckPack(constants.AcceptGroupInviteAckOpt, true, ar.msg)
	ret, err := json.Marshal(ackPack)
	if err != nil {
		log.Info.Println("serialize accept group invite ack pack object err:", err)
		return
	}

	// 发送回执
	err = r.GetConnection().SendMsg(constants.TCPNetwork, constants.AckOption, 101, ret)
	if err != nil {
		log.Error.Println("accept group invite send ack message to client err:", err)
	}
}