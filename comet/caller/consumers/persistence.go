package consumers

import (
	"chalurania/comet/constants"
	"chalurania/comet/packet"
	"chalurania/comet/variable"
	"chalurania/service/dao"
	"chalurania/service/log"
	"chalurania/service/model"
	"chalurania/service/scrypt"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
)

// 消费者订阅消息处理回调函数
func Consume(msg redis.Message) error {
	// json 解析 data wrap 数据
	var dw packet.DataPersistWrap
	err := json.Unmarshal(msg.Data, &dw)
	if err != nil {
		fmt.Printf("unmarshal data wrap err=%v\n", err)
	}

	switch dw.Opt {
	case constants.SignUpPersistenceOpt:
		// 插入新用户
		err = insertUser(dw.Model)
		if err != nil {
			return err
		}
	case constants.FriendRequestPersistOpt:
		// 插入新好友请求
		err = insertFriendRequest(dw.Model)
		if err != nil {
			return err
		}
	case constants.MessagePersistenceOpt:
		// 插入新消息
		err = insertMessage(dw.Model)
		if err != nil {
			return err
		}
	case constants.GroupRequestPersistOpt:
		// 插入群组请求
		err = insertGroupInvite(dw.Model)
		if err != nil {
			return err
		}
	}

	return nil
}

// 插入新用户
func insertUser(m []byte) error {
	// json 解析 user 数据
	var u model.User
	err := json.Unmarshal(m, &u)
	if err != nil {
		log.Error.Printf("unmarshal user err: %v\n", err)
	}
	// 加密密码
	u.Password = scrypt.Crypto(u.Password)

	// 添加用户
	userDAO := dao.NewUserDAO(variable.GoDB)
	_, err = userDAO.AddUser(u)
	if err != nil {
		return errors.New("insert user error")
	}

	return nil
}

// 插入新好友请求
func insertFriendRequest(m []byte) error {
	// json 解析 friend request 数据
	var fr model.FriendRequest
	err := json.Unmarshal(m, &fr)
	if err != nil {
		log.Error.Printf("unmarshal friend request err: %v\n", err)
	}

	// 添加好友请求
	friendRequestDAO := dao.NewFriendRequestDAO(variable.GoDB)
	_, err = friendRequestDAO.AddFriendRequest(fr)
	if err != nil {
		return errors.New("insert friend request error")
	}

	return nil
}


// 插入新消息
func insertMessage(m []byte) error {
	// json 解析 message 数据
	var msg model.Message
	err := json.Unmarshal(m, &msg)
	if err != nil {
		log.Error.Printf("unmarshal message err: %v\n", err)
	}

	// 添加消息
	messageDAO := dao.NewMessageDAO(variable.GoDB)
	_, err = messageDAO.AddMessage(msg)
	if err != nil {
		return errors.New("insert message error")
	}

	return nil
}

// 插入群组邀请
func insertGroupInvite(m []byte) error {
	// json 解析 group invite 数据
	var gi model.GroupInvite
	err := json.Unmarshal(m, &gi)
	if err != nil {
		log.Error.Printf("unmarshal group invite err: %v\n", err)
	}

	// 添加群组邀请
	groupInviteDAO := dao.NewGroupInviteDAO(variable.GoDB)
	_, err = groupInviteDAO.AddGroupInvite(gi)
	if err != nil {
		return errors.New("insert group invite error")
	}

	return nil
}