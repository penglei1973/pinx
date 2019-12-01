package main

import (
	"fmt"
	"net/http"
	"pinx/pinterface"
	"pinx/pnet"
	"runtime/pprof"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	p := pprof.Lookup("goroutine")
	p.WriteTo(w, 1)
}

type PingRouter struct {
	pnet.BaseRouter // 继承
}

// 处理conn业务之前的hook的方法
func (this *PingRouter) PreHandle(request pinterface.IRequest) {
	/*
		time.Sleep(time.Second * 2)
		fmt.Println("ping PreHandle...")
		_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping ... \n"))
		if err != nil {
			fmt.Println("call back ping error")
		}
	*/
}

// 处理conn业务方法
func (this *PingRouter) Handle(request pinterface.IRequest) {
	//fmt.Println("ping handle")
	//fmt.Println("recv from client : msgId = ", request.GetMsgId())
	//fmt.Printf("recv from client : data = %s\n", request.GetData())
	request.GetConnection().SendMsg(1, []byte("hello world"))
	// dp := pnet.NewDataPack()
	// msg, _ := dp.Pack(pnet.NewMsgPackage(0, []byte("call back!!!!")))
	// _, err := request.GetConnection().GetTCPConnection().Write(msg)
	/*
		if err != nil {
			fmt.Println("call back ping error")
		}
	*/
}

//处理完conn的hook
func (this *PingRouter) PostHandle(request pinterface.IRequest) {
	/*
		fmt.Println("ping POstHandle...")
		_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping ... \n"))
		if err != nil {
			fmt.Println("call back ping error")
		}
	*/
}

func conBegin(conn pinterface.IConnection) {
	fmt.Println("conn begin")
	err := conn.SendMsg(2, []byte("connn begin .........."))
	if err != nil {
		fmt.Println(err)
	}
}

func conEnd(conn pinterface.IConnection) {
	fmt.Println("conn end")
}

func main() {
	// 创建一个server实例
	s := pnet.NewServer("[pinx with router1.0]")
	s.SetOnConnStart(conBegin)
	s.SetOnConnStop(conEnd)
	s.AddRouter(1, &PingRouter{})
	http.HandleFunc("/", handler)
	go func() {
		fmt.Println("pprof start...")
		fmt.Println(http.ListenAndServe(":9876", nil))
	}()
	// 开启服务
	s.Serve()
}
