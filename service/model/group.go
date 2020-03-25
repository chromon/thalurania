package model

import "time"

// 群组
type Group struct {
	// 群组主键 Id
	Id int64 `json:"id"`

	// 群组 Id
	GroupId int64 `json:"group_id"`

	// 群组名称
	Name string `json:"name"`

	// 群简介
	Introduction string `json:"introduction"`

	// 群成员数量
	UserCount int32 `json:"user_count"`

	// 群组类型
	Type int32 `json:"type"`

	// 群附加属性
	Extra string `json:"extra"`

	// 创建时间
	CreateTime time.Time `json:"create_time"`

	// 更新时间
	UpdateTime time.Time `json:"update_time"`
}