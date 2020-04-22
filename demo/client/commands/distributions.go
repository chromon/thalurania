package commands

import "flag"

// 命令分发
func CommandDistribute(m map[string]*flag.Flag) int32 {

	// 注册
	if RegCommand(m) {
		return 1
	} else if LoginCommand(m) {
		return 2
	}

	return 0
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