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

// 群组成员列表
func GroupMembers(m map[string]*flag.Flag, conn net.Conn) {
	IdWorker, _ := sequence.NewWorker(0)

	var g model.Group
	groupId, err := strconv.ParseInt(strings.TrimSpace(m["gn"].Value.String()), 10, 64)
	if err != nil {
		fmt.Println("parse group id err:", err)
	}

	g = model.Group{GroupId: groupId}

	// 序列化群组对象
	ret, err := json.Marshal(g)
	if err != nil {
		log.Info.Println("serialize group object err:", err)
	}

	// 封包群组对象消息并发送
	dp := packet.NewDataPack()
	msg, _ := dp.Pack(constants.TCPNetwork, constants.GroupMembersOption, packet.NewMessage(IdWorker.GetId(), ret))
	_, err = conn.Write(msg)
	if err != nil {
		log.Error.Println("client group list write message err:", err)
	}
}