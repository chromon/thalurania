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
// 注册：tim -r -u 用户名 -p
// 登录：tim -l -u 用户名 -p
// 登出：tim -q
// 搜索：tim -s -u 用户名
// 		tim -s -n 用户 id
// 		tim -s -g 用户组 id
// 好友请求：tim -add -u 用户名
//		   tim -add -n 用户 id
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
	// 用户 id
	c.FlagSet.String("n", "", "account id")
	// 请求添加好友
	c.FlagSet.Bool("add", false, "request to be friends")
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
