package logic

import (
	"chalurania/comet/constants"
	"chalurania/comet/packet"
	"chalurania/service/log"
	"chalurania/service/sequence"
	"net"
)

// 已加入的群组列表
func GroupList(conn net.Conn) {
	IdWorker, _ := sequence.NewWorker(0)

	// 封包用户对象消息并发送
	dp := packet.NewDataPack()
	msg, _ := dp.Pack(constants.TCPNetwork, constants.GroupListOption,
		packet.NewMessage(IdWorker.GetId(), []byte("query group list from server")))
	_, err := conn.Write(msg)
	if err != nil {
		log.Error.Println("client query group list write message err:", err)
		return
	}
}