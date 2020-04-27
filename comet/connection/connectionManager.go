package connection

import (
	"chalurania/api"
	"chalurania/service/log"
	"errors"
	"sync"
)

// 管理连接模块
type ConnectionManager struct {
	// 管理的连接信息
	connections map[int64]api.IConnection

	// 读写连接的读写锁
	connectionLock sync.RWMutex
}

// 创建连接管理
func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[int64] api.IConnection),
	}
}

// 添加连接
func (cm *ConnectionManager) AddConnection(conn api.IConnection) {
	// 保护共享资源 map，加写锁
	cm.connectionLock.Lock()
	defer cm.connectionLock.Unlock()

	// 将 conn 连接添加到连接管理中
	cm.connections[conn.GetConnId()] = conn
	log.Info.Println("connection add conn id:", conn.GetConnId(),
		"server connections count:", cm.GetConnectionSize())
}

// 删除连接仅从 map 中删除，并未停止连接业务
func (cm *ConnectionManager) RemoveConnection(conn api.IConnection) {
	// 保护共享资源，加写锁
	cm.connectionLock.Lock()
	defer cm.connectionLock.Unlock()

	// 删除连接信息
	delete(cm.connections, conn.GetConnId())

	log.Info.Println("connection removed conn id:", conn.GetConnId(),
		"server connections count:", cm.GetConnectionSize())
}

// 由 connId 获取连接
func (cm *ConnectionManager) GetConnection(connId int64) (conn api.IConnection, err error) {
	// 保护共享资源，加读锁
	cm.connectionLock.RLock()
	defer cm.connectionLock.RUnlock()

	if conn, ok := cm.connections[connId]; ok {
		return conn, nil
	}

	return nil, errors.New("connection not found")
}

// 当前连接数量
func (cm *ConnectionManager) GetConnectionSize() int {
	return len(cm.connections)
}

// 删除并停止所有连接
func (cm *ConnectionManager) ClearConnection() {
	// 保护共享资源，加写锁
	cm.connectionLock.Lock()
	defer cm.connectionLock.Unlock()

	// 停止并删除全部连接信息
	for connId, conn := range cm.connections {
		// 停止连接
		conn.Stop()
		// 删除连接
		delete(cm.connections, connId)
	}

	log.Info.Println("clear all connections success, connection count:", cm.GetConnectionSize())
}