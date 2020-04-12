package comet

import (
	"chalurania/api"
	"chalurania/comet/connection"
	"chalurania/comet/router"
	"chalurania/service/config"
	"chalurania/service/log"
	"fmt"
	"net"
	"time"
)

// server 服务类
type Server struct {
	// 服务器 Id
	Id string

	// 服务器名称
	Name string

	// 协议版本
	netWork string

	// 服务器绑定的 IP
	IP string

	// 服务器绑定的端口号
	Port int

	// 消息管理模块，用来绑定 request id 和对应的处理方法
	RouterManager api.IRouterManager

	// 连接管理器
	ConnManager api.IConnectionManager

	// Server 连接创建时 Hook 函数
	OnConnStart func (conn api.IConnection)

	// Server 连接断开时的 Hook 函数
	OnConnStop func (conn api.IConnection)
}

// 初始化服务器
func NewServer() api.IServer {
	s := &Server{
		Name:          config.GlobalObj.Name,
		netWork:       "tcp",
		IP:            config.GlobalObj.Host,
		Port:          config.GlobalObj.Port,
		RouterManager: router.NewRouterManager(),
		ConnManager:   connection.NewConnectionManager(),
	}
	return s
}

// 启动服务器
func (s *Server) Start() {
	log.Info.Printf("server starting at IP: %s, Port: %d", s.IP, s.Port)

	// 服务端监听协程
	go func() {

		// 启动 worker 工作池
		s.RouterManager.StartWorkerPool()

		// 获取服务器监听地址
		addrStr := fmt.Sprintf("%s:%d", s.IP, s.Port)
		addr, err := net.ResolveTCPAddr(s.netWork, addrStr)
		if err != nil {
			log.Error.Println("resolveTCPAddr err:", err)
			return
		}

		// 监听服务器地址
		listener, err := net.ListenTCP(s.netWork, addr)
		if err != nil {
			log.Error.Println("listenTCP err:", err)
			return
		}

		// 服务器正在监听
		log.Info.Println("server start success, listening...")

		// TODO 自动生成 Id
		var cid uint32 = 0

		// 与客户端建立连接
		for {
			// 阻塞等待建立连接请求
			conn, err := listener.AcceptTCP()
			if err != nil {
				log.Error.Println("acceptTCP err:", err)
				// 连接失败继续等待下一次连接
				continue
			}

			// 设置服务器最大连接，如果超过最大连接，则丢弃当前连接
			if s.ConnManager.GetConnectionSize() >= config.GlobalObj.MaxConn {
				conn.Close()
				continue
			}

			// 处理新连接请求
			currentConn := connection.NewConnection(s, conn, cid, s.RouterManager)
			cid ++

			// 启动当前连接的处理业务
			go currentConn.Start()
		}
	}()
}

// 停止服务器
func (s *Server) Stop() {
	log.Info.Println("server stop success!")
	// 清理连接
	s.ConnManager.ClearConnection()
}

// 服务器服务方法
func (s *Server) Serve() {
	// 启动服务器
	s.Start()

	// 阻塞，否则主 go 退出，服务器监听 go 也会退出
	for {
		time.Sleep(time.Second)
	}
}

// 给当前服务注册路由方法，供客户端连接处理使用
func (s *Server)AddRouter(requestId uint32, router api.IRouter) {
	s.RouterManager.AddRouter(requestId, router)
}

// 得到连接管理器
func (s *Server) GetConnManager() api.IConnectionManager {
	return s.ConnManager
}

// 设置 server 的连接创建时调用的 Hook 函数
func (s *Server) SetOnConnStart(hookFunc func (api.IConnection)) {
	s.OnConnStart = hookFunc
}

// 设置 server 的连接断开时调用的 Hook 函数
func (s *Server) SetOnConnStop(hookFunc func (api.IConnection)) {
	s.OnConnStop = hookFunc
}

// 调用连接 OnConnStart Hook 函数
func (s *Server) CallOnConnStart(conn api.IConnection) {
	if s.OnConnStart != nil {
		//log.Info.Println("call on connection start...")
		s.OnConnStart(conn)
	}
}

// 调用连接 OnConnStop Hook 函数
func (s *Server) CallOnConnStop(conn api.IConnection) {
	if s.OnConnStop != nil {
		log.Info.Println("call on connection stop...")
		s.OnConnStop(conn)
	}
}