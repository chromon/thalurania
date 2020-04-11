package api

// 消息封装抽象接口
type IMessage interface {
	// 获取消息 Id
	GetMsgId() int64

	// 获取消息数据段长度
	GetDataLen() uint32

	// 获取消息内容
	GetData() []byte

	// 设置消息 Id
	SetMsgId(int64)

	// 设置消息内容
	SetData([]byte)

	// 设置消息数据段长度
	SetDataLen(uint32)
}