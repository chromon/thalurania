package router

import (
	"chalurania/api"
)

// 路由功能，用于自定义连接处理业务
// 实现 router 时，根据需要对基类方法进行重写
type Router struct {}

// 根据需要自定义实现方法
func (r *Router) PreHandle(req api.IRequest) {}
func (r *Router) Handle(req api.IRequest) {}
func (r *Router) PostHandle(req api.IRequest) {}