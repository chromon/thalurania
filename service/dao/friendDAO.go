package dao

import (
	"chalurania/service/db/conn"
	"chalurania/service/log"
	"chalurania/service/model"
	"database/sql"
)

// 好友数据访问对象
type FriendDAO struct {
	GoDB *conn.GoDB
}

// 创建好友访问对象
func NewFriendDAO(goDB *conn.GoDB) *FriendDAO {
	return &FriendDAO{
		GoDB: goDB,
	}
}

// 插入新好友
func (u *FriendDAO) AddFriend(f model.Friend) int64 {
	insertId, err := u.GoDB.Insert("insert into friend values(?, ?, ?)",
		nil, f.UserId, f.FriendId)
	if err != nil {
		log.Info.Println("insert friend err:", err)
		return 0
	}

	return insertId
}

// 查询接收到的好友请求
func (u *FriendDAO) QueryFriend(user model.User) (*sql.Rows, error) {
	// 查询
	rows, err := u.GoDB.Query("select friend_id from friend where user_id=?", user.UserId)
	return rows, err
}

// 查询双方是否是好友
func (u *FriendDAO) QueryFriendById(user, friend model.User) bool {
	// 查询
	row := u.GoDB.QueryRow("select * from friend where user_id=? && friend_id=?", user.UserId, friend.UserId)

	var friendship model.Friend

	err := row.Scan(&friendship.Id, &friendship.UserId, &friendship.FriendId)
	if err != nil {
		log.Error.Println("query friend err:", err)
		return false
	}

	return true
}