package api

// 连接管理
type IConnectionManager interface {
	// 添加连接
	AddConnection(IConnection)

	// 删除连接
	RemoveConnection(IConnection)

	// 由 connId 获取连接
	GetConnection(uint32) (IConnection, error)

	// 当前连接数量
	GetConnectionSize() int

	// 删除并品质所有连接
	ClearConnection()
}