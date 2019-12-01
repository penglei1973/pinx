package pnet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"pinx/pinterface"
	"pinx/utils"
)

type Connection struct {
	TcpServer    pinterface.IServer    // 当前conn属于哪个server, 在conn初始化的时候添加即可
	Conn         *net.TCPConn          // 当前连接的socket TCP套接字
	ConnID       uint32                // 当前连接的ID SessionID
	IsClosed     bool                  // 当前连接的关闭状态
	ExitBuffchan chan bool             // 告知该连接已经退出/停止
	MsgHandle    pinterface.IMsgHandle // 连接处理的router
	MsgChan      chan []byte           // 无缓冲管道，用于读写两个goroutine之间消息通信
	MsgBuffChan  chan []byte           //有缓冲管道
}

// 创建连接的方法
func NewConntion(server pinterface.IServer, conn *net.TCPConn, connID uint32, msghandler pinterface.IMsgHandle) *Connection {
	c := &Connection{
		TcpServer:    server,
		Conn:         conn,
		ConnID:       connID,
		IsClosed:     false,
		ExitBuffchan: make(chan bool, 1),
		MsgHandle:    msghandler,
		MsgChan:      make(chan []byte),
		MsgBuffChan:  make(chan []byte, utils.GlobalObject.MaxMsgChanLen),
	}

	// 将新创建的Conn添加到连接管理中
	c.TcpServer.GetConnMgr().Add(c)

	return c
}

// Write Goroutine
func (c *Connection) StartWriter() {
	fmt.Println("Writer Goroutine is running")
	defer func() {
		fmt.Println(c.RemoteAddr().String(), "[conn write exit]")
		// c.Stop()
	}()

	for {
		select {
		case data, ok := <-c.MsgBuffChan:
			if ok {
				// 数据要写给客户端
				if _, err := c.Conn.Write(data); err != nil {
					fmt.Println("send data error : ", err)
					return
				}
			}
		case data, ok := <-c.MsgChan:
			if ok {
				// 数据要写给客户端
				if _, err := c.Conn.Write(data); err != nil {
					fmt.Println("send data error : ", err)
					return
				}
			}
		case <-c.ExitBuffchan:
			return
		}
	}
}

// Read的Goroutine
func (c *Connection) StartReader() {

	fmt.Println("Reader Gorountine is running")

	defer func() {
		fmt.Println(c.RemoteAddr().String(), " [reader conn exit]")
		c.Stop()
	}()

	for {
		fmt.Println("read begin...")
		dp := NewDataPack()
		headData := make([]byte, dp.GetHeadLen())
		_, err := io.ReadFull(c.GetTCPConnection(), headData) // readfull 把headdata填充满
		if err != nil {
			fmt.Println("read head error")
			c.ExitBuffchan <- true
			continue
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
			if msg, ok := msgHead.(*Message); ok {
				msg.Data = make([]byte, msg.GetDataLen())
				if _, err := io.ReadFull(c.GetTCPConnection(), msg.Data); err != nil {
					fmt.Println("server unpack data err : ", err)
					c.ExitBuffchan <- true
					continue
				}

				// fmt.Printf("Recv Msg : ID = %d, len = %d, Data = %s", msg.Id, msg.DataLen, string(msg.Data))
				// 得到当前用户的Request
				req := Request{
					conn: c,
					msg:  msg,
				}

				// 从路由Routers 中找到注册绑定Conn的对应handle
				if utils.GlobalObject.WorkerPoolSize > 0 {
					c.MsgHandle.SendMsgToTaskQueue(&req)
				} else {
					go c.MsgHandle.DoMsgHandler(&req)
				}
			}
		} else {
			fmt.Println("len <= 0")
			continue
		}

	}
}

// 启动连接， 让当前来捏开始工作
func (c *Connection) Start() {
	//开启处理连接读取到和护短数据之后的请求业务
	go c.StartReader()
	go c.StartWriter()

	for {
		select {
		case <-c.ExitBuffchan:
			// 得到退出消息不再阻塞

			return
		}
	}
}

// 停止连接， 结束当前来凝结状态
func (c *Connection) Stop() {
	// 1. 如果当前连接已经关闭
	if c.IsClosed {
		return
	}
	c.IsClosed = true

	// TODO Connection Stop() 如果用户注册了该连接的关闭回调业务，那么在此应该显示调用

	// 2. 关闭socket
	c.Conn.Close()

	// 通知从缓冲队列读数据的业务，该连接已经关闭
	c.ExitBuffchan <- true

	// 将连接从连接管理器中删除
	c.TcpServer.GetConnMgr().Remove(c)

	// 关闭该连接全部管道
	close(c.ExitBuffchan)
	close(c.MsgBuffChan)
	close(c.MsgChan)

}

// 从当前连接获取原始的socket TCPCONN
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// 获取当前连接
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// 获取 远程客户端地址信息
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// 发送数据给客户端
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.IsClosed == true {
		return errors.New("Connection closed when send msg")
	}

	// 封包发送
	dp := NewDataPack()
	msg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("pack error msg id = ", msgId)
		return errors.New("pack error msg")
	}

	c.MsgChan <- msg
	/*
		if _, err := c.Conn.Write(msg); err != nil {
			fmt.Println("write msg id", msgId, " error")
			c.ExitBuffchan <- true
			return errors.New("conn Write error")
		}
	*/

	return nil
}

func (c *Connection) SendBuffMsg(msgId uint32, data []byte) error {
	if c.IsClosed == true {
		return errors.New("Connection closed when send msg")
	}

	// 封包发送
	dp := NewDataPack()
	msg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("pack error msg id = ", msgId)
		return errors.New("pack error msg")
	}

	c.MsgBuffChan <- msg
	/*
		if _, err := c.Conn.Write(msg); err != nil {
			fmt.Println("write msg id", msgId, " error")
			c.ExitBuffchan <- true
			return errors.New("conn Write error")
		}
	*/

	return nil
}
