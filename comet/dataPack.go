package comet

import (
	"bytes"
	"chalurania/api"
	"chalurania/service/config"
	"chalurania/service/log"
	"encoding/binary"
	"errors"
)

// 数据包
type DataPack struct {}

// 初始化数据包
func NewDataPack() *DataPack {
	return &DataPack{}
}

// 获取包头长度
func (dp *DataPack) GetHeaderLen() uint32 {
	// Id uint32 (4个字节) + DataLen uint32 (4个字节)
	return 8
}

// 封包
func (dp *DataPack) Pack(msg api.IMessage) ([]byte, error) {
	// 创建存放二进制 bytes 字节的缓冲
	dataBuf := bytes.NewBuffer([]byte{})

	// 序列化，将数据转换成 byte 字节流，即将 Message id 写入 dataBuf 中
	err := binary.Write(dataBuf, binary.LittleEndian, msg.GetMsgId())
	if err != nil {
		log.Error.Println("Binary write message id err:", err)
		return nil, err
	}

	// 将 data length 写入 dataBuf 中
	err = binary.Write(dataBuf, binary.LittleEndian, msg.GetDataLen())
	if err != nil {
		log.Error.Println("Binary write data length err:", err)
		return nil, err
	}

	// 将 data 写入 dataBuf 中
	err = binary.Write(dataBuf, binary.LittleEndian, msg.GetData())
	if err != nil {
		log.Error.Println("Binary write data err:", err)
		return nil, err
	}

	return dataBuf.Bytes(), nil
}

// 拆包
func (dp *DataPack) Unpack(binaryData []byte) (api.IMessage, error) {
	// 创建一个输出二进制数据的 IOReader
	dataBuf := bytes.NewReader(binaryData)

	// 解析 header 信息，得到 Message Id 和 dataLen
	msg := &Message{}

	// 读取 Message Id
	err := binary.Read(dataBuf, binary.LittleEndian, &msg.Id)
	if err != nil {
		log.Error.Println("Binary read message id err:", err)
		return nil, err
	}

	// 读取 data length
	err = binary.Read(dataBuf, binary.LittleEndian, &msg.DataLen)
	if err != nil {
		log.Error.Println("Binary read data length err:", err)
		return nil, err
	}

	// 判断 data length 长度是否超出允许的最大数据包长度
	if config.GlobalObj.MaxPacketSize > 0 &&
		msg.DataLen > config.GlobalObj.MaxPacketSize {
		log.Error.Println("Too large message data received")
		return nil, errors.New("too large message data received")
	}

	// 拆包分两次，第一次解析 header，之后调用者再根据 dataLen 继续从 io 流中读取 body 中的数据
	return msg, nil
}
