package api

// 封装客户端请求信息
type IRequest interface {
	// 获取请求 Id
	GetRequestId() uint32

	// 获取请求连接信息
	GetConnection() IConnection

	// 获取请求消息的数据
	GetData() []byte

	// 得到请求的消息 ID
	GetMsgID() int64
}