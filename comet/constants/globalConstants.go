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

	// 好友列表
	FriendListOption = 108

	// 私聊
	PrivateChatOption = 109

	// 离线消息
	OfflineMessageOption = 110

	// 创建群组
	CreateGroupOption = 111

	// 群组邀请
	GroupInviteOption = 112

	// 群组邀请列表
	GroupInviteListOption = 113

	// 接受群组邀请
	AcceptGroupInviteOption = 114
)

// 数据持久化协议指令
const (
	// 注册信息
	SignUpPersistenceOpt = 201

	// 添加好友请求
	FriendRequestPersistOpt = 202

	// 添加消息
	MessagePersistenceOpt = 203

	// 群组邀请
	GroupRequestPersistOpt = 204
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

	// 接受好友请求
	AcceptFriendRepAckOpt = 308

	// 好友列表
	FriendListAckOpt = 309

	// 发送消息
	SendMessageAckOpt = 310

	// 离线消息
	OfflineMsgAckOpt = 311

	// 创建群组
	CreateGroupAckOpt = 312

	// 群组邀请
	GroupRequestAckOpt = 313

	// 群组邀请列表
	GroupInviteAckOpt = 314

	// 接收群组邀请
	AcceptGroupInviteAckOpt = 315
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

	// 发送消息
	SendMessage = 503

	// 登录成功
	LoginSuccess = 504

	// 发送群组请求
	SendGroupRequest = 505
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

	// 好友列表
	FriendListCommand = 611

	// 通过用户名与好友私聊
	ChatWithFriendByNameCommand = 612

	// 通过用户 id 与好友私聊
	ChatWithFriendByIdCommand = 613

	// 通过用户名查询离线消息
	OfflineMsgByNameCommand = 615

	// 通过用户 id 查询离线消息
	OfflineMsgByIdCommand = 616

	// 创建群组
	CreateGroupCommand = 617

	// 搜索用户组
	SearchGroupCommand = 618

	// 通过用户名邀请到群组
	GroupInviteByNameCommand = 619

	// 通过用户 id 邀请到群组
	GroupInviteByIdCommand = 620

	// 群组邀请列表
	GroupInviteListCommand = 621

	// 接受群组邀请
	AcceptGroupInviteCommand = 622
)