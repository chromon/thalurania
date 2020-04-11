package model

import "time"

// 用户
type User struct {
	// 用户主键 Id
	Id int64 `json:"id"`

	// 用户 Id
	UserId int64 `json:"user_id"`

	// 用户密码
	Password string `json:"password"`

	// 昵称
	Nickname string `json:"nickname"`

	// 性别
	Gender int32 `json:"gender"`

	// 附加信息
	Extra string `json:"extra"`

	// 创建时间
	CreateTime time.Time `json:"create_time"`

	// 更新时间
	UpdateTime time.Time `json:"update_time"`
}