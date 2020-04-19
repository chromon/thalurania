package packet

// 用于包装实体类，在消息队列中实现异步数据存储
// opt: 1, 添加用户信息
type DataWrap struct {
	// 数据操作指令
	Opt int32 `json:"opt"`

	// 接收到的 json 数据
	Model []byte `json:"model"`
}

func NewDataWrap(opt int32, model []byte) *DataWrap {
	return &DataWrap{
		Opt:   opt,
		Model: model,
	}
}

func (d *DataWrap) GetOpt() int32 {
	return d.Opt
}

func (d *DataWrap) GetModel() []byte {
	return d.Model
}
