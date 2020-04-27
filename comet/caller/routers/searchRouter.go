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
type SearchRouter struct {
	router.Router
	success bool
	user model.User
	info string
}

func (sr *SearchRouter) Handle(r api.IRequest) {
	var stp packet.ServerTransPack
	err := json.Unmarshal(r.GetData(), &stp)
	if err != nil {
		log.Error.Printf("unmarshal server trans pack err: %v\n", err)
	}

	// 用户信息
	var user model.User
	err = json.Unmarshal(stp.Data, &user)
	if err != nil {
		log.Error.Printf("unmarshal user err: %v\n", err)
	}

	switch stp.Opt {
	case constants.SearchUsernameCommand:
		// 搜索用户名
		// 搜索失败信息
		sr.info = "username - '" + user.Username + "'"

		// 校验用户信息
		userDAO := dao.NewUserDAO(variable.GoDB)
		exist, u := userDAO.QueryUserByName(user)
		if exist {
			sr.success = true
			sr.user = *u
		} else {
			sr.success = false
		}
	case constants.SearchUserIdCommand:
		// 搜索用户 id
		// 搜索失败信息
		sr.info = "userId - '" + strconv.FormatInt(user.UserId,10) + "'"

		// 校验用户信息
		userDAO := dao.NewUserDAO(variable.GoDB)
		exist, u := userDAO.QueryUserById(user)
		if exist {
			sr.success = true
			sr.user = *u
		} else {
			sr.success = false
		}
	}
}

// 回执消息
func (sr *SearchRouter) PostHandle(r api.IRequest) {
	// 反向客户端发送 ack 数据
	var searchMsg []byte
	if sr.success {
		// 序列化用户对象
		ret, err := json.Marshal(sr.user)
		if err != nil {
			log.Info.Println("serialize user object err:", err)
		}
		searchMsg = ret
	} else {
		searchMsg = []byte("oops, the search info " + sr.info +" not exist")
	}

	// 包装 ack
	ackPack := packet.NewServerAckPack(constants.SearchAckOpt, sr.success, searchMsg)
	ret, err := json.Marshal(ackPack)
	if err != nil {
		log.Info.Println("serialize search ack pack object err:", err)
		return
	}

	// 发送回执
	err = r.GetConnection().SendMsg(constants.TCPNetwork, constants.AckOption, 101, ret)
	if err != nil {
		log.Error.Println("search send ack message to client err:", err)
	}
}