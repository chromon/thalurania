package router

import (
	"chalurania/api"
	"chalurania/service/config"
	"chalurania/service/log"
	"reflect"
	"strconv"
)

type RouterManager struct {
	// 存放每个 id 对应的处理方法
	Routers map[uint32] api.IRouter

	// 业务工作池数量
	WorkerPoolSize int64

	// 任务队列，request 请求信息的 channel 集合，worker 会从对应的队列中获取客户端请求数据并处理
	TaskQueue []chan api.IRequest
}

// 创建消息管理
func NewRouterManager() *RouterManager {
	return &RouterManager{
		Routers: make(map[uint32] api.IRouter),
		WorkerPoolSize: config.GlobalObj.WorkerPoolSize,
		// TaskQueue 中的每个队列应该是和一个 Worker 对应，所以数量一致
		TaskQueue: make([]chan api.IRequest, config.GlobalObj.WorkerPoolSize),
	}
}

// 以非阻塞式处理消息
func (rm *RouterManager) ManageRequest(request api.IRequest) {
	manager, ok := rm.Routers[request.GetRequestId()]
	if !ok {
		log.Error.Println("manager request not found request id:", request.GetRequestId())
		return
	}

	// 执行相应的处理方法
	manager.PreHandle(request)
	manager.Handle(request)
	manager.PostHandle(request)
}

// 为消息请求添加具体的处理逻辑
func (rm *RouterManager) AddRouter(requestId uint32, router api.IRouter) {
	// 判断当前 request 绑定的处理方法是否存在
	if _, ok := rm.Routers[requestId]; ok {
		log.Error.Println("repeated router api, request id:", strconv.Itoa(int(requestId)))
	}

	// 添加 RequestId 与 router api 对应关系
	rm.Routers[requestId] = router

	log.Info.Printf("add router api [%s] mapping request id - %d", reflect.TypeOf(router).Elem().Name(), requestId)
}

// 启动一个 Worker
func (rm *RouterManager) StartWorker(workerId int, taskQueue chan api.IRequest) {
	log.Info.Println("worker id -", workerId, "is started")
	// 循环等待队列中的消息
	for {
		select {
		// 当有消息时则去除队列中的 Request，并执行绑定的业务方法
		case request := <- taskQueue:
			rm.ManageRequest(request)
		}
	}
}

// 启动 worker 工作池, 每一个 worker 分配一个 TaskQueue
func (rm *RouterManager) StartWorkerPool() {
	// 遍历需要启动的 worker 数量，并依次启动
	for i := 0; i < int(rm.WorkerPoolSize); i++ {
		// 给当前 worker 对应的任务队列开辟空间
		rm.TaskQueue[i] = make(chan api.IRequest, config.GlobalObj.MaxWorkerTaskLen)
		// 启动当前 worker，阻塞等待对应的任务队列是否有消息传递进来
		go rm.StartWorker(i, rm.TaskQueue[i])
	}
}

// 将消息交给 TaskQueue，由 Worker 进行处理
func (rm *RouterManager) SendRequestToTaskQueue(request api.IRequest) {
	// 根据 Conn Id 来分配当前的连接应该由哪个 Worker 负责处理
	// 轮询的平均分配法，得到需要处理此连接的 worker Id
	workerId := request.GetConnection().GetConnId() % rm.WorkerPoolSize
	//log.Info.Println("add Conn id:", request.GetConnection().GetConnId(),
	//	" request Message id:", request.GetMsgID(), " to worker id:", workerId)

	// 将请求消息发送给任务队列
	rm.TaskQueue[workerId] <- request
}