package model

// 群组离线消息
type GroupOfflineMessage struct {

	// 群组离线消息 id
	Id int64 `json:"id"`

	// 用户 id
	UserId int64 `json:"user_id"`

	// 群组 id
	GroupId int64 `json:"group_id"`

	// 最后一条消息的 id
	MessageId int64 `json:"message_id"`
}