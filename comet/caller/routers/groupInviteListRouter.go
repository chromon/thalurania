package routers

import (
	"bytes"
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

// 群组邀请列表路由
type GroupInviteListRouter struct {
	router.Router
	sign bool
	msg string
}

func (gr *GroupInviteListRouter) Handle(r api.IRequest) {
	// 当前用户信息
	user, err := r.GetConnection().GetProperty("user")
	if err != nil {
		log.Error.Println("conn get user property err:", err)
		return
	}
	u := user.(*model.User)

	// 群组邀请
	groupInviteDAO := dao.NewGroupInviteDAO(variable.GoDB)
	// 查询群组邀请数量
	count := groupInviteDAO.QueryGroupInviteCount(*u)
	if count < 1 {
		gr.msg = "no group invitations"
		return
	}

	// 查询群组邀请
	rows, err := groupInviteDAO.QueryGroupInvite(*u)
	defer func() {
		if err = rows.Close(); err != nil {
			panic(err)
		}
	}()

	userDAO := dao.NewUserDAO(variable.GoDB)
	var info string
	var bt bytes.Buffer

	// 遍历群组邀请
	for rows.Next() {
		var gi model.GroupInvite
		err = rows.Scan(&gi.Id, &gi.UserId, &gi.FriendId, &gi.GroupId, &gi.Del)
		if err != nil {
			log.Error.Println("scan group invite id err:", err)
			return
		}

		// 查询邀请人信息
		friend := model.User{UserId: gi.UserId}
		exist, f := userDAO.QueryUserById(friend)
		if !exist {
			log.Error.Println("query user failed")
		}

		info = "friend " + f.Username + " (" + strconv.FormatInt(f.UserId, 10) + ") invite you to join group (" + strconv.FormatInt(gi.GroupId, 10) + ")"

		bt.WriteString(info)
		bt.WriteString(",")
	}

	gr.msg = bt.String()
	gr.sign = true
}

// 回执消息
func (gr *GroupInviteListRouter) PostHandle(r api.IRequest) {
	// 包装 ack
	ackPack := packet.NewServerAckPack(constants.GroupInviteAckOpt, gr.sign, []byte(gr.msg))
	ret, err := json.Marshal(ackPack)
	if err != nil {
		log.Info.Println("serialize query group invite ack pack object err:", err)
		return
	}

	// 发送回执
	err = r.GetConnection().SendMsg(constants.TCPNetwork, constants.AckOption, 101, ret)
	if err != nil {
		log.Error.Println("query group invite send ack message to client err:", err)
	}
}