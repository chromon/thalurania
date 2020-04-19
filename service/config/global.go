package config

import (
	"chalurania/api"
	"chalurania/service/log"
	"encoding/json"
	"io/ioutil"
	"os"
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

	// 业务工作池的数量
	WorkerPoolSize uint32

	// 业务工作池对应的任务队列最大任务数
	MaxWorkerTaskLen uint32

	// SendBufMsg 发送消息的缓冲最大长度
	MaxMsgChanLen uint32

	// 配置文件路径
	ConfigFilePath string

	// redis network
	RedisNetwork string

	// redis address
	RedisAddress string

	// redis password
	RedisPassword string

	// redis database
	RedisDatabase int

	// mysql username
	DBUserName string

	// mysql password
	DBPassword string

	// mysql host
	DBHost string

	// mysql port
	DBPort string

	// mysql database name
	DBName string
}

// 全局对象
var GlobalObj *Global

// 判断一个文件是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// 读取用户配置文件
func (g *Global) Reload() {
	if exist, _ := PathExists(g.ConfigFilePath); exist == false {
		log.Error.Println("Config file", g.ConfigFilePath, "isn't exist")
		return
	}

	data, err := ioutil.ReadFile(g.ConfigFilePath)
	if err != nil {
		log.Error.Println("Read config file err:", err)
	}

	// 解析 json 数据到对象中
	err = json.Unmarshal(data, g)
	if err != nil {
		log.Error.Println("JSON unmarshal err:", err)
	}
}

// 全局对象初始化
func init() {
	// 设置全局对象默认值
	GlobalObj = &Global{
		Name:             "IM Server",
		Version:          "v0.1",
		Host:             "127.0.0.1",
		Port:             8080,
		MaxPacketSize:    4096,
		MaxConn:          10000,
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
		MaxMsgChanLen:    1024,
		ConfigFilePath:   "conf/config.json",

		RedisNetwork:  "tcp",
		RedisAddress:  "127.0.0.1:6379",
		RedisPassword: "",
		RedisDatabase: 0,

		DBUserName: "root",
		DBPassword: "root",
		DBHost:     "127.0.0.1",
		DBPort:     "3306",
		DBName:     "thalurania",
	}

	// 从配置文件中重写加载用户自定义配置参数
	//GlobalObj.Reload()
}
