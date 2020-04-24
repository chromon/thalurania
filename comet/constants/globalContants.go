package constants

// 协议指令（路由请求 id）
const (
	// 回执
	AckOption = 100

	// 注册
	SignUpOption = 101

	// 登录
	LoginOption = 102
)

// 数据持久化协议指令
const (
	// 注册信息
	SignUpPersistenceOpt = 201
)

// 回执信息指令
const (
	// 注册
	SignUpAckOpt = 301

	// 登录
	LoginAckOpt = 302
)

// 网络协议
const (
	// tcp
	TCPNetwork = 401
)