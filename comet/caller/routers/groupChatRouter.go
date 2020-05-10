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
	"time"
)

type GroupChatRouter struct {
	router.Router
	sign bool
	msg []byte
}

func (gc *GroupChatRouter) Handle(r api.IRequest) {
	// 当前用户信息
	user, err := r.GetConnection().GetProperty("user")
	if err != nil {
		log.Error.Println("conn get user property err:", err)
		return
	}
	u := user.(*model.User)

	var message model.Message
	err = json.Unmarshal(r.GetData(), &message)
	if err != nil {
		log.Error.Printf("unmarshal message err: %v\n", err)
	}

	// 消息内容
	content := message.Content

	// 由群组 id 查询群组
	groupDAO := dao.NewGroupDAO(variable.GoDB)
	group := model.Group{GroupId: message.ReceiverId}
	gExist, g := groupDAO.QueryGroupByGroupId(group)
	if !gExist {
		gc.msg = []byte("group id not exist, check it again")
		return
	}

	// 查询当前用户是否在群组中
	groupUserDAO := dao.NewGroupUserDAO(variable.GoDB)
	guExist, _ := groupUserDAO.QueryGroupUserById(*u, *g)
	if !guExist {
		// 当前用户不是群成员
		gc.msg = []byte("not a member of this group")
		return
	}

	// 保存群组消息
	messageDAO := dao.NewMessageDAO(variable.GoDB)
	message = model.Message{Seq: r.GetMsgID(), Content: content, MessageTypeId: 2, SenderType: 1, SenderId: u.UserId,
		ReceiverType: 2, ReceiverId: g.GroupId, ToUserIds: "", SendTime: time.Now(), Status: 2, CreateTime: time.Now(), UpdateTime: time.Now()}
	insertId, err := messageDAO.AddMessage(message)
	if err != nil {
		log.Error.Println("save message err:", err)
		return
	}
	log.Info.Println("insertId:", insertId)

	// 保存当前用户在群组中接收的最后一条消息的位置
	groupOfflineMsgDAO := dao.NewGroupOfflineMessageDAO(variable.GoDB)
	groupOfflineMsg := model.GroupOfflineMessage{UserId: u.UserId, GroupId: g.GroupId, MessageId: insertId}
	_, err = groupOfflineMsgDAO.AddGroupOfflineMessage(groupOfflineMsg)

	// 向群组频道发送消息
	info := []byte("[new] message from " + u.Username +" (" + strconv.FormatInt(u.UserId,10) +") in group (" + strconv.FormatInt(g.GroupId, 10) + "): \n" + message.Content)

	// 打包 ack 消息
	serverTransPack := packet.NewServerTransPack(constants.SendGroupMessage, info)
	ret, err := json.Marshal(serverTransPack)
	if err != nil {
		log.Info.Println("serialize server trans pack (send group message) object err:", err)
		return
	}

	// publish 消息(pack), consumer 将消息发送给好友
	groupChanName := "GroupChannel:" + strconv.FormatInt(g.GroupId,10)
	_, err = variable.RedisPool.Publish(groupChanName, string(ret))
	if err != nil {
		log.Error.Println("redis pool publish to group channel err:", err)
		return
	}

	gc.msg = []byte("send group message success")
}


// 回执消息
func (gc *GroupChatRouter) PostHandle(r api.IRequest) {
	// 包装 ack
	ackPack := packet.NewServerAckPack(constants.SendGroupMessageAckOpt, true, gc.msg)
	ret, err := json.Marshal(ackPack)
	if err != nil {
		log.Info.Println("serialize group chat ack pack object err:", err)
		return
	}

	// 发送回执
	err = r.GetConnection().SendMsg(constants.TCPNetwork, constants.AckOption, 101, ret)
	if err != nil {
		log.Error.Println("group chat send ack message to client err:", err)
	}
}