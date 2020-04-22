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
func (c *Command) CommandInit() {
	// 注册
	c.FlagSet.Bool("r", false, "register an account")
	// 账户名称
	c.FlagSet.String("u", "", "account name")
	// 账户密码
	c.FlagSet.Bool("p", false, "account password")
	// 登录
	c.FlagSet.Bool("l", false, "login an account")
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
