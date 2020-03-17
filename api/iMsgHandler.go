package api

// 消息管理
type IMsgHandler interface {
	// 以非阻塞式处理消息
	DoMsgHandler(IRequest)

	// 为消息添加具体的处理逻辑
	AddRouter(uint32, IRouter)
}