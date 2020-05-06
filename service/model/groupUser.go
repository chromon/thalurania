package model

import "time"

// 群成员
type GroupUser struct {
	// 群组成员组件 Id
	Id int64 `json:"id"`

	// 群组 Id
	GroupId int64 `json:"group_id"`

	// 成员 Id
	UserId int64 `json:"user_id"`

	// 成员在群组中昵称
	Label string `json:"label"`

	// 附加属性
	Extra string `json:"extra"`

	// 创建时间
	CreateTime time.Time `json:"create_time"`

	// 更新时间
	UpdateTime time.Time `json:"update_time"`
}