package logic

import (
	"chalurania/comet/constants"
	"chalurania/comet/packet"
	"chalurania/service/log"
	"chalurania/service/model"
	"chalurania/service/sequence"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"
)

func Search(m map[string]*flag.Flag, conn net.Conn, opt int32) {
	IdWorker, _ := sequence.NewWorker(0)

	var u model.User

	// 创建用户对象
	switch opt {
	case constants.SearchUsernameCommand:
		// 搜索用户名
		u = model.User{Username: m["u"].Value.String()}
	case constants.SearchUserIdCommand:
		// 搜索用户 id
		userId, err := strconv.ParseInt(strings.TrimSpace(m["n"].Value.String()), 10, 64)
		if err != nil {
			fmt.Println("parse user id err:", err)
		}
		u = model.User{UserId: userId}
	}

	// 序列化用户对象
	ret, err := json.Marshal(u)
	if err != nil {
		log.Info.Println("serialize user object err:", err)
	}

	serverTransPack := packet.NewServerTransPack(opt, []byte(ret))
	stp, err := json.Marshal(serverTransPack)
	if err != nil {
		log.Info.Println("serialize server trans pack (search) object err:", err)
		return
	}

	// 封包用户对象消息并发送
	dp := packet.NewDataPack()
	msg, _ := dp.Pack(constants.TCPNetwork, constants.SearchOption, packet.NewMessage(IdWorker.GetId(), stp))
	_, err = conn.Write(msg)
	if err != nil {
		log.Error.Println("client search write message err:", err)
	}
}