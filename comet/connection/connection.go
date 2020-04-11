package connection

import (
	"chalurania/api"
	"chalurania/comet/packet"
	"chalurania/comet/router"
	"chalurania/service/config"
	"chalurania/service/log"
	"errors"
	"io"
	"net"
	"sync"
)

type Connection struct {
	// 当前 connection 隶属于的 server
	TCPServer api.IServer

	// 当前连接 socket TCP 套接字
	Conn *net.TCPConn

	// 连接 Id
	ConnId uint32

	// 连接是否关闭
	isClosed bool

	// 通知连接是否退出/停止的 channel
	ExitChan chan bool

	// 消息管理，id 与对应处理方法
	RouterManager api.IRouterManager

	// 无缓冲消息通信管道，用于读写两个协程之间通信
	MessageChan chan []byte

	// 有缓冲消息通信管道，用户读写两个协程之间通信
	MessageBufChan chan []byte

	// 连接属性
	property map[string]interface{}

	// 保护连接属性修改的锁
	propertyLock sync.RWMutex
}

// 创建连接
func NewConnection(server api.IServer, conn *net.TCPConn, connId uint32,
	requestManager api.IRouterManager) *Connection {
	c := &Connection{
		TCPServer:      server,
		Conn:           conn,
		ConnId:         connId,
		isClosed:       false,
		ExitChan:       make(chan bool, 1),
		RouterManager:  requestManager,
		MessageChan:    make(chan []byte),
		MessageBufChan: make(chan []byte, config.GlobalObj.MaxMsgChanLen),
		property:       make(map[string]interface{}),
	}

	// 将新建的 conn 连接添加到连接管理器中
	c.TCPServer.GetConnManager().AddConnection(c)

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
		dp := packet.NewDataPack()

		// 读取客户端 Message header
		header := make([]byte, dp.GetHeaderLen())
		// io.ReadFull 读取正好 len(headerLen) 长度的字节
		_, err := io.ReadFull(c.GetTCPConnection(), header)
		if err != nil {
			// log.Error.Println("IO read message header err:", err)
			c.ExitChan <- true
			break
		}

		// 拆包，得到 network, operation, message id 和 data length
		_, operation, msg, err := dp.Unpack(header)
		if err != nil {
			log.Error.Println("Unpack header err:", err)
			c.ExitChan <- true
			break
		}

		// 根据 dataLen 读取 data
		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			_, err := io.ReadFull(c.GetTCPConnection(), data)
			if err != nil {
				log.Error.Println("IO read data err:", err)
				c.ExitChan <- true
				break
			}
			msg.SetData(data)
		}

		// 得到当前客户端请求数据的协议指令 (对应请求 requestId)
		var rId = operation

		req := router.Request{
			RequestId: rId,
			Conn:      c,
			Message:   msg,
		}

		if config.GlobalObj.WorkerPoolSize > 0 {
			// 已经启动工作池机制，将消息发送给 worker 处理
			c.RouterManager.SendRequestToTaskQueue(&req)
		} else {
			// 从 router 中找到注册绑定 conn 的对应 handle
			go c.RouterManager.ManageRequest(&req)
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
		case data, ok := <- c.MessageBufChan:
			if ok {
				// 有数据要写给客户端
				if _, err := c.Conn.Write(data); err != nil {
					log.Error.Println("Writer write buf data err:", err)
					return
				}
			} else {
				log.Info.Println("Message buffer chan closed")
				break
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

	// 执行用户传进来创建连接时需要处理的业务 hook 函数
	c.TCPServer.CallOnConnStart(c)

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

	// 执行用户注册的关闭连接时回调函数
	c.TCPServer.CallOnConnStop(c)

	// 关闭连接
	err := c.Conn.Close()
	if err != nil {
		log.Info.Println("Conn close err:", err)
	}

	// 通知 channel 连接关闭
	c.ExitChan <- true

	// 将连接从连接管理器中删除
	c.TCPServer.GetConnManager().RemoveConnection(c)

	close(c.ExitChan)
	close(c.MessageChan)
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
func (c *Connection) SendMsg(netWork uint32, operation uint32,
	id int64, data []byte) error {
	if c.isClosed {
		return errors.New("connection closed when send message")
	}

	// 将 data 封包，并发送
	dp := packet.NewDataPack()
	dataBuf, err := dp.Pack(netWork, operation, packet.NewMessage(id, data))
	if err != nil {
		log.Error.Println("pack message id:", id, " err:", err)
		return err
	}

	// 发送到客户端
	c.MessageChan <- dataBuf

	return nil
}

// 将 Message 数据发送到远程 TCP 客户端（有缓冲）
func (c *Connection) SendBufMsg(netWork uint32, operation uint32,
	msgId int64, data []byte) error {
	if c.isClosed {
		return errors.New("connection closed when send message")
	}

	// 将 data 封包，并发送
	dp := packet.NewDataPack()
	dataBuf, err := dp.Pack(netWork, operation, packet.NewMessage(msgId, data))
	if err != nil {
		log.Error.Println("pack message id:", msgId, " err:", err)
		return err
	}

	// 发送到客户端
	c.MessageBufChan <- dataBuf

	return nil
}


// 设置连接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.property[key] = value
}

// 获取连接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	}

	return nil, errors.New("no property found")
}

// 移除连接属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}