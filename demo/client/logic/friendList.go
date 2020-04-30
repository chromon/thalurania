package logic

import (
	"chalurania/comet/constants"
	"chalurania/comet/packet"
	"chalurania/service/log"
	"chalurania/service/sequence"
	"net"
)

// 好友列表
func FriendList(conn net.Conn) {
	IdWorker, _ := sequence.NewWorker(0)

	// 封包用户对象消息并发送
	dp := packet.NewDataPack()
	msg, _ := dp.Pack(constants.TCPNetwork, constants.FriendListOption,
		packet.NewMessage(IdWorker.GetId(), []byte("request friend list from server")))
	_, err := conn.Write(msg)
	if err != nil {
		log.Error.Println("client friend list write message err:", err)
		return
	}
}