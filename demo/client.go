package main

import (
	"chalurania/service/log"
	"net"
	"time"
)

func main() {
	log.Info.Println("Client start")

	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Error.Println("Net dial err:", err)
		return
	}

	for {
		_, err := conn.Write([]byte("HelloWorld"))
		if err != nil {
			log.Error.Println("Conn write err:", err)
			return
		}

		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			log.Error.Println("Conn Read err:", err)
			return
		}

		log.Info.Printf("Server feedback: %s, len: %d", buf[:n], n)
		time.Sleep(time.Second)
	}
}