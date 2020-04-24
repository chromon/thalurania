package main

import (
	"bufio"
	"chalurania/service/log"
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	log.Info.Println("Client start")

	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Error.Println("Net dial err:", err)
		return
	}

	// scanner 用户读取客户端命令
	scanner := bufio.NewScanner(os.Stdin)

	go func() {
		for {
			// 读取命令
			fmt.Print("~ ")
			scanner.Scan()
			if err := scanner.Err(); err != nil {
				_, err = fmt.Fprintln(os.Stderr, "error:", err)
			}
			_, err := conn.Write([]byte(scanner.Text()))
			if err != nil {
				log.Error.Println("Conn write err:", err)
				return
			}
			//time.Sleep(time.Second)

		}
	}()

	go func() {

		for {
			buf := make([]byte, 1024)
			n, err := conn.Read(buf)
			if err != nil {
				log.Error.Println("Conn Read err:", err)
				return
			}
			//fmt.Printf("\r\033[k")
			fmt.Printf("\b\bServer feedback: %s, len: %d \n", buf[:n], n)
			//fmt.Println()
			fmt.Print("~ ")
		}

	}()

	for {
		time.Sleep(time.Second)
	}
}