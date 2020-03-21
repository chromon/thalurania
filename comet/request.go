package comet

import "chalurania/api"

// 客户端消息请求信息
type Request struct {
	// 请求 Id
	requestId uint32

	// 与客户端建立好的连接
	conn api.IConnection

	// 客户端请求的数据
	message api.IMessage
}

// 获取请求 Id
func (r *Request) GetRequestId() uint32 {
	return r.requestId
}

// 获取请求连接信息
func (r *Request) GetConnection() api.IConnection {
	return r.conn
}

// 获取请求消息数据
func (r *Request) GetData() []byte {
	return r.message.GetData()
}

// 得到请求的消息 ID
func (r *Request) GetMsgID() uint32 {
	return r.message.GetMsgId()
}