package pnet

import (
	"errors"
	"fmt"
	"pinx/pinterface"
	"sync"
)

type ConnManager struct {
	Connections map[uint32]pinterface.IConnection // 管理的连接信息
	ConnLock    sync.RWMutex                      // 读写连接的读写锁
}

// 创建一个连接管理模块
func NewConnManager() *ConnManager {
	return &ConnManager{
		Connections: make(map[uint32]pinterface.IConnection),
	}
}

// 添加连接
func (connMgr *ConnManager) Add(conn pinterface.IConnection) {
	// 保护共享资源map 加写锁
	connMgr.ConnLock.Lock()
	defer connMgr.ConnLock.Unlock()

	// 将conn连接添加到ConnMananger中
	connMgr.Connections[conn.GetConnID()] = conn

	fmt.Printf("connection add to ConnManager successfully: conn num = %d \n", connMgr.Len())
}

// 删除连接
func (connMgr *ConnManager) Remove(conn pinterface.IConnection) {
	// 保护共享资源Map 加写锁
	connMgr.ConnLock.Lock()
	defer connMgr.ConnLock.Unlock()

	// 删除连接信息
	delete(connMgr.Connections, conn.GetConnID())

	fmt.Printf("Connections Remove ConnID = %d successfully : conn num = %d\n", conn.GetConnID(), connMgr.Len())
}

// 利用ConnID获取连接
func (connMgr *ConnManager) Get(connID uint32) (pinterface.IConnection, error) {
	// 保护共享资源Map 加读锁
	connMgr.ConnLock.RLock()
	defer connMgr.ConnLock.RUnlock()

	if conn, ok := connMgr.Connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not found")
	}
}

//获取当前连接
func (connMgr *ConnManager) Len() int {
	return len(connMgr.Connections)
}

// 删除并停止所有连接
func (connMgr *ConnManager) ClearConn() {
	// 保护共享资源map加写锁
	connMgr.ConnLock.Lock()
	defer connMgr.ConnLock.Unlock()

	// 停止并删除全部的连接信息
	for connID, conn := range connMgr.Connections {
		conn.Stop()
		delete(connMgr.Connections, connID)
	}

	fmt.Printf("clear all connections successfully : conn num = %d", connMgr.Len())
}
