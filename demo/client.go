package main

import (
	"bufio"
	"chalurania/comet/constants"
	"chalurania/comet/packet"
	"chalurania/demo/client/commands"
	"chalurania/demo/client/logic"
	"chalurania/service/log"
	"chalurania/service/model"
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
			case constants.SearchUsernameCommand:
				// 通过用户名进行搜索
				logic.Search(c.CommandMap, conn, constants.SearchUsernameCommand)
			case constants.SearchUserIdCommand:
				// 通过用户 id 进行搜索
				logic.Search(c.CommandMap, conn, constants.SearchUserIdCommand)
			case constants.AddUserByNameCommand:
				// 通过用户名添加好友
				logic.FriendRequest(c.CommandMap, conn, constants.AddUserByNameCommand)
			case constants.AddUserByIdCommand:
				// 通过用户 id 添加好友
				logic.FriendRequest(c.CommandMap, conn, constants.AddUserByIdCommand)
			case constants.FriendReqListCommand:
				// 好友请求列表
				logic.FriendReqList(conn)
			case constants.AcceptFriendByNameCommand:
				// 通过用户名接受好友请求
				logic.AcceptFriend(c.CommandMap, conn, constants.AcceptFriendByNameCommand)
			case constants.AcceptFriendByIdCommand:
				// 通过用户 id 接受好友请求
				logic.AcceptFriend(c.CommandMap, conn, constants.AcceptFriendByIdCommand)
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
				fmt.Printf("\b\blost connect from remote server, ")
				os.Exit(1)
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
				case constants.SignUpAckOpt:
					// 注册
					fmt.Printf("\b\b%s \n", ackPack.Data)
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
				case constants.SearchAckOpt:
					if ackPack.Sign {
						var user model.User
						err = json.Unmarshal(ackPack.Data, &user)
						if err != nil {
							log.Error.Printf("unmarshal user err: %v\n", err)
						}
						fmt.Printf("\b\b[id: %d, username: \"%s\", nickname: \"%s\"]\n", user.UserId, user.Username, user.Nickname)
					} else {
						fmt.Printf("\b\b%s \n", ackPack.Data)
					}
				case constants.FriendRequestAckOpt:
					fmt.Printf("\b\b%s \n", ackPack.Data)
				case constants.FriendReqListAckOpt:
					if ackPack.Sign {
						var frArray [2][]byte
						err = json.Unmarshal(ackPack.Data, &frArray)
						if err != nil {
							log.Error.Printf("unmarshal friend request list err: %v\n", err)
						}

						fmt.Printf("\b\bfriend request to: \n")
						var sentFrMap map[string][]byte
						err = json.Unmarshal(frArray[0], &sentFrMap)
						for _, value := range sentFrMap {
							var u model.User
							err = json.Unmarshal(value, &u)

							fmt.Println("\t", u.Username, "(", u.UserId, ")")
						}

						fmt.Printf("\b\bfriend request from: \n")
						var receivedFrMap map[string][]byte
						err = json.Unmarshal(frArray[1], &receivedFrMap)
						for _, value := range receivedFrMap {
							var u model.User
							err = json.Unmarshal(value, &u)

							fmt.Println("\t", u.Username, "(", u.UserId, ")")
						}


					} else {
						fmt.Printf("\b\b%s \n", ackPack.Data)
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
