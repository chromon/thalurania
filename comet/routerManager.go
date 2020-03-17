package comet

import (
	"chalurania/api"
	"chalurania/service/log"
	"strconv"
)

type RouterManager struct {
	// 存放每个 id 对应的处理方法
	Routers map[uint32] api.IRouter
}

// 创建消息管理
func NewRouterManager() *RouterManager {
	return &RouterManager{
		Routers: make(map[uint32] api.IRouter),
	}
}

// 以非阻塞式处理消息
func (rm *RouterManager) ManageRequest(request api.IRequest) {
	manager, ok := rm.Routers[request.GetMsgID()]
	if !ok {
		log.Error.Println("Manager request not found message id:", request.GetMsgID())
		return
	}

	// 执行相应的处理方法
	manager.PostHandle(request)
	manager.Handle(request)
	manager.PostHandle(request)
}

// 为消息添加具体的处理逻辑
func (rm *RouterManager) AddRouter(msgId uint32, router api.IRouter) {
	// 判断当前 message 绑定的处理方法是否存在
	if _, ok := rm.Routers[msgId]; ok {
		log.Error.Println("Repeated router api, message id:", strconv.Itoa(int(msgId)))
	}

	// 添加 msgId 与 router api 对应关系
	rm.Routers[msgId] = router
	log.Info.Println("Add router api message id:", msgId)
}