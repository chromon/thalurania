package api

import "net"

// 连接接口
type IConnection interface {
	// 启动连接
	Start()

	// 停止连接
	Stop()

	// 获取当前连接 socket TCPConn
	GetTCPConnection() *net.TCPConn

	// 获取当前连接 Id
	GetConnId() uint32

	// 获取远程客户端地址信息
	GetRemoteAddr() net.Addr

	// 将 Message 数据发送到远程 TCP 客户端
	SendMsg(uint32, []byte) error
}