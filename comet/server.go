package comet

import (
	"chalurania/api"
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

	// 当前 server 由用户自定义绑定的回调 router，即当前连接的实际处理业务
	Router api.IRouter
}

// 初始化服务器
func NewServer() api.IServer {

	s := &Server{
		Name: config.GlobalObj.Name,
		netWork: "tcp",
		IP: config.GlobalObj.Host,
		Port: config.GlobalObj.Port,
		Router: nil,
	}
	return s
}

// 启动服务器
func (s *Server) Start() {
	log.Info.Printf("Server starting at IP: %s, Port: %d", s.IP, s.Port)

	// 服务端监听协程
	go func() {
		// 获取服务器监听地址
		addrStr := fmt.Sprintf("%s:%d", s.IP, s.Port)
		addr, err := net.ResolveTCPAddr(s.netWork, addrStr)
		if err != nil {
			log.Error.Println("ResolveTCPAddr err:", err)
			return
		}

		// 监听服务器地址
		listener, err := net.ListenTCP(s.netWork, addr)
		if err != nil {
			log.Error.Println("ListenTCP err:", err)
			return
		}

		// 服务器正在监听
		log.Info.Println("Server start success, listening...")

		// TODO 自动生成 Id
		var cid uint32 = 0

		// 与客户端建立连接
		for {
			// 阻塞等待建立连接请求
			conn, err := listener.AcceptTCP()
			if err != nil {
				log.Error.Println("AcceptTCP err:", err)
				// 连接失败继续等待下一次连接
				continue
			}

			// 处理新连接请求
			currentConn := NewConnection(conn, cid, s.Router)
			cid ++

			// 启动当前连接的处理业务
			go currentConn.Start()
		}
	}()
}

// 停止服务器
func (s *Server) Stop() {
	log.Info.Println("Server stop success!")
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
func (s *Server)AddRouter(router api.IRouter) {
	s.Router = router
	log.Info.Println("Add router success")
}