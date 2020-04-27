package logic

import (
	"bufio"
	"chalurania/comet/constants"
	"chalurania/comet/packet"
	"chalurania/service/log"
	"chalurania/service/model"
	"chalurania/service/sequence"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
)

// 登录
func Login(m map[string]*flag.Flag, conn net.Conn) {
	scanner := bufio.NewScanner(os.Stdin)
	IdWorker, _ := sequence.NewWorker(0)

	// 读取密码
	fmt.Print("~ ")
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		_, err = fmt.Fprintln(os.Stderr, "error:", err)
	}
	pwd := scanner.Text()

	// 创建用户对象
	u := model.User{Username: m["u"].Value.String(), Password: pwd}
	// 序列化用户对象
	ret, err := json.Marshal(u)
	if err != nil {
		log.Info.Println("serialize user object err:", err)
	}

	// 封包用户对象消息并发送
	dp := packet.NewDataPack()
	msg, _ := dp.Pack(constants.TCPNetwork, constants.LoginOption, packet.NewMessage(IdWorker.GetId(), ret))
	_, err = conn.Write(msg)
	if err != nil {
		log.Error.Println("client login write message err:", err)
	}
}