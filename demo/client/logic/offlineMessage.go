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

// 查询离线消息
func OfflineMessage(m map[string]*flag.Flag, conn net.Conn, opt int32) {
	IdWorker, _ := sequence.NewWorker(0)

	var u model.User

	// 创建用户对象
	switch opt {
	case constants.OfflineMsgByNameCommand:
		// 搜索用户名
		u = model.User{Username: m["u"].Value.String()}
	case constants.OfflineMsgByIdCommand:
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
		log.Info.Println("serialize server trans pack (offline message) object err:", err)
		return
	}

	// 封包用户对象消息并发送
	dp := packet.NewDataPack()
	msg, _ := dp.Pack(constants.TCPNetwork, constants.OfflineMessageOption, packet.NewMessage(IdWorker.GetId(), stp))
	_, err = conn.Write(msg)
	if err != nil {
		log.Error.Println("client offline message write message err:", err)
	}
}