package main

import (
	"chalurania/service/log"
	"fmt"
	"net"
	"time"
)

func main() {
	// 服务端监听协程
	go func() {
		// 获取服务器监听地址
		addrStr := fmt.Sprintf("%s:%d", "127.0.0.1", 8080)
		addr, err := net.ResolveTCPAddr("tcp", addrStr)
		if err != nil {
			log.Error.Println("ResolveTCPAddr err:", err)
			return
		}

		// 监听服务器地址
		listener, err := net.ListenTCP("tcp", addr)
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
					fmt.Println(n)
					// 简单将数据回写
					if _, err := conn.Write(buf[:n]); err != nil {
						log.Error.Println("Conn write err:", err)
						continue
					}
				}
			}()
		}
	}()

	for {
		time.Sleep(time.Second)
	}
}
