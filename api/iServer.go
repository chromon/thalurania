package api

type IServer interface {
	// 启动服务器
	Start()

	// 停止服务器
	Stop()

	// 开启服务器业务
	Serve()

	// 添加路由
	AddRouter(uint32, IRouter)

	// 得到连接管理器
	GetConnManager() IConnectionManager

	// 设置 server 的连接创建时调用的 Hook 函数
	SetOnConnStart(func (IConnection))

	// 设置 server 的连接断开时调用的 Hook 函数
	SetOnConnStop(func (IConnection))

	// 调用连接 OnConnStart Hook 函数
	CallOnConnStart(IConnection)

	// 调用连接 OnConnStop Hook 函数
	CallOnConnStop(IConnection)
}