package constants

// 协议指令（路由请求 id）
const (
	// 回执
	AckOption = 100

	// 注册
	SignUpOption = 101

	// 登录
	LoginOption = 102

	// 登出
	LogoutOption = 103
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

	// 设备下线
	DeviceOffline = 303

	// 登出
	LogoutAckOpt = 304
)

// 网络协议
const (
	// tcp
	TCPNetwork = 401
)

// 服务器数据传输指令
const (
	// 强制下线
	KickOut = 501
)

// 客户端命令
const (
	// 无效命令
	ErrorCommand = 600

	// 注册
	RegisterCommand = 601

	// 登录
	LoginCommand = 602

	// 登出
	LogoutCommand = 603
)