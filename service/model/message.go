package model

import "time"

// 消息
type Message struct {

	// 消息主键 Id
	Id int64 `json:"id"`

	// 消息序列号
	Seq int64 `json:"seq"`

	// 消息内容
	Content string `json:"content"`

	// 消息所属类型 Id
	MessageTypeId int64 `json:"message_type_id"`

	// 发送者类型
	SenderType int32 `json:"sender_type"`

	// 发送者 Id
	SenderId int64 `json:"sender_id"`

	// 接受者类型
	ReceiverType int32 `json:"receiver_type"`

	// 接受者 Id
	ReceiverId int64 `json:"receiver_id"`

	// at 用户 id
	ToUserIds string `json:"to_user_ids"`

	// 消息发送时间
	SendTime time.Time `json:"send_time"`

	// 消息状态
	Status int32 `json:"status"`

	// 消息创建时间
	CreateTime time.Time `json:"create_time"`

	// 消息更新时间
	UpdateTime time.Time `json:"update_time"`
}