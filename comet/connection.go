package comet

import (
	"chalurania/api"
	"chalurania/service/config"
	"chalurania/service/log"
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

	// 连接业务的处理方法 router
	Router api.IRouter
}

// 创建连接
func NewConnection(conn *net.TCPConn, connId uint32,
	router api.IRouter) *Connection {
	c := &Connection{
		Conn: conn,
		ConnId: connId,
		isClosed: false,
		ExitChan: make(chan bool, 1),
		Router: router,
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
		buf := make([]byte, config.GlobalObj.MaxPacketSize)
		_, err := c.Conn.Read(buf)
		if err != nil {
			log.Error.Println("Conn read buf err:", err)
			c.ExitChan <- true
			continue
		}

		// 得到当前客户端请求数据
		req := Request{
			conn: c,
			data: buf,
		}

		// 从 router 中找到注册绑定 conn 的对应 handle
		go func (r api.IRequest) {
			// 执行注册的路由方法
			c.Router.PreHandle(r)
			c.Router.Handle(r)
			c.Router.PostHandle(r)
		}(&req)
	}

}

// 启动连接
func (c *Connection) Start() {
	// 连接读取客户端数据并处理数据
	go c.StartReader()

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