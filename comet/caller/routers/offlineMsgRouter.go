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

// 搜索路由
type OfflineMsgRouter struct {
	router.Router
	msg []byte
	flag bool
}

func (or *OfflineMsgRouter) Handle(r api.IRequest) {
	// 当前用户信息
	user, err := r.GetConnection().GetProperty("user")
	if err != nil {
		log.Error.Println("conn get user property err:", err)
		or.msg = []byte("server error, could not get user info")
		return
	}
	u := user.(*model.User)

	var stp packet.ServerTransPack
	err = json.Unmarshal(r.GetData(), &stp)
	if err != nil {
		log.Error.Printf("unmarshal server trans pack err: %v\n", err)
		return
	}

	// 解析好友用户信息
	var friend model.User
	err = json.Unmarshal(stp.Data, &friend)
	if err != nil {
		log.Error.Printf("unmarshal friend err: %v\n", err)
		return
	}

	userDAO := dao.NewUserDAO(variable.GoDB)

	// 校验好友用户信息
	var f *model.User
	var exist bool
	switch stp.Opt {
	case constants.OfflineMsgByNameCommand:
		// 通过用户名与好友聊天
		// 查询好友信息（根据用户名）
		exist, f = userDAO.QueryUserByName(friend)
		if !exist {
			or.msg = []byte("wrong information entered, friend name (" + friend.Username + ") not found")
			return
		}
	case constants.OfflineMsgByIdCommand:
		// 通过用户 id 与好友聊天
		exist, f = userDAO.QueryUserById(friend)
		if !exist {
			or.msg = []byte("wrong information entered, friend id (" + strconv.FormatInt(friend.UserId,10) + ") not found")
			return
		}
	}

	// 校验是否是真实好友
	friendDAO := dao.NewFriendDAO(variable.GoDB)
	isFriend := friendDAO.QueryFriendById(*u, *f)
	if !isFriend {
		or.msg = []byte("no friendship exists")
		return
	}

	messageDAO := dao.NewMessageDAO(variable.GoDB)
	// 查询好友离线消息数量
	count := messageDAO.QueryOfflineMsgCount(*u, *f)
	if count < 1 {
		// 当前好友不存在离线消息
		or.msg = []byte("no offline messages")
		return
	}

	// 查询好友发送的离线信息
	rows, err := messageDAO.QueryOfflineMessage(*u, *f)
	defer func() {
		if err = rows.Close(); err != nil {
			panic(err)
		}
	}()
	if err != nil {
		log.Error.Println("query offline message err:", err)
	}

	// 离线消息 map
	offlineMsgMap := make(map[string][]byte)

	// 遍历离线消息
	for rows.Next() {
		var om model.Message
		err = rows.Scan(&om.Id, &om.Seq, &om.Content, &om.MessageTypeId, &om.SenderType, &om.SenderId, &om.ReceiverType,
			&om.ReceiverId, &om.ToUserIds, &om.SendTime, &om.Status, &om.CreateTime, &om.UpdateTime)
		if err != nil {
			log.Error.Println("rows scan offline message err:", err)
			return
		}

		// value - 消息信息
		omJson, err := json.Marshal(om)
		if err != nil {
			log.Info.Println("serialize offline message object err:", err)
		}

		offlineMsgMap[strconv.FormatInt(om.Seq, 10)] = omJson

		// 更新离线消息状态
		om.Status = 2
		messageDAO.UpdateMessage(om)
	}

	mapJson, err := json.Marshal(offlineMsgMap)
	if err != nil {
		log.Info.Println("serialize offline message map object err:", err)
	}
	or.msg = mapJson

	or.flag = true
}


// 回执消息
func (or *OfflineMsgRouter) PostHandle(r api.IRequest) {
	// 包装 ack
	ackPack := packet.NewServerAckPack(constants.OfflineMsgAckOpt, or.flag, or.msg)
	ret, err := json.Marshal(ackPack)
	if err != nil {
		log.Info.Println("serialize offline message ack pack object err:", err)
		return
	}

	// 发送回执
	err = r.GetConnection().SendMsg(constants.TCPNetwork, constants.AckOption, 101, ret)
	if err != nil {
		log.Error.Println("offline send ack message to client err:", err)
	}
}