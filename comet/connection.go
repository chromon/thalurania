package comet

import (
	"chalurania/api"
	"chalurania/service/config"
	"chalurania/service/log"
	"errors"
	"io"
	"net"
)

type Connection struct {
	// 当前连接 socket TCP 套接字
	Conn *net.TCPConn

	// 连接 Id
	ConnId uint32

	// 连接是否关闭
	isClosed bool

	// 通知连接是否退出/停止的 channel
	ExitChan chan bool

	// 消息管理，id 与对应处理方法
	RequestManager api.IRequestManager

	// 消息通信管道，用于读写两个协程之间通信
	MessageChan chan []byte
}

// 创建连接
func NewConnection(conn *net.TCPConn, connId uint32,
	requestManager api.IRequestManager) *Connection {
	c := &Connection{
		Conn:           conn,
		ConnId:         connId,
		isClosed:       false,
		ExitChan:       make(chan bool, 1),
		RequestManager: requestManager,
		MessageChan:    make(chan []byte),
	}
	return c
}

// 处理 conn 读取数据的协程
func (c *Connection) StartReader() {
	log.Info.Println("Reader goroutine running...")

	defer log.Info.Println(c.GetRemoteAddr().String(), " conn reader exit")
	defer c.Stop()

	// 循环读取数据
	for {
		// 创建数据包
		dp := NewDataPack()

		// 读取客户端 Message header
		header := make([]byte, dp.GetHeaderLen())
		// io.ReadFull 读取正好 len(headerLen) 长度的字节
		_, err := io.ReadFull(c.GetTCPConnection(), header)
		if err != nil {
			log.Error.Println("IO read message header err:", err)
			c.ExitChan <- true
			continue
		}

		// 拆包，得到 message id 和 data length
		msg, err := dp.Unpack(header)
		if err != nil {
			log.Error.Println("Unpack header err:", err)
			c.ExitChan <- true
			continue
		}

		// 根据 dataLen 读取 data
		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			_, err := io.ReadFull(c.GetTCPConnection(), data)
			if err != nil {
				log.Error.Println("IO read data err:", err)
				c.ExitChan <- true
				continue
			}
			msg.SetData(data)
		}

		// 得到当前客户端请求数据
		req := Request{
			conn: c,
			msg: msg,
		}

		if config.GlobalObj.WorkerPoolSize > 0 {
			// 已经启动工作池机制，将消息发送给 worker 处理
			c.RequestManager.SendRequestToTaskQueue(&req)
		} else {
			// 从 router 中找到注册绑定 conn 的对应 handle
			go c.RequestManager.ManageRequest(&req)
		}
	}
}

// 处理 conn 写入数据协程
func (c *Connection) StartWriter() {
	log.Info.Println("Writer goroutine running...")
	defer log.Info.Println(c.GetRemoteAddr().String(), " conn writer exit")

	for {
		select {
		case data := <- c.MessageChan:
			// 有数据写给客户端
			if _, err := c.Conn.Write(data); err!= nil {
				log.Error.Println("Writer write data err:", err)
				return
			}
		case <- c.ExitChan:
			// conn 已关闭
			return
		}
	}
}

// 启动连接
func (c *Connection) Start() {
	// 连接读取客户端数据并处理数据
	go c.StartReader()

	// 向客户端写入数据
	go c.StartWriter()

	for {
		select {
		case <- c.ExitChan:
			// 得到退出消息，不再阻塞
			return
		}
	}
}

// 停止连接
func (c *Connection) Stop() {
	if c.isClosed {
		return
	}

	c.isClosed = true

	// 关闭连接
	err := c.Conn.Close()
	if err != nil {
		log.Info.Println("Conn close err:", err)
	}

	// 通知 channel 连接关闭
	c.ExitChan <- true
	close(c.ExitChan)
}

// 获取当前连接
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// 获取当前连接 Id
func (c *Connection) GetConnId() uint32 {
	return c.ConnId
}

// 获取远程客户端地址信息
func (c *Connection) GetRemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// 将 Message 数据发送到远程 TCP 客户端
func (c *Connection) SendMsg(id uint32, data []byte) error {
	if c.isClosed {
		return errors.New("connection closed when send message")
	}

	// 将 data 封包，并发送
	dp := NewDataPack()
	msg, err := dp.Pack(NewMessage(id, data))
	if err != nil {
		log.Error.Println("Pack messsage id:", id, " err:", err)
		return err
	}

	// 发送到客户端
	c.MessageChan <- msg

	return nil
}