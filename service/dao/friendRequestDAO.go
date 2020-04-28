package dao

import (
	"chalurania/service/db/conn"
	"chalurania/service/model"
)

// 好友请求数据访问对象
type FriendRequestDAO struct {
	GoDB *conn.GoDB
}

// 创建好友请求访问对象
func NewFriendRequestDAO(goDB *conn.GoDB) *FriendRequestDAO {
	return &FriendRequestDAO{
		GoDB: goDB,
	}
}

// 插入新好友请求
func (u *FriendRequestDAO) AddFriendRequest(fr model.FriendRequest) (int64, error) {
	insertId, err := u.GoDB.Insert("insert into friend_request values(?, ?, ?)",
		nil, fr.UserId, fr.FriendId)
	if err != nil {
		return 0, err
	}

	return insertId, nil
}