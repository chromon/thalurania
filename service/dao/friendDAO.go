package dao

import (
	"chalurania/service/db/conn"
	"chalurania/service/log"
	"chalurania/service/model"
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