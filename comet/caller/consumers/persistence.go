package consumers

import (
	"chalurania/comet/packet"
	"chalurania/comet/variable"
	"chalurania/service/dao"
	"chalurania/service/log"
	"chalurania/service/model"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
)

func Consume(msg redis.Message) error {
	log.Info.Println("hehehe")
	log.Info.Printf("recv msg: %s", msg.Data)

	// json 解析 data wrap 数据
	var dw packet.DataWrap
	err := json.Unmarshal(msg.Data, &dw)
	if err != nil {
		fmt.Printf("unmarshal dw err=%v\n", err)
	}
	fmt.Printf("dw=%v dw.Opt=%v dw.Model=%v \n", dw, dw.Opt, dw.Model)

	switch dw.Opt {
	case 1:
		// json 解析 user 数据
		var u model.User
		err = json.Unmarshal(dw.Model, &u)
		if err != nil {
			fmt.Printf("unmarshal user err=%v\n", err)
		}
		fmt.Printf("u=%v u.nickname=%v u.password=%v\n", u, u.Nickname, u.Password)

		userDAO := dao.NewUserDAO(variable.GoDB)
		_, err := userDAO.AddUser(u)
		if err != nil {
			return errors.New("insert user err")
		}
	}

	return nil
}
