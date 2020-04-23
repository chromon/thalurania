package packet

// 回执信息
type AckPack struct {
	// ack 指令
	Opt int32 `json:"opt"`

	Sign bool `json:"sign"`

	Data []byte `json:"data"`
}

func NewAckPack(opt int32, sign bool, data []byte) *AckPack {
	return &AckPack{
		Opt: opt,
		Sign: sign,
		Data: data,
	}
}