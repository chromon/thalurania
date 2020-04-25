package packet

// 回执信息
type ServerAckPack struct {
	// ack 指令
	Opt int32 `json:"opt"`

	// 标记符号，返回客户端操作命令是否成功
	Sign bool `json:"sign"`

	// 消息数据
	Data []byte `json:"data"`
}

func NewServerAckPack(opt int32, sign bool, data []byte) *ServerAckPack {
	return &ServerAckPack{
		Opt: opt,
		Sign: sign,
		Data: data,
	}
}