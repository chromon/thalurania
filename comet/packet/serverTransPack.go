package packet

// 数据传输包，用于服务器中各系统之间信息传输，主要用在 PubSub 订阅消息传输，复用 user consumer
type ServerTransPack struct {
	// 频道中传输消息指令（类型）
	Opt int32 `json:"opt"`

	// 消息数据
	Data []byte `json:"data"`
}

func NewServerTransPack(opt int32, data []byte) *ServerTransPack {
	return &ServerTransPack{
		Opt: opt,
		Data: data,
	}
}