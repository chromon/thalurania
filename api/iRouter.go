package api

// 路由接口
type IRouter interface {
	// 处理连接业务之前的回调方法
	PreHandle(r IRequest)

	// 处理连接业务的方法
	Handle(r IRequest)

	// 处理连接业务之后的回调方法
	PostHandle(r IRequest)
}