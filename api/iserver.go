package api

type IServer interface {
	// 启动服务器
	Start()

	// 停止服务器
	Stop()

	// 开启服务器业务
	Serve()

	// 添加路由
	AddRouter(msgId uint32, router IRouter)
}