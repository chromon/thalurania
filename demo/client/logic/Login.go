package logic

import (
	"bufio"
	"chalurania/comet/packet"
	"chalurania/service/log"
	"chalurania/service/model"
	"chalurania/service/sequence"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

// 登录
func Login(m map[string]*flag.Flag, conn net.Conn) {
	scanner := bufio.NewScanner(os.Stdin)
	IdWorker, _ := sequence.NewWorker(0)

	// 读取密码
	fmt.Print("~ ")
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		_, err = fmt.Fprintln(os.Stderr, "error:", err)
	}
	pwd := scanner.Text()

	// 创建用户对象
	u := model.User{Nickname: m["u"].Value.String(), Password: pwd}
	// 序列化用户对象
	ret, err := json.Marshal(u)
	if err != nil {
		log.Info.Println("serialize user object err:", err)
		return
	}

	// 封包用户对象消息并发送
	dp := packet.NewDataPack()
	msg, _ := dp.Pack(1, 2, packet.NewMessage(IdWorker.GetId(), ret))
	_, err = conn.Write(msg)
	if err != nil {
		log.Error.Println("client login write message err:", err)
		return
	}

	// 读取流中的消息回执 ack 数据包 header 部分
	header := make([]byte, dp.GetHeaderLen())
	_, err = io.ReadFull(conn, header)
	if err != nil {
		log.Error.Println("client login read ack header err:", err)
		return
	}

	// ack 拆包
	_, _, receiveMsg, err := dp.Unpack(header)
	if err != nil {
		log.Error.Println("unpack login ack header err:", err)
		return
	}

	if receiveMsg.GetDataLen() > 0 {
		msg := receiveMsg.(*packet.Message)
		msg.Data = make([]byte, msg.GetDataLen())

		_, err := io.ReadFull(conn, msg.Data)
		if err != nil {
			log.Error.Println("client unpack login ack data err:", err)
			return
		}
		fmt.Println(string(msg.Data))
	}
}