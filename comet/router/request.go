package router

import "chalurania/api"

// 客户端消息请求信息
type Request struct {
	// 请求 Id
	RequestId uint32

	// 与客户端建立好的连接
	Conn api.IConnection

	// 客户端请求的数据
	Message api.IMessage
}

// 获取请求 Id
func (r *Request) GetRequestId() uint32 {
	return r.RequestId
}

// 获取请求连接信息
func (r *Request) GetConnection() api.IConnection {
	return r.Conn
}

// 获取请求消息数据
func (r *Request) GetData() []byte {
	return r.Message.GetData()
}

// 得到请求的消息 ID
func (r *Request) GetMsgID() int64 {
	return r.Message.GetMsgId()
}