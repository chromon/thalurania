package packet

// 消息
type Message struct {
	// 消息 Id
	Id int64

	// 消息长度
	DataLen uint32

	// 消息内容
	Data []byte
}

// 创建消息
func NewMessage(id int64, data []byte) *Message {
	return &Message{
		Id: id,
		DataLen: uint32(len(data)),
		Data: data,
	}
}

// 获取消息 Id
func (m *Message) GetMsgId() int64 {
	return m.Id
}

// 获取消息数据段长度
func (m *Message) GetDataLen() uint32 {
	return m.DataLen
}

// 获取消息内容
func (m *Message) GetData() []byte {
	return m.Data
}

// 设置消息 Id
func (m *Message) SetMsgId(id int64) {
	m.Id = id
}

// 设置消息内容
func (m *Message) SetData(data []byte) {
	m.Data = data
}

// 设置消息数据段长度
func (m *Message) SetDataLen(dataLen uint32) {
	m.DataLen = dataLen
}