package commands

import "flag"

// 命令分发
func CommandDistribute(m map[string]*flag.Flag) int32 {

	// 注册
	if RegCommand(m) {
		return 1
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