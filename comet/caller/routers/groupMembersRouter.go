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

// 接受群组邀请
type GroupMembersRouter struct {
	router.Router
	msg []byte
	sign bool
}

func (gr *GroupMembersRouter) Handle(r api.IRequest) {
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
		gr.msg = []byte("group id not found, check it again")
		return
	}

	// 判断当前用户是否是群成员
	groupUserDAO := dao.NewGroupUserDAO(variable.GoDB)
	guExist, _ := groupUserDAO.QueryGroupUserById(*u, *g)
	if !guExist {
		// 当前用户不是群成员
		gr.msg = []byte("not a member of this group")
		return
	}

	// 查询群组成员数量
	count := groupUserDAO.QueryGroupUserCount(*g)
	if count < 1 {
		// 不存在群成员
		gr.msg = []byte("error, no group members found")
		return
	}

	// 查询群成员列表
	rows, err := groupUserDAO.QueryGroupUsers(*g)
	defer func() {
		if err = rows.Close(); err != nil {
			panic(err)
		}
	}()

	userDAO := dao.NewUserDAO(variable.GoDB)
	var info string
	var bt bytes.Buffer

	bt.WriteString("group (" + strconv.FormatInt(g.GroupId, 10) + ") members:")
	bt.WriteString(",")

	// 遍历群组成员信息
	for rows.Next() {
		var groupUser model.GroupUser
		err = rows.Scan(&groupUser.Id, &groupUser.GroupId, &groupUser.UserId, &groupUser.Label, &groupUser.Extra, &groupUser.CreateTime, &groupUser.UpdateTime)
		if err != nil {
			log.Error.Println("scan group user err:", err)
			return
		}

		var gu model.User
		gu.UserId = groupUser.UserId
		_, userInfo := userDAO.QueryUserById(gu)

		info = userInfo.Username + " (" + strconv.FormatInt(userInfo.UserId, 10) + ")"
		bt.WriteString(info)
		bt.WriteString(",")
	}

	gr.msg = []byte(bt.String())
	gr.sign = true
}

// 回执消息
func (gr *GroupMembersRouter) PostHandle(r api.IRequest) {
	// 包装 ack
	ackPack := packet.NewServerAckPack(constants.GroupMembersAckOpt, gr.sign, gr.msg)
	ret, err := json.Marshal(ackPack)
	if err != nil {
		log.Info.Println("serialize group members ack pack object err:", err)
		return
	}

	// 发送回执
	err = r.GetConnection().SendMsg(constants.TCPNetwork, constants.AckOption, 101, ret)
	if err != nil {
		log.Error.Println("group members send ack message to client err:", err)
	}
}