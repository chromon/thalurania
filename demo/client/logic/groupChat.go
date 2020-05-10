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

// 群组聊天
func GroupChat(m map[string]*flag.Flag, conn net.Conn) {

	IdWorker, _ := sequence.NewWorker(0)

	groupId, err := strconv.ParseInt(strings.TrimSpace(m["gn"].Value.String()), 10, 64)
	if err != nil {
		fmt.Println("parse group id err:", err)
	}
	message := model.Message{ReceiverId: groupId, Content: m["m"].Value.String()}

	// 序列化消息对象
	ret, err := json.Marshal(message)
	if err != nil {
		log.Info.Println("serialize message object err:", err)
	}

	// 封包用户对象消息并发送
	dp := packet.NewDataPack()
	msg, _ := dp.Pack(constants.TCPNetwork, constants.GroupChatOption, packet.NewMessage(IdWorker.GetId(), ret))
	_, err = conn.Write(msg)
	if err != nil {
		log.Error.Println("client group chat write message err:", err)
	}
}