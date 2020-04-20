package consumers

import (
	"chalurania/comet/packet"
	"chalurania/comet/variable"
	"chalurania/service/dao"
	"chalurania/service/model"
	"chalurania/service/scrypt"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
)

func Consume(msg redis.Message) error {
	// json 解析 data wrap 数据
	var dw packet.DataWrap
	err := json.Unmarshal(msg.Data, &dw)
	if err != nil {
		fmt.Printf("unmarshal dw err=%v\n", err)
	}

	switch dw.Opt {
	case 1:
		// json 解析 user 数据
		var u model.User
		err = json.Unmarshal(dw.Model, &u)
		if err != nil {
			fmt.Printf("unmarshal user err=%v\n", err)
		}
		// 加密密码
		u.Password = scrypt.Crypto(u.Password)

		// 添加用户
		userDAO := dao.NewUserDAO(variable.GoDB)
		_, err := userDAO.AddUser(u)
		if err != nil {
			return errors.New("insert user err")
		}
	}

	return nil
}
