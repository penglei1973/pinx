package main

import (
	"fmt"
	"io"
	"net"
	"pinx/pnet"
	"time"
)

func main() {
	for i := 0; i < 10000; i++ {
		go func() {
			fmt.Println("client test start ...")

			conn, err := net.Dial("tcp", "127.0.0.1:9999")
			if err != nil {
				fmt.Println(err)
				return
			}

			for {
				dp := pnet.NewDataPack()

				msg, _ := dp.Pack(pnet.NewMsgPackage(1, []byte("test message...")))
				_, err := conn.Write(msg)
				if err != nil {
					fmt.Println(err)
					return
				}

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
							break
						}
					}
				}
				time.Sleep(time.Second)

			}
		}()
	}
	for {
		time.Sleep(time.Second * 100)
	}
}
