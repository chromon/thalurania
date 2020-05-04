package logic

import (
	"chalurania/comet/constants"
	"chalurania/comet/packet"
	"chalurania/service/log"
	"chalurania/service/sequence"
	"net"
)

// 创建群组
func CreateGroup(conn net.Conn) {
	IdWorker, _ := sequence.NewWorker(0)

	// 封包用户对象消息并发送
	dp := packet.NewDataPack()
	msg, _ := dp.Pack(constants.TCPNetwork, constants.CreateGroupOption,
		packet.NewMessage(IdWorker.GetId(), []byte("create group")))
	_, err := conn.Write(msg)
	if err != nil {
		log.Error.Println("client create group write message err:", err)
		return
	}
}