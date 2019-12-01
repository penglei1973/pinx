package utils

import (
	"encoding/json"
	"io/ioutil"
	"pinx/pinterface"
)

type GlobalObj struct {
	// Server
	TcpServer pinterface.IServer // 全局server对象
	Host      string             //主机IP
	TcpPort   int                // 服务器端口号
	Name      string             //服务器名称

	// Pinx
	Version          string // pinx版本号
	MaxPacketSize    uint32 //数据包的最大值
	MaxConn          int    // 服务器允许最大连接数
	WorkerPoolSize   uint32 //业务工作Worker池的数量
	MaxWorkerTaskLen uint32 //业务工作Worker对应负责的任务队列的最大任务存储数量
	MaxMsgChanLen    uint32 //SendBuffMsg发送消息的缓冲最大长度

	// config file path
	ConFilePath string
}

var GlobalObject *GlobalObj

func init() {
	// 初始化GlobalObject对象， 设置一些默认值
	GlobalObject = &GlobalObj{
		Name:             "ServerApp",
		Version:          "v6.6",
		TcpPort:          7777,
		Host:             "0.0.0.0",
		MaxConn:          12000,
		MaxPacketSize:    4096,
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
		MaxMsgChanLen:    1024,
		ConFilePath:      "conf/pinx.json",
	}

	GlobalObject.Reload()
}

// 获取用户配置文件
func (this *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/pinx.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}
