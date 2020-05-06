package model

type GroupInvite struct {
	// 主键 Id
	Id int64 `json:"id"`
	
	// 用户 id
	UserId int64 `json:"user_id"`
	
	// 好友 id
	FriendId int64 `json:"friend_id"`
	
	// 群组 id
	GroupId int64 `json:"group_id"`
	
	// 是否已删除 0：否，1：是
	Del int32 `json:"del"`
}