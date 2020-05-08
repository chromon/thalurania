package commands

import (
	"flag"
	"fmt"
)

// 命令
type Command struct {
	// FlagSet 对象
	FlagSet *flag.FlagSet
	// 实际的命令集合
	CommandMap map[string]*flag.Flag
}

// 新建命令
func NewCommand(name string) *Command {
	c := &Command{
		FlagSet: flag.NewFlagSet(name, 0),
		CommandMap: make(map[string]*flag.Flag),
	}

	return c
}

// 初始化命令
// 注册：tim -r -u [用户名] -p
// 登录：tim -l -u [用户名] -p
// 登出：tim -q
// 搜索：tim -s -u [用户名]
// 		tim -s -n [用户 id]
// 		tim -s -g -gn [用户组 id]
// 好友请求：tim -add -u [用户名]
//		   tim -add -n [用户 id]
// 好友请求列表：tim -fr -list
// 接受好友请求：tim -accept -u [用户名]
//			   tim -accept -n [用户 id]
// 好友列表：tim -f -list
// 与好友聊天：tim -chat -u [用户名] -m [消息内容]
// 			 tim -chat -n [用户 id] -m [消息内容]
// 查看离线消息：tim -o -u [用户名]
//			   tim -o -n [用户 id]
// 创建群组：tim -add -g
// 群组邀请：tim -i -u [用户名] -gn [群组 id]
// 		   tim -i -n [用户 id] -gn [群组 id]
// 群组邀请列表：tim -g -i -list
// 接受群组邀请：tim -accept -gn [群组 id]
// 加入的群组列表：tim -g -list
// 群成员列表：tim -g -gn [群组 id] -list
func (c *Command) CommandInit() {
	// 注册
	c.FlagSet.Bool("r", false, "register an account")
	// 账户名称
	c.FlagSet.String("u", "", "account name")
	// 账户密码
	c.FlagSet.Bool("p", false, "account password")
	// 登录
	c.FlagSet.Bool("l", false, "logic an account")
	// 登出
	c.FlagSet.Bool("q", false, "logout an account")
	// 搜索
	c.FlagSet.Bool("s", false, "search")
	// 用户
	c.FlagSet.String("n", "", "account id")
	// 请求添加好友
	c.FlagSet.Bool("add", false, "request to be friends or group")
	// 请求好友
	c.FlagSet.Bool("req", false, "friend request")
	// 列表
	c.FlagSet.Bool("list", false, "list")
	// 接受好友请求
	c.FlagSet.Bool("accept", false, "accept request")
	// 好友
	c.FlagSet.Bool("f", false, "friend")
	// 聊天
	c.FlagSet.Bool("chat", false, "chat with friend or group")
	// 消息内容
	c.FlagSet.String("m", "", "message")
	// 查看离线消息
	c.FlagSet.Bool("o", false, "offline message")
	// 群组
	c.FlagSet.Bool("g", false, "group")
	// 群组 id
	c.FlagSet.String("gn", "", "group id")
	// 群组邀请
	c.FlagSet.Bool("i", false, "group invite")
}

// 解析命令
func (c *Command) ParseCommand(arguments []string) {
	err := c.FlagSet.Parse(arguments)
	if err != nil {
		fmt.Println("parse command err:", err)
	}
}

// 遍历命令并将其添加到实际的命令集合中
func (c *Command) VisitCommand() {
	fn := func (f *flag.Flag) {
		c.CommandMap[f.Name] = f
	}

	c.FlagSet.Visit(fn)
}
