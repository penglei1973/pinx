package pnet

import "pinx/pinterface"

type Request struct {
	conn pinterface.IConnection // 和客户端建立好的连接
	msg  pinterface.IMessage    // 客户端请求的数据
}

// 获取请求连接信息
func (r *Request) GetConnection() pinterface.IConnection {
	return r.conn
}

// 获取请求消息的数据
func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

// 获取请求的消息id
func (r *Request) GetMsgId() uint32 {
	return r.msg.GetMsgId()
}
