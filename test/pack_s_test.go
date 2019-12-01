package test

import (
	"fmt"
	"io"
	"net"
	"pinx/pnet"
	"testing"
)

func TestPackS(t *testing.T) {
	// 创建socket Tcp Server
	listener, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("server listen err:", err)
		return
	} else {
		fmt.Println("127.0.0.1:7777 is Listening...")
	}

	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println("accept err: ", err)
		} else {
			fmt.Println("conn is comming")
		}

		go func(conn net.Conn) {
			dp := pnet.NewDataPack()

			for {
				headData := make([]byte, dp.GetHeadLen())
				_, err := io.ReadFull(conn, headData) // readfull 把headdata填充满
				if err != nil {
					fmt.Println("read head error")
					break
				} else {
					fmt.Println("head is ", headData)
				}

				// 将headData字节流 拆包到msg中
				msgHead, err := dp.Unpack(headData)
				if err != nil {
					fmt.Println("unpack err: ", err)
					return
				}

				if msgHead.GetDataLen() > 0 {
					if msg, ok := msgHead.(*pnet.Message); ok {
						msg.Data = make([]byte, msg.GetDataLen())
						if _, err := io.ReadFull(conn, msg.Data); err != nil {
							fmt.Println("server unpack data err : ", err)
							return
						}

						fmt.Printf("Recv Msg : ID = %d, len = %d, Data = %s", msg.Id, msg.DataLen, string(msg.Data))
					}
				}
			}
		}(conn)
	}
}
