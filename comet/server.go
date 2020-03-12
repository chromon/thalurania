package comet

import (
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
}

func NewServer(name string) *Server {
	s := &Server{
		Name: name,
		netWork: "tcp",
		IP: "127.0.0.1",
		Port: 8080,
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
		log.Info.Println("Start server success, listening...")

		// 与客户端建立连接
		for {
			// 阻塞等待建立连接请求
			conn, err := listener.AcceptTCP()
			if err != nil {
				log.Error.Println("AcceptTCP err:", err)
				// 连接失败继续等待下一次连接
				continue
			}

			// 从客户的读取数据协程
			go func() {
				// 循环读取
				for {
					buf := make([]byte, 1024)
					n, err := conn.Read(buf)
					if err != nil {
						log.Error.Println("Conn read err:", err)
						// 读取失败继续等待下一次读取
						continue
					}

					// 简单将数据回写
					if _, err := conn.Write(buf[:n]); err != nil {
						log.Error.Println("Conn write err:", err)
						continue
					}
				}
			}()
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