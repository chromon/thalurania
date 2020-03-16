package config

import (
	"chalurania/api"
	"chalurania/service/log"
	"encoding/json"
	"io/ioutil"
)

// 服务器全局参数
type Global struct {
	// 全局 Server 对象
	TCPServer api.IServer

	// 服务器主机
	Host string

	// 服务器端口号
	Port int

	// 服务器名称
	Name string

	// 服务器版本号
	Version string

	// 数据包最大值
	MaxPacketSize uint32

	// 服务器允许的最大连接数
	MaxConn int
}

// 全局对象
var GlobalObj *Global

// 全局对象初始化
func init() {
	// 设置全局对象默认值
	GlobalObj = &Global{
		Name: "IM Server",
		Version: "v0.1",
		Host: "127.0.0.1",
		Port: 8080,
		MaxPacketSize: 4096,
		MaxConn: 10000,
	}

	// 从配置文件中重写加载用户自定义配置参数
	GlobalObj.Reload()
}

// 读取用户配置文件
func (g *Global) Reload() {
	data, err := ioutil.ReadFile("conf/config.json")
	if err != nil {
		log.Error.Println("Read config file err:", err)
	}

	// 解析 json 数据到对象中
	err = json.Unmarshal(data, &GlobalObj)
}