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

func GroupInvite(m map[string]*flag.Flag, conn net.Conn, opt int32) {
	IdWorker, _ := sequence.NewWorker(0)

	var u model.User

	// 创建用户对象
	switch opt {
	case constants.GroupInviteByNameCommand:
		// 搜索用户名, 使用 extra 字段暂代群组 id
		u = model.User{Username: m["u"].Value.String(), Extra: strings.TrimSpace(m["gn"].Value.String())}
	case constants.GroupInviteByIdCommand:
		// 搜索用户 id
		userId, err := strconv.ParseInt(strings.TrimSpace(m["n"].Value.String()), 10, 64)
		if err != nil {
			fmt.Println("parse user id err:", err)
		}
		u = model.User{UserId: userId, Extra: strings.TrimSpace(m["gn"].Value.String())}
	}

	// 序列化用户对象
	ret, err := json.Marshal(u)
	if err != nil {
		log.Info.Println("serialize user object err:", err)
	}

	serverTransPack := packet.NewServerTransPack(opt, []byte(ret))
	stp, err := json.Marshal(serverTransPack)
	if err != nil {
		log.Info.Println("serialize server trans pack (group invite) object err:", err)
		return
	}

	// 封包用户对象消息并发送
	dp := packet.NewDataPack()
	msg, _ := dp.Pack(constants.TCPNetwork, constants.GroupInviteOption, packet.NewMessage(IdWorker.GetId(), stp))
	_, err = conn.Write(msg)
	if err != nil {
		log.Error.Println("client friend request write message err:", err)
	}
}