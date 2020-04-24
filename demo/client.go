package main

import (
	"bufio"
	"chalurania/comet/packet"
	"chalurania/demo/client/commands"
	"chalurania/demo/client/logic"
	"chalurania/service/log"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"
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

	// 读取命令并发送到服务器
	go func() {
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
				fmt.Println("commands not found")
			case 1:
				// 注册命令
				logic.SignUp(c.CommandMap, conn)
			case 2:
				// 登录命令
				logic.Login(c.CommandMap, conn)
			}
		}
	}()

	// 接收服务器发送的消息
	go func() {
		for {
			dp := packet.NewDataPack()
			// 读取流中的消息回执 ack 数据包 header 部分
			header := make([]byte, dp.GetHeaderLen())
			_, err = io.ReadFull(conn, header)
			if err != nil {
				log.Error.Println("client read ack header err:", err)
				return
			}

			// ack 拆包
			_, _, receiveMsg, err := dp.Unpack(header)
			if err != nil {
				log.Error.Println("unpack data header err:", err)
				return
			}

			if receiveMsg.GetDataLen() > 0 {
				msg := receiveMsg.(*packet.Message)
				msg.Data = make([]byte, msg.GetDataLen())

				_, err := io.ReadFull(conn, msg.Data)
				if err != nil {
					log.Error.Println("client unpack data err:", err)
					return
				}
				//fmt.Println(string(msg.Data))
				fmt.Printf("\b\bServer feedback: %s \n", string(msg.Data))
				fmt.Print("~ ")
				//os.Exit(0)
			}
		}
	}()

	for {
		time.Sleep(time.Second)
	}

	//for {
	//	// 创建命令并初始化
	//	c := commands.NewCommand("tim")
	//	c.CommandInit()
	//
	//	// 命令提示符
	//	fmt.Print("~ ")
	//	// 读取命令
	//	scanner.Scan()
	//	if err := scanner.Err(); err != nil {
	//		_, err = fmt.Fprintln(os.Stderr, "error:", err)
	//	}
	//	args = scanner.Text()
	//
	//	// 解析命令
	//	c.ParseCommand(strings.Split(args, " ")[1:])
	//	c.VisitCommand()
	//
	//	// 命令分发
	//	switch commands.CommandDistribute(c.CommandMap) {
	//	case 0:
	//		fmt.Println("commands not found")
	//	case 1:
	//		// 注册命令
	//		logic.SignUp(c.CommandMap, conn)
	//	case 2:
	//		// 登录命令
	//		s, err := logic.Login(c.CommandMap, conn)
	//		if err != nil {
	//			fmt.Println("login err:", err)
	//		}
	//		sign = s
	//	}
	//
	//	//if sign {
	//	//	sign = false
	//	//	// 读取信息
	//	//	go func() {
	//	//		for {
	//	//			dp := packet.NewDataPack()
	//	//			// 读取流中的消息回执 ack 数据包 header 部分
	//	//			header := make([]byte, dp.GetHeaderLen())
	//	//			_, err = io.ReadFull(conn, header)
	//	//			if err != nil {
	//	//				log.Error.Println("client read ack header err:", err)
	//	//				return
	//	//			}
	//	//
	//	//			// ack 拆包
	//	//			_, _, receiveMsg, err := dp.Unpack(header)
	//	//			if err != nil {
	//	//				log.Error.Println("unpack data header err:", err)
	//	//				return
	//	//			}
	//	//
	//	//			if receiveMsg.GetDataLen() > 0 {
	//	//				msg := receiveMsg.(*packet.Message)
	//	//				msg.Data = make([]byte, msg.GetDataLen())
	//	//
	//	//				_, err := io.ReadFull(conn, msg.Data)
	//	//				if err != nil {
	//	//					log.Error.Println("client unpack data err:", err)
	//	//					return
	//	//				}
	//	//				fmt.Println(string(msg.Data))
	//	//				os.Exit(0)
	//	//			}
	//	//		}
	//	//	}()
	//	//}
	//}
}
