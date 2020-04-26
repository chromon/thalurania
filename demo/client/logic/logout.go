package logic

import (
	"chalurania/comet/constants"
	"chalurania/comet/packet"
	"chalurania/service/log"
	"chalurania/service/sequence"
	"net"
)

// 登出
func Logout(conn net.Conn) {
	IdWorker, _ := sequence.NewWorker(0)

	// 封包用户对象消息并发送
	dp := packet.NewDataPack()
	msg, _ := dp.Pack(constants.TCPNetwork, constants.LogoutOption,
		packet.NewMessage(IdWorker.GetId(), []byte("logout from server")))
	_, err := conn.Write(msg)
	if err != nil {
		log.Error.Println("client logout write message err:", err)
		return
	}
}