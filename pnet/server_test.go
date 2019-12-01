package pnet

import (
	"fmt"
	"net"
	"testing"
	"time"
)

func ClintTest() {
	fmt.Println("Client Test ... start")

	// 等待三秒，让服务器起来
	time.Sleep(time.Second * 3)

	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}

	for {
		_, err := conn.Write([]byte("hello pinx"))
		if err != nil {
			fmt.Println("write error :", err)
			return
		}

		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read buf error")
			if err != nil {
				fmt.Println("read buf error")
				return
			}
		}

		fmt.Printf("Server call back : %s, cnt = %d\n", buf[:cnt], cnt)
		time.Sleep(time.Second * 1)
	}
}

// Server 模块测试
func TestServer(t *testing.T) {
	server := NewServer("[pinx test]")
	go ClintTest()
	server.Serve()
}
