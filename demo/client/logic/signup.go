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

// 注册
func SignUp(m map[string]*flag.Flag, conn net.Conn) {
	IdWorker, _ := sequence.NewWorker(0)
	scanner := bufio.NewScanner(os.Stdin)
	var pwd string

	for {
		// 读取密码
		fmt.Print("~ ")
		scanner.Scan()
		if err := scanner.Err(); err != nil {
			_, err = fmt.Fprintln(os.Stderr, "error:", err)
		}
		pwd = scanner.Text()

		// 读取确认密码
		fmt.Print("~ ")
		scanner.Scan()
		if err := scanner.Err(); err != nil {
			_, err = fmt.Fprintln(os.Stderr, "error:", err)
		}
		cpwd := scanner.Text()

		if pwd == cpwd {
			break
		} else {
			log.Info.Println("the passwords you entered did not match, try again")
		}
	}

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
	msg, _ := dp.Pack(1, 1, packet.NewMessage(IdWorker.GetId(), ret))
	_, err = conn.Write(msg)
	if err != nil {
		log.Error.Println("client sign up write message err:", err)
		return
	}

	// 读取流中的消息回执 ack 数据包 header 部分
	header := make([]byte, dp.GetHeaderLen())
	_, err = io.ReadFull(conn, header)
	if err != nil {
		log.Error.Println("client sign up read ack header err:", err)
		return
	}

	// ack 拆包
	_, _, receiveMsg, err := dp.Unpack(header)
	if err != nil {
		log.Error.Println("unpack sign up ack header err:", err)
		return
	}

	if receiveMsg.GetDataLen() > 0 {
		msg := receiveMsg.(*packet.Message)
		msg.Data = make([]byte, msg.GetDataLen())

		_, err := io.ReadFull(conn, msg.Data)
		if err != nil {
			log.Error.Println("client unpack sign up ack data err:", err)
			return
		}
		fmt.Println(string(msg.Data))
	}
}