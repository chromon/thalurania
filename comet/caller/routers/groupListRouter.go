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

// 群组列表
type GroupListRouter struct {
	router.Router
	sign bool
	msg []byte
}

func (gr *GroupListRouter) Handle(r api.IRequest) {
	// 当前用户信息
	user, err := r.GetConnection().GetProperty("user")
	if err != nil {
		log.Error.Println("conn get user property err:", err)
		return
	}
	u := user.(*model.User)

	groupUserDAO := dao.NewGroupUserDAO(variable.GoDB)

	// 查询是否加入过群组
	count := groupUserDAO.QueryGroupCountByUser(*u)
	if count < 1 {
		gr.msg = []byte("never joined any group")
		return
	}

	// 查询群组
	rows, err := groupUserDAO.QueryGroupByUser(*u)
	defer func() {
		if err = rows.Close(); err != nil {
			panic(err)
		}
	}()

	groupDAO := dao.NewGroupDAO(variable.GoDB)
	var info string
	var bt bytes.Buffer

	// 遍历群组
	for rows.Next() {
		var gu model.GroupUser
		err = rows.Scan(&gu.Id, &gu.GroupId, &gu.UserId, &gu.Label, &gu.Extra, &gu.CreateTime, &gu.UpdateTime)
		if err != nil {
			log.Error.Println("scan group user id err:", err)
			return
		}

		// 查询对应群组信息
		group := model.Group{GroupId: gu.GroupId}
		_, g := groupDAO.QueryGroupByGroupId(group)

		info = g.Name + " (" + strconv.FormatInt(g.GroupId, 10) + ")"
		bt.WriteString(info)
		bt.WriteString(",")
	}

	gr.msg = []byte(bt.String())
	gr.sign = true
}

// 回执消息
func (gr *GroupListRouter) PostHandle(r api.IRequest) {
	// 包装 ack
	ackPack := packet.NewServerAckPack(constants.GroupListAckOpt, gr.sign, gr.msg)
	ret, err := json.Marshal(ackPack)
	if err != nil {
		log.Info.Println("serialize group list ack pack object err:", err)
		return
	}

	// 发送回执
	err = r.GetConnection().SendMsg(constants.TCPNetwork, constants.AckOption, 101, ret)
	if err != nil {
		log.Error.Println("group list send ack message to client err:", err)
	}
}