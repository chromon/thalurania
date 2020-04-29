package dao

import (
	"chalurania/service/db/conn"
	"chalurania/service/model"
	"database/sql"
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
	insertId, err := u.GoDB.Insert("insert into friend_request values(?, ?, ?, ?)",
		nil, fr.UserId, fr.FriendId, 0)
	if err != nil {
		return 0, err
	}

	return insertId, nil
}

// 查询已发送的好友请求
func (u *FriendRequestDAO) QuerySentFriendReq(user model.User) (*sql.Rows, error) {
	// 查询
	row, err := u.GoDB.Query("select * from friend_request where user_id=?", user.UserId)
	return row, err
}

// 查询接收到的好友请求
func (u *FriendRequestDAO) QueryReceiveFriendReq(user model.User) (*sql.Rows, error) {
	// 查询
	row, err := u.GoDB.Query("select * from friend_request where friend_id=?", user.UserId)
	return row, err
}