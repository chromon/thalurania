package main

import (
	"bufio"
	"chalurania/comet/constants"
	"chalurania/comet/packet"
	"chalurania/demo/client/commands"
	"chalurania/demo/client/logic"
	"chalurania/service/log"
	"encoding/json"
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
			case constants.ErrorCommand:
				fmt.Println("commands not found")
			case constants.RegisterCommand:
				// 注册命令
				logic.SignUp(c.CommandMap, conn)
			case constants.LoginCommand:
				// 登录命令
				logic.Login(c.CommandMap, conn)
			case constants.LogoutCommand:
				// 登出命令
				logic.Logout(conn)
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

			// 将服务端返回的 ack header 信息拆包
			_, _, receiveMsg, err := dp.Unpack(header)
			if err != nil {
				log.Error.Println("unpack data header err:", err)
				return
			}

			if receiveMsg.GetDataLen() > 0 {
				msg := receiveMsg.(*packet.Message)
				msg.Data = make([]byte, msg.GetDataLen())

				// 读取消息内容
				_, err := io.ReadFull(conn, msg.Data)
				if err != nil {
					log.Error.Println("client unpack data err:", err)
					return
				}

				// 解析 json 类型消息为 ack 包
				var ackPack packet.ServerAckPack
				err = json.Unmarshal(msg.Data, &ackPack)
				if err != nil {
					log.Error.Printf("unmarshal ack pack err=%v\n", err)
				}

				switch ackPack.Opt {
				case constants.LoginAckOpt:
					// 登录
					fmt.Printf("\b\b%s \n", ackPack.Data)
				case constants.DeviceOffline:
					// 被动离线
					fmt.Printf("\b\b%s \n", ackPack.Data)
					os.Exit(0)
				case constants.LogoutAckOpt:
					if ackPack.Sign {
						fmt.Printf("\b\b%s \n", ackPack.Data)
						os.Exit(0)
					} else {
						fmt.Printf("\b\b%s \n", ackPack.Data)
						os.Exit(1)
					}
				}

				fmt.Print("~ ")
			}
		}
	}()

	for {
		time.Sleep(time.Second)
	}
}
