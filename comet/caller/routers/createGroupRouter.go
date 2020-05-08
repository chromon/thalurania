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

type CreateGroupRouter struct {
	router.Router
	msg []byte
}

func (cr *CreateGroupRouter) Handle(r api.IRequest) {
	// 当前用户对象
	user, err := r.GetConnection().GetProperty("user")
	if err != nil {
		log.Error.Println("conn get user property err:", err)
		return
	}
	u := user.(*model.User)

	// 新建空白群组
	groupId := variable.IdWorker.GetId()
	group := model.Group{GroupId: groupId, Name: "group_" + strconv.FormatInt(groupId, 10)}
	// 插入群组
	groupDAO := dao.NewGroupDAO(variable.GoDB)
	insertId, err := groupDAO.AddGroup(group)
	if err != nil {
		cr.msg = []byte("create group failed, try again")
		return
	}
	log.Info.Println(insertId)

	// 查询新建组信息
	group.Id = insertId
	log.Info.Println(group.Id)
	exist, g := groupDAO.QueryGroupById(group)
	if !exist {
		log.Error.Println("query new group error")
		cr.msg = []byte("system error, try again later")
		return
	}

	// 将当前用户添加到群组中
	groupUserDAO := dao.NewGroupUserDAO(variable.GoDB)
	_, err = groupUserDAO.AddGroupUser(*g, *u)
	if err != nil {
		log.Error.Println("add user to group err:", err)
		cr.msg = []byte("system error, try again later")
		return
	}

	cr.msg = []byte("create group (" + strconv.FormatInt(g.GroupId, 10) + ") success")

}

// 回执消息
func (cr *CreateGroupRouter) PostHandle(r api.IRequest) {
	// 包装 ack
	ackPack := packet.NewServerAckPack(constants.CreateGroupAckOpt, true, cr.msg)
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