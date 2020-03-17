package comet

import (
	"chalurania/api"
	"chalurania/service/log"
	"strconv"
)

type MsgHandler struct {
	// 存放每个 id 对应的处理方法
	Apis map[uint32] api.IRouter
}

// 创建消息管理
func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis: make(map[uint32] api.IRouter),
	}
}

// 以非阻塞式处理消息
func (mh *MsgHandler) DoMsgHandler(request api.IRequest) {
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		log.Error.Println("Message handler not found message id:", request.GetMsgID())
		return
	}

	// 执行相应的处理方法
	handler.PostHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// 为消息添加具体的处理逻辑
func (mh *MsgHandler) AddRouter(msgId uint32, router api.IRouter) {
	// 判断当前 Message 绑定的处理方法是否存在
	if _, ok := mh.Apis[msgId]; ok {
		log.Error.Println("Repeated router api, message id:", strconv.Itoa(int(msgId)))
	}

	// 添加 message 与 router api 对应关系
	mh.Apis[msgId] = router
	log.Info.Println("Add router api message id:", msgId)
}