package commands

import (
	"chalurania/comet/constants"
	"flag"
)

// 命令分发
func CommandDistribute(m map[string]*flag.Flag) int32 {
	
	if RegCommand(m) {
		// 注册命令
		return constants.RegisterCommand
	} else if LoginCommand(m) {
		// 登录命令
		return constants.LoginCommand
	} else if LogoutCommand(m) {
		// 登出命令
		return constants.LogoutCommand
	} else if SearchUsernameCommand(m) {
		// 搜索用户名
		return constants.SearchUsernameCommand
	} else if SearchGroupIdCommand(m) {
		// 搜索用户组
		return constants.SearchGroupCommand
	} else if SearchUserIdCommand(m) {
		// 搜索用户 id
		return constants.SearchUserIdCommand
	} else if AddFriendByNameCommand(m) {
		// 通过用户名添加好友
		return constants.AddUserByNameCommand
	} else if AddFriendByIdCommand(m) {
		// 通过用户 id 添加好友
		return constants.AddUserByIdCommand
	} else if FriReqListCommand(m) {
		// 好友请求列表
		return constants.FriendReqListCommand
	} else if AcceptFriendByNameCommand(m) {
		// 通过用户名接受好友请求
		return constants.AcceptFriendByNameCommand
	} else if AcceptFriendByIdCommand(m) {
		// 通过用户 id 接受好友请求
		return constants.AcceptFriendByIdCommand
	} else if FriendListCommand(m) {
		// 好友列表
		return constants.FriendListCommand
	} else if ChatWithFriendByNameCommand(m) {
		// 通过用户名与好友私聊
		return constants.ChatWithFriendByNameCommand
	} else if ChatWithFriendByIdCommand(m) {
		// 通过用户 id 与好友私聊
		return constants.ChatWithFriendByIdCommand
	} else if OfflineMsgByNameCommand(m) {
		// 通过用户名查询离线消息
		return constants.OfflineMsgByNameCommand
	} else if OfflineMsgByIdCommand(m) {
		// 通过用户 id 查询离线消息
		return constants.OfflineMsgByIdCommand
	} else if CreateGroupCommand(m) {
		// 创建群组
		return constants.CreateGroupCommand
	} else if GroupInviteByNameCommand(m) {
		// 通过用户名邀请到群组
		return constants.GroupInviteByNameCommand
	} else if GroupInviteByIdCommand(m) {
		// 通过用户 id 邀请到群组
		return constants.GroupInviteByIdCommand
	} else if GroupInviteListCommand(m) {
		// 群组邀请列表
		return constants.GroupInviteListCommand
	} else if AcceptGroupInviteCommand(m) {
		// 接受群组邀请
		return constants.AcceptGroupInviteCommand
	} else if GroupMemberCommand(m) {
		// 群成员列表
		return constants.GroupMemberCommand
	} else if GroupListCommand(m) {
		// 已加入的群组列表
		return constants.GroupListCommand
	} else if GroupChatCommand(m) {
		// 群组聊天
		return constants.GroupChatCommand
	}

	// 无效命令
	return constants.ErrorCommand
}

// 注册命令
func RegCommand(m map[string]*flag.Flag) bool {
	_, r := m["r"]
	_, u := m["u"]
	_, p := m["p"]

	if r && u && p {
		return true
	}
	return false
}

// 登录命令
func LoginCommand(m map[string]*flag.Flag) bool {
	_, l := m["l"]
	_, u := m["u"]
	_, p := m["p"]

	if l && u && p {
		return true
	}
	return false
}

// 登出命令
func LogoutCommand(m map[string]*flag.Flag) bool {
	_, q := m["q"]
	return q
}

// 搜索用户名命令
func SearchUsernameCommand(m map[string]*flag.Flag) bool {
	_, s := m["s"]
	_, u := m["u"]

	if s && u {
		return true
	}
	return false
}

// 搜索用户 id 命令
func SearchUserIdCommand(m map[string]*flag.Flag) bool {
	_, s := m["s"]
	_, n := m["n"]

	if s && n {
		return true
	}
	return false
}

// 搜索用户名命令
func SearchGroupIdCommand(m map[string]*flag.Flag) bool {
	_, s := m["s"]
	_, g := m["g"]
	_, gn := m["gn"]

	if s && g && gn {
		return true
	}
	return false
}

// 通过用户名添加好友
func AddFriendByNameCommand(m map[string]*flag.Flag) bool {
	_, add := m["add"]
	_, u := m["u"]

	if add && u {
		return true
	}
	return false
}

// 通过用户名添加好友
func AddFriendByIdCommand(m map[string]*flag.Flag) bool {
	_, add := m["add"]
	_, n := m["n"]

	if add && n {
		return true
	}
	return false
}

// 好友请求列表
func FriReqListCommand(m map[string]*flag.Flag) bool {
	_, fr := m["req"]
	_, list := m["list"]

	if fr && list {
		return true
	}
	return false
}

// 通过用户名接受好友请求
func AcceptFriendByNameCommand(m map[string]*flag.Flag) bool {
	_, accept := m["accept"]
	_, u := m["u"]

	if accept && u {
		return true
	}
	return false
}

// 通过用户 id 接受好友请求
func AcceptFriendByIdCommand(m map[string]*flag.Flag) bool {
	_, accept := m["accept"]
	_, n := m["n"]

	if accept && n {
		return true
	}
	return false
}

// 好友列表
func FriendListCommand(m map[string]*flag.Flag) bool {
	_, f := m["f"]
	_, list := m["list"]

	if f && list {
		return true
	}
	return false
}

// 通过用户名私聊
func ChatWithFriendByNameCommand(m map[string]*flag.Flag) bool {
	_, chat := m["chat"]
	_, u := m["u"]
	_, msg := m["m"]

	if chat && u && msg {
		return true
	}
	return false
}

// 通过用户 id 私聊
func ChatWithFriendByIdCommand(m map[string]*flag.Flag) bool {
	_, chat := m["chat"]
	_, n := m["n"]
	_, msg := m["m"]

	if chat && n && msg {
		return true
	}
	return false
}

// 通过用户名查询离线消息
func OfflineMsgByNameCommand(m map[string]*flag.Flag) bool {
	_, o := m["o"]
	_, u := m["u"]

	if o && u {
		return true
	}
	return false
}

// 通过用户 id 查询离线消息
func OfflineMsgByIdCommand(m map[string]*flag.Flag) bool {
	_, o := m["o"]
	_, n := m["n"]

	if o && n {
		return true
	}
	return false
}

// 创建群组
func CreateGroupCommand(m map[string]*flag.Flag) bool {
	_, add := m["add"]
	_, g := m["g"]

	if add && g {
		return true
	}
	return false
}

// 群组通过用户名邀请
func GroupInviteByNameCommand(m map[string]*flag.Flag) bool {
	_, i := m["i"]
	_, u := m["u"]
	_, gn := m["gn"]

	if i && u && gn {
		return true
	}
	return false
}

// 群组通过用户 id 邀请
func GroupInviteByIdCommand(m map[string]*flag.Flag) bool {
	_, i := m["i"]
	_, n := m["n"]
	_, gn := m["gn"]

	if i && n && gn {
		return true
	}
	return false
}

// 群组邀请列表
func GroupInviteListCommand(m map[string]*flag.Flag) bool {
	_, g := m["g"]
	_, i := m["i"]
	_, list := m["list"]

	if g && i && list {
		return true
	}
	return false
}

// 通过群组 id 接受群组邀请
func AcceptGroupInviteCommand(m map[string]*flag.Flag) bool {
	_, accept := m["accept"]
	_, gn := m["gn"]

	if accept && gn {
		return true
	}
	return false
}

// 通过群组 id 查询群成员列表
func GroupMemberCommand(m map[string]*flag.Flag) bool {
	_, g := m["g"]
	_, gn := m["gn"]
	_, list := m["list"]

	if g && gn && list {
		return true
	}
	return false
}

// 已加入的群组列表
func GroupListCommand(m map[string]*flag.Flag) bool {
	_, g := m["g"]
	_, list := m["list"]

	if g && list {
		return true
	}
	return false
}

// 群组聊天
func GroupChatCommand(m map[string]*flag.Flag) bool {
	_, chat := m["chat"]
	_, gn := m["gn"]
	_, msg := m["m"]

	if chat && gn && msg {
		return true
	}
	return false
}