package arguments

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
func (c *Command) CommandInit() {
	c.FlagSet.Int("p", 8080, "server port")
	c.FlagSet.String("s", "aaa", "string value")
	c.FlagSet.Bool("b", false, "open file")
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
		fmt.Println("--", f)
		fmt.Println("--", f.Name)
		fmt.Println("--", f.Value.String())

		c.CommandMap[f.Name] = f
	}

	c.FlagSet.Visit(fn)
}
