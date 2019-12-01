package pinterface

import "net"

// 定义连接接口
type IConnection interface {
	Start()                                  // 启动连接， 让当前来捏开始工作
	Stop()                                   // 停止连接， 结束当前来凝结状态
	GetTCPConnection() *net.TCPConn          // 从当前连接获取原始的socket TCPCONN
	GetConnID() uint32                       // 获取当前连接
	RemoteAddr() net.Addr                    // 获取 远程客户端地址信息
	SendMsg(msgId uint32, data []byte) error // 发送数据给客户端
	SendBuffMsg(uint32, []byte) error        // 发送给客户端带缓存
}

// 定义一个同一处理连接业务接口
type HandFunc func(*net.TCPConn, []byte, int) error

// arg[0] : socket原生连接
// arg[1] : 客户端请求数据
// arg[2] : 数据的长度
