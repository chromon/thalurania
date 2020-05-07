package logic

import (
	"chalurania/comet/constants"
	"chalurania/comet/packet"
	"chalurania/service/log"
	"chalurania/service/sequence"
	"net"
)

// 群组请求列表
func GroupInviteList(conn net.Conn) {
	IdWorker, _ := sequence.NewWorker(0)

	// 封包用户对象消息并发送
	dp := packet.NewDataPack()
	msg, _ := dp.Pack(constants.TCPNetwork, constants.GroupInviteListOption,
		packet.NewMessage(IdWorker.GetId(), []byte("query group invite list from server")))
	_, err := conn.Write(msg)
	if err != nil {
		log.Error.Println("client query group invite list write message err:", err)
		return
	}
}