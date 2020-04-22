package consumers

import (
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
	var dw packet.DataWrap
	err := json.Unmarshal(msg.Data, &dw)
	if err != nil {
		fmt.Printf("unmarshal data wrap err=%v\n", err)
	}

	switch dw.Opt {
	case 1:
		// 插入新用户
		err = insertUser(dw.Model)
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
		log.Error.Printf("unmarshal user err=%v\n", err)
	}
	// 加密密码
	u.Password = scrypt.Crypto(u.Password)

	// 添加用户
	userDAO := dao.NewUserDAO(variable.GoDB)
	_, err = userDAO.AddUser(u)
	if err != nil {
		return errors.New("insert user err")
	}

	return nil
}
