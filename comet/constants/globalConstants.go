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

	// 搜索用户名
	SearchOption = 104

	// 添加好友
	FriendRequestOption = 105

	// 好友请求列表
	FriendReqListOption = 106

	// 接受好友请求
	AcceptFriendOption = 107
)

// 数据持久化协议指令
const (
	// 注册信息
	SignUpPersistenceOpt = 201

	// 添加好友请求
	FriendRequestPersistOpt = 202
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

	// 搜索
	SearchAckOpt = 305

	// 添加好友请求
	FriendRequestAckOpt = 306

	// 好友请求列表
	FriendReqListAckOpt = 307
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

	// 发送好友请求
	SendFriendRequest = 502
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

	// 搜索用户名
	SearchUsernameCommand = 604

	// 搜索用户 id
	SearchUserIdCommand = 605

	// 通过用户名添加好友
	AddUserByNameCommand = 606

	// 通过用户 id 添加好友
	AddUserByIdCommand = 607

	// 好友请求列表
	FriendReqListCommand = 608

	// 通过用户名接受好友请求
	AcceptFriendByNameCommand = 609

	// 通过用户 id 接受好友请求
	AcceptFriendByIdCommand = 610
)