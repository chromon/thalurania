package api

// 消息管理
type IRequestManager interface {
	// 以非阻塞式处理消息
	ManageRequest(IRequest)

	// 为消息添加具体的处理逻辑
	AddRouter(uint32, IRouter)

	// 启动一个 Worker
	StartWorker(int, chan IRequest)

	// 启动 worker 工作池
	StartWorkerPool()

	// 将消息交给 TaskQueue，由 Worker 进行处理
	SendRequestToTaskQueue(IRequest)
}