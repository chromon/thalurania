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
			case constants.SearchGroupCommand:
				// 搜索群组
				logic.Search(c.CommandMap, conn, constants.SearchGroupCommand)
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
			case constants.FriendListCommand:
				// 好友列表
				logic.FriendList(conn)
			case constants.ChatWithFriendByNameCommand:
				// 通过好友用户名聊天
				logic.PrivateChat(c.CommandMap, conn, constants.ChatWithFriendByNameCommand)
			case constants.ChatWithFriendByIdCommand:
				// 通过好友用户 id 聊天
				logic.PrivateChat(c.CommandMap, conn, constants.ChatWithFriendByIdCommand)
			case constants.OfflineMsgByNameCommand:
				// 通过用户名查询离线消息
				logic.OfflineMessage(c.CommandMap, conn, constants.OfflineMsgByNameCommand)
			case constants.OfflineMsgByIdCommand:
				// 通过用户 id 查询离线消息
				logic.OfflineMessage(c.CommandMap, conn, constants.OfflineMsgByIdCommand)
			case constants.CreateGroupCommand:
				// 创建群组
				logic.CreateGroup(conn)
			case constants.GroupInviteByNameCommand:
				// 通过用户名邀请到群组
				logic.GroupInvite(c.CommandMap, conn, constants.GroupInviteByNameCommand)
			case constants.GroupInviteByIdCommand:
				// 通过用户 id 邀请到群组
				logic.GroupInvite(c.CommandMap, conn, constants.GroupInviteByIdCommand)
			case constants.GroupInviteListCommand:
				// 群组邀请列表
				logic.GroupInviteList(conn)
			case constants.AcceptGroupInviteCommand:
				// 接受群组邀请
				logic.AcceptGroup(c.CommandMap, conn)
			case constants.GroupMemberCommand:
				// 群成员列表
				logic.GroupMembers(c.CommandMap, conn)
			case constants.GroupListCommand:
				// 已加入的群组列表
				logic.GroupList(conn)
			case constants.GroupChatCommand:
				// 群组聊天
				logic.GroupChat(c.CommandMap, conn)
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
				fmt.Println("unpack data header err:", err)
				return
			}

			if receiveMsg.GetDataLen() > 0 {
				msg := receiveMsg.(*packet.Message)
				msg.Data = make([]byte, msg.GetDataLen())

				// 读取消息内容
				_, err := io.ReadFull(conn, msg.Data)
				if err != nil {
					fmt.Println("client unpack data err:", err)
					return
				}

				// 解析 json 类型消息为 ack 包
				var ackPack packet.ServerAckPack
				err = json.Unmarshal(msg.Data, &ackPack)
				if err != nil {
					fmt.Printf("unmarshal ack pack err=%v\n", err)
				}

				switch ackPack.Opt {
				case constants.SignUpAckOpt:
					// 注册
					fmt.Printf("\b\b%s \n", ackPack.Data)
				case constants.LoginAckOpt:
					// 登录
					if ackPack.Sign {
						var offlineMsgMap map[string]string
						err = json.Unmarshal(ackPack.Data, &offlineMsgMap)
						//fmt.Println(offlineMsgMap)
						if len(offlineMsgMap) > 0 {
							for key, value := range offlineMsgMap {
								fmt.Printf("\b\b[offline] you have %s unread messages from %s\n", value, key)
							}
						}
					} else {
						fmt.Printf("\b\b%s \n", ackPack.Data)
					}
				case constants.DeviceOffline:
					// 被动离线
					fmt.Printf("\b\b%s \n", ackPack.Data)
					os.Exit(0)
				case constants.LogoutAckOpt:
					// 登出
					if ackPack.Sign {
						fmt.Printf("\b\b%s \n", ackPack.Data)
						os.Exit(0)
					} else {
						fmt.Printf("\b\b%s \n", ackPack.Data)
						os.Exit(1)
					}
				case constants.SearchAckOpt:
					// 搜索
					fmt.Printf("\b\b%s \n", ackPack.Data)
				case constants.FriendRequestAckOpt:
					// 添加好友请求
					fmt.Printf("\b\b%s \n", ackPack.Data)
				case constants.FriendReqListAckOpt:
					// 好友请求列表
					if ackPack.Sign {
						var frArray [2][]byte
						err = json.Unmarshal(ackPack.Data, &frArray)
						if err != nil {
							fmt.Printf("unmarshal friend request list err: %v\n", err)
						}

						fmt.Printf("\b\b  friend request to: \n")
						var sentFrMap map[string][]byte
						err = json.Unmarshal(frArray[0], &sentFrMap)
						for _, value := range sentFrMap {
							var u model.User
							err = json.Unmarshal(value, &u)

							fmt.Println("    ", u.Username, "(", u.UserId, ")")
						}

						fmt.Printf("\b\b  friend request from: \n")
						var receivedFrMap map[string][]byte
						err = json.Unmarshal(frArray[1], &receivedFrMap)
						for _, value := range receivedFrMap {
							var u model.User
							err = json.Unmarshal(value, &u)

							fmt.Println("    ", u.Username, "(", u.UserId, ")")
						}
					} else {
						fmt.Printf("\b\b%s \n", ackPack.Data)
					}
				case constants.AcceptFriendRepAckOpt:
					// 接受好友请求
					fmt.Printf("\b\b%s \n", ackPack.Data)
				case constants.FriendListAckOpt:
					// 好友列表
					if ackPack.Sign {
						var friendMap map[string][]byte
						err = json.Unmarshal(ackPack.Data, &friendMap)
						if err != nil {
							fmt.Printf("unmarshal friend map err: %v\n", err)
						}

						fmt.Printf("\b\b  friend list: \n")
						for _, value := range friendMap {
							var u model.User
							err = json.Unmarshal(value, &u)

							fmt.Println("    ", u.Username, "(", u.UserId, ")")
						}
					} else {
						fmt.Printf("\b\b%s \n", ackPack.Data)
					}
				case constants.SendMessageAckOpt:
					// 发送消息
					fmt.Printf("\b\b%s \n", ackPack.Data)
				case constants.OfflineMsgAckOpt:
					// 离线消息
					if ackPack.Sign {
						var offlineMsgMap map[string][]byte
						err = json.Unmarshal(ackPack.Data, &offlineMsgMap)

						for _, value := range offlineMsgMap {
							// 离线消息
							var message model.Message
							err = json.Unmarshal(value, &message)
							if err != nil {
								fmt.Printf("unmarshal messages err: %v\n", err)
							}

							fmt.Printf("\b\b\"%s\" -- %v\n", message.Content, message.CreateTime)
						}
					} else {
						fmt.Printf("\b\b%s \n", ackPack.Data)
					}
				case constants.CreateGroupAckOpt:
					// 创建群组
					fmt.Printf("\b\b%s \n", ackPack.Data)
				case constants.GroupRequestAckOpt:
					// 群组邀请
					fmt.Printf("\b\b%s \n", ackPack.Data)
				case constants.GroupInviteAckOpt:
					// 群组邀请列表
					if ackPack.Sign {
						res := strings.Split(string(ackPack.Data), ",")
						for i := 0; i < len(res) - 1; i++ {
							fmt.Printf("\b\b%s \n", res[i])
						}
					} else {
						fmt.Printf("\b\b%s \n", ackPack.Data)
					}
				case constants.AcceptGroupInviteAckOpt:
					// 接受群组邀请
					fmt.Printf("\b\b%s \n", ackPack.Data)
				case constants.GroupMembersAckOpt:
					// 群组成员列表
					if ackPack.Sign {
						res := strings.Split(string(ackPack.Data), ",")
						for i := 0; i < len(res) - 1; i++ {
							fmt.Printf("\b\b%s \n", res[i])
						}
					} else {
						fmt.Printf("\b\b%s \n", ackPack.Data)
					}
				case constants.GroupListAckOpt:
					// 已加入群组列表
					if ackPack.Sign {
						res := strings.Split(string(ackPack.Data), ",")
						for i := 0; i < len(res) - 1; i++ {
							fmt.Printf("\b\b%s \n", res[i])
						}
					} else {
						fmt.Printf("\b\b%s \n", ackPack.Data)
					}
				case constants.SendGroupMessageAckOpt:
					// 群组聊天
					fmt.Printf("\b\b%s \n", ackPack.Data)
				}

				fmt.Printf("\b\b~ ")
			}
		}
	}()

	for {
		time.Sleep(time.Second)
	}
}
