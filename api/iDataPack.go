package api

// 数据包
// 在 TCP 连接中的数据流，为数据添加头部信息
type IDataPack interface {

	// 获取包头长度
	GetHeaderLen() uint32

	// 封包
	Pack(IMessage) ([]byte, error)

	// 拆包
	Unpack([]byte) (IMessage, error)
}