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
	} else if SearchUserIdCommand(m) {
		// 搜索用户 id
		return constants.SearchUserIdCommand
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

// 搜索用户名命令
func SearchUserIdCommand(m map[string]*flag.Flag) bool {
	_, s := m["s"]
	_, n := m["n"]

	if s && n {
		return true
	}
	return false
}