package packet

import (
	"bytes"
	"chalurania/api"
	"chalurania/service/config"
	"chalurania/service/log"
	"encoding/binary"
	"errors"
)

// 数据包
type DataPack struct {

	// 协议版本号 -- 1: tcp，2: websocket
	// NetWork uint32

	// 协议指令 (对应请求 requestId)
	// 1: 上行消息，2: 下行消息，3: auth认证，4: auth认证返回，
	// 5: 客户端请求心跳，6: 服务端心跳答复
	// 0: 回执信息，1: 注册，2：登录
	// Operation uint32

	// 消息id，服务端返回和客户端发送一一对应
	// MessageId uint32

	// 数据长度
	// DataLen uint32

	// 消息体
	// MessageData []byte

	// |NetWork|Operation|MessageId|DataLen|MessageData|
}

// 初始化数据包
func NewDataPack() *DataPack {
	return &DataPack{
	}
}

// 获取包头长度
func (dp *DataPack) GetHeaderLen() uint32 {
	// NetWork + Operation + MessageId + DataLen (4 * 2 + 8 + 4个字节)
	return 20
}

// 封包
func (dp *DataPack) Pack(netWork uint32, operation uint32, msg api.IMessage) ([]byte, error) {
	// 创建存放二进制 bytes 字节的缓冲
	dataBuf := bytes.NewBuffer([]byte{})

	// 序列化，将数据转换成 byte 字节流
	// 将协议版本写入 dataBuf 中
	err := binary.Write(dataBuf, binary.LittleEndian, netWork)
	if err != nil {
		log.Error.Println("binary write network err:", err)
		return nil, err
	}

	// 将协议指令 Id 写入 dataBuf 中
	err = binary.Write(dataBuf, binary.LittleEndian, operation)
	if err != nil {
		log.Error.Println("binary write operation err:", err)
	}

	// 将 Message id 写入 dataBuf 中
	err = binary.Write(dataBuf, binary.LittleEndian, msg.GetMsgId())
	if err != nil {
		log.Error.Println("binary write message id err:", err)
		return nil, err
	}

	// 将 data length 写入 dataBuf 中
	err = binary.Write(dataBuf, binary.LittleEndian, msg.GetDataLen())
	if err != nil {
		log.Error.Println("binary write data length err:", err)
		return nil, err
	}

	// 将 data 写入 dataBuf 中
	err = binary.Write(dataBuf, binary.LittleEndian, msg.GetData())
	if err != nil {
		log.Error.Println("binary write data err:", err)
		return nil, err
	}

	return dataBuf.Bytes(), nil
}

// 拆包
func (dp *DataPack) Unpack(binaryData []byte) (uint32, uint32, api.IMessage, error) {
	// 创建一个输出二进制数据的 IOReader
	dataBuf := bytes.NewReader(binaryData)

	var netWork uint32 = 0
	var operation uint32 = 0

	// 解析 header 信息，得到 Message Id 和 dataLen
	msg := &Message{}

	// 读取 Network
	err := binary.Read(dataBuf, binary.LittleEndian, &netWork)
	if err != nil {
		log.Error.Println("binary read network err:", err)
		return 0, 0, nil, err
	}

	// 读取 Operation
	err = binary.Read(dataBuf, binary.LittleEndian, &operation)
	if err != nil {
		log.Error.Println("binary read operation err:", err)
		return 0, 0, nil, err
	}

	// 读取 Message Id
	err = binary.Read(dataBuf, binary.LittleEndian, &msg.Id)
	if err != nil {
		log.Error.Println("binary read message id err:", err)
		return 0, 0, nil, err
	}

	// 读取 data length
	err = binary.Read(dataBuf, binary.LittleEndian, &msg.DataLen)
	if err != nil {
		log.Error.Println("binary read data length err:", err)
		return 0, 0, nil, err
	}

	// 判断 data length 长度是否超出允许的最大数据包长度
	if config.GlobalObj.MaxPacketSize > 0 &&
		msg.DataLen > config.GlobalObj.MaxPacketSize {
		log.Error.Println("too large message data received")
		return 0, 0, nil, errors.New("too large message data received")
	}

	// 拆包分两次，第一次解析 header，之后调用者再根据 dataLen 继续从 io 流中读取 body 中的数据
	return uint32(netWork), uint32(operation), msg, nil
}
