package main

import (
	"bufio"
	"chalurania/demo/client/commands"
	"chalurania/demo/client/logic"
	"chalurania/service/log"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// 创建连接
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Error.Println("Net dial err:", err)
		return
	}

	// scanner 用户读取客户端命令
	scanner := bufio.NewScanner(os.Stdin)
	// 命令参数
	var args string

	fmt.Println("welcome to thalurania im")

	for {
		// 创建命令并初始化
		c := commands.NewCommand("tim")
		c.CommandInit()

		// 命令提示符
		fmt.Print("~ ")
		// 读取命令
		scanner.Scan()
		if err := scanner.Err(); err != nil {
			_, err = fmt.Fprintln(os.Stderr, "error:", err)
		}
		args = scanner.Text()

		// 解析命令
		c.ParseCommand(strings.Split(args, " ")[1:])
		c.VisitCommand()

		// 命令分发
		switch commands.CommandDistribute(c.CommandMap) {
		case 0:
			fmt.Println("combined commands not found")
		case 1:
			// 注册命令
			logic.SignUp(c.CommandMap, conn)
		}
	}










	//for {
	//
	//
	//
	//
	//	// 发送封包消息
	//	dp := packet.NewDataPack()
	//	msg, _ := dp.Pack(1, 1, packet.NewMessage(102, []byte("First message to server1")))
	//	_, err := conn.Write(msg)
	//	if err != nil {
	//		log.Error.Println("Client write message err:", err)
	//		return
	//	}
	//
	//	// 读取流中的数据包 header 部分
	//	header := make([]byte, dp.GetHeaderLen())
	//	_, err = io.ReadFull(conn, header)
	//	if err != nil {
	//		log.Error.Println("Client read header err:", err)
	//		break
	//	}
	//
	//	// 拆包
	//	network, operation, receiveMsg, err := dp.Unpack(header)
	//	if err != nil {
	//		log.Error.Println("Unpack err:", err)
	//		return
	//	}
	//
	//	if receiveMsg.GetDataLen() > 0 {
	//		msg := receiveMsg.(*packet.Message)
	//		msg.Data = make([]byte, msg.GetDataLen())
	//
	//		_, err := io.ReadFull(conn, msg.Data)
	//		if err != nil {
	//			log.Error.Println("Server unpack data err:", err)
	//			return
	//		}
	//		log.Info.Printf("Server feedback message id: %d - %s, len: %d, network: %d, opertion: %d", msg.Id, msg.Data, msg.DataLen, network, operation)
	//	}
	//
	//	time.Sleep(time.Second)
	//}
	//
	//
	//g := Girl{"satori", 16, "f"}
	//
	////ret, err := json.MarshalIndent(g, "", " ")
	//ret, err := json.Marshal(g)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//} else {
	//	fmt.Println(string(ret))
	//}
	//
	////创建一个变量
	//g2 := Girl{}
	////传入json字符串，和指针
	//err = json.Unmarshal(ret, &g2)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(g2)  //{satori 16 f 东方地灵殿 false}
	//fmt.Println(g2.Name, g2.Age) // satori 16
}
