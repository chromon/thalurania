package consumers

import (
	"bytes"
	"chalurania/api"
	"chalurania/comet/constants"
	"chalurania/comet/packet"
	"chalurania/service/log"
	"chalurania/service/model"
	"encoding/json"
	"github.com/gomodule/redigo/redis"
	"strconv"
)

// 处理 channel 订阅到的信息
type UserConsume struct {
	User *model.User
	Request api.IRequest
}

func NewUserConsume(u *model.User, r api.IRequest) *UserConsume {
	return &UserConsume{
		User: u,
		Request: r,
	}
}

func (uc *UserConsume) Consume() func(redis.Message) error {
	return func(msg redis.Message) error {
		log.Info.Printf("user consume recv msg: %s", msg.Data)

		// 服务器内部数据传输包，用于区分 channel 消息的类型（踢人，聊天...）
		var stp packet.ServerTransPack
		err := json.Unmarshal(msg.Data, &stp)
		if err != nil {
			log.Error.Printf("unmarshal server trans pack err: %v\n", err)
			return err
		}

		var ackPack *packet.ServerAckPack
		switch stp.Opt {
		case constants.KickOut:
			// 踢人，迫使另一设备下线
			ackPack = packet.NewServerAckPack(constants.DeviceOffline, true, stp.Data)
		case constants.SendFriendRequest:
			// 发送添加好友信息
			var friend model.User
			err = json.Unmarshal(stp.Data, &friend)
			if err != nil {
				log.Error.Printf("unmarshal stp pack err: %v\n", err)
			}

			// 拼接通知信息
			var bt bytes.Buffer
			bt.WriteString("[NEW] friend request from: ")
			bt.WriteString(friend.Username)
			bt.WriteString(" (")
			bt.WriteString(strconv.FormatInt(friend.UserId,10))
			bt.WriteString(")")

			ackPack = packet.NewServerAckPack(constants.FriendRequestAckOpt, true, []byte(bt.String()))
		case constants.SendMessage:
			// 向好友发送消息
			ackPack = packet.NewServerAckPack(constants.SendMessageAckOpt, true, stp.Data)
		case constants.LoginSuccess:
			// 向用户发送登录成功消息，使用 false 表明直接显示信息，不需要解析
			ackPack = packet.NewServerAckPack(constants.LoginAckOpt, false, stp.Data)
		case constants.SendGroupRequest:
			// 向用户发送群组邀请
			ackPack = packet.NewServerAckPack(constants.GroupRequestAckOpt, true, stp.Data)
		case constants.SendGroupMessage:
			// 群组聊天
			ackPack = packet.NewServerAckPack(constants.SendGroupMessageAckOpt, true, stp.Data)
		}

		// 序列化 ack 并向客户端发送
		ret, err := json.Marshal(ackPack)
		if err != nil {
			log.Info.Println("serialize logic ack pack object err:", err)
			return err
		}

		// 向客户端发送信息
		err = uc.Request.GetConnection().SendMsg(constants.TCPNetwork, constants.AckOption, 101, ret)
		if err != nil {
			log.Error.Println("user consumer message to client err:", err)
			return err
		}
		return nil
	}
}
