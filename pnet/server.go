package pnet

import (
	"fmt"
	"net"
	"pinx/pinterface"
	"pinx/utils"
	"time"
)

//ISERVER 接口的实现, 定义一个Server服务类
type Server struct {
	Name      string                  // 服务器名称
	IPVersion string                  // tcp4 or other
	IP        string                  // 服务绑定ip地址
	Port      int                     // 服务绑定端口
	Cid       uint32                  // 服务器连接请求数handler gorountine 开启的数量
	MsgRouter pinterface.IMsgHandle   //当前server绑定回调router
	ConnMar   pinterface.IConnManager // 连接管理模块

	OnConnStart func(conn pinterface.IConnection)
	OnConnStop  func(conn pinterface.IConnection)
}

func NewServer(name string) pinterface.IServer {
	s := &Server{
		Name:      utils.GlobalObject.Name,
		IPVersion: "tcp4",
		IP:        utils.GlobalObject.Host,
		Port:      utils.GlobalObject.TcpPort,
		Cid:       0,
		MsgRouter: NewMsgHandle(),
		ConnMar:   NewConnManager(),
	}

	return s
}

//实现 pinterface.IServer 的方法

func (s *Server) Start() {
	fmt.Printf("[START] Server listenner at IP: %s, Port %d, is starting\n", s.IP, s.Port)
	fmt.Printf("Version : %s, MaxConn: %d, MaxPacketSize:%d\n", utils.GlobalObject.Version, utils.GlobalObject.MaxConn, utils.GlobalObject.MaxPacketSize)

	// 开启Linstergoroutine
	go func() {
		// 0. 启动worker工作池机制
		s.MsgRouter.StartWorkerPool()

		// 1. 获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr err: ", err)
			return
		}

		// 2. 监听服务器地址
		listenner, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen", s.IPVersion, " err: ", err)
			return
		}

		// 已经监听成功
		fmt.Println("start pinx server ", s.Name, "successfu, listening...")

		// 3. 启动server网络连接业务
		for {
			// 3.1 阻塞顶戴客户端简历连接请求
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err ", err)
				continue
			} else {
				fmt.Printf("%s is comming\n", conn.RemoteAddr())
			}

			// 3.2 TODO Server.Start() 设置服务器最大连接限制, 如果超过最大连接，那么则关闭此新的连接
			if s.ConnMar.Len() > utils.GlobalObject.MaxConn {
				conn.Close()
				continue
			}

			// 3.3 TODO Server.Start() 处理该新连接请求的业务方法， 此时应该有 handler 和conn是绑定的æ
			dealConn := NewConntion(s, conn, s.Cid, s.MsgRouter)
			s.Cid++

			// 测试echo 服务器
			go dealConn.Start()
			/*
				go func() {
					//不断的循环从客户端获取数据
					for {
						buf := make([]byte, 512)
						cnt, err := conn.Read(buf)
						if err != nil {
							fmt.Println("recv buf err ", err)
							continue
						}
						//echo
						if _, err := conn.Write(buf[:cnt]); err != nil {
							fmt.Println("write back buf err ", err)
							continue
						}
					}
				}()
			*/
		}

	}()
}
func (s *Server) Stop() {
	fmt.Println("[STOP] Pinx server , name :", s.Name)

	// TODO Server.Stop() 将其他需要清理的连接信息或者其他信息也要一并停止或者清理
	s.ConnMar.ClearConn()
}
func (s *Server) Serve() {
	s.Start()

	// TODO Server.Serve() 是否在启动服务的时候 还要处理其他的事情

	// 阻塞
	for {
		time.Sleep(time.Second * 100)
	}
}

func (s *Server) AddRouter(msgId uint32, router pinterface.IRouter) {
	s.MsgRouter.AddRouter(msgId, router)
	fmt.Println("Add Router successful")
}
func (s *Server) GetConnMgr() pinterface.IConnManager {
	return s.ConnMar
}

func (s *Server) SetOnConnStart(hookFunc func(pinterface.IConnection)) {
	s.OnConnStart = hookFunc
}

func (s *Server) SetOnConnStop(hookFunc func(pinterface.IConnection)) {
	s.OnConnStop = hookFunc
}

func (s *Server) CallOnConnStart(conn pinterface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Printf("callonconnstart...")
		s.OnConnStart(conn)
	}
}
func (s *Server) CallOnConnStop(conn pinterface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Printf("callonconnstop...")
		s.OnConnStop(conn)
	}
}
