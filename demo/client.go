package main

import (
	"chalurania/comet"
	"chalurania/service/log"
	"io"
	"net"
	"time"
)

func main() {
	log.Info.Println("Client start")

	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Error.Println("Net dial err:", err)
		return
	}

	for {
		// 发送封包消息
		dp := comet.NewDataPack()
		msg, _ := dp.Pack(comet.NewMessage(1, []byte("First message to server1")))
		_, err := conn.Write(msg)
		if err != nil {
			log.Error.Println("Client write message err:", err)
			return
		}

		// 读取流中的数据包 header 部分
		header := make([]byte, dp.GetHeaderLen())
		_, err = io.ReadFull(conn, header)
		if err != nil {
			log.Error.Println("Client read header err:", err)
			break
		}

		// 拆包
		receiveMsg, err := dp.Unpack(header)
		if err != nil {
			log.Error.Println("Unpack err:", err)
			return
		}

		if receiveMsg.GetDataLen() > 0 {
			msg := receiveMsg.(*comet.Message)
			msg.Data = make([]byte, msg.GetDataLen())

			_, err := io.ReadFull(conn, msg.Data)
			if err != nil {
				log.Error.Println("Server unpack data err:", err)
				return
			}
			log.Info.Printf("Server feedback message id: %d - %s, len: %d", msg.Id, msg.Data, msg.DataLen)
		}

		time.Sleep(time.Second)
	}
}