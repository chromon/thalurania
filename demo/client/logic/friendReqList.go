package logic

import (
	"chalurania/comet/constants"
	"chalurania/comet/packet"
	"chalurania/service/log"
	"chalurania/service/sequence"
	"net"
)

// 查询好友请求列表
func FriendReqList(conn net.Conn) {
	IdWorker, _ := sequence.NewWorker(0)

	dp := packet.NewDataPack()
	msg, _ := dp.Pack(constants.TCPNetwork, constants.FriendReqListOption,
		packet.NewMessage(IdWorker.GetId(), []byte("query friend request list")))
	_, err := conn.Write(msg)
	if err != nil {
		log.Error.Println("client logout write message err:", err)
		return
	}
}