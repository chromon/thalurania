package dao

import (
	"chalurania/service/db/conn"
	"chalurania/service/log"
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
	row, err := u.GoDB.Query("select * from friend_request where user_id=? and del=0", user.UserId)
	return row, err
}

// 查询接收到的好友请求
func (u *FriendRequestDAO) QueryReceiveFriendReq(user model.User) (*sql.Rows, error) {
	// 查询
	row, err := u.GoDB.Query("select * from friend_request where friend_id=? and del=0", user.UserId)
	return row, err
}


// 查询收到的好友请求是否存在
func (u *FriendRequestDAO) QueryFriendReq(user, friend model.User) (bool, *model.FriendRequest) {
	// 查询
	row := u.GoDB.QueryRow("select * from friend_request where user_id=? && friend_id=? && del=0", user.UserId, friend.UserId)

	var friendRequest model.FriendRequest
	err := row.Scan(&friendRequest.Id, &friendRequest.UserId, &friendRequest.FriendId, &friendRequest.Del)
	if err != nil {
		return false, &friendRequest
	}

	return true, &friendRequest
}

// 更新用户请求
func (u *FriendRequestDAO) UpdateFriendReq(fr model.FriendRequest) int64 {
	affNum, err := u.GoDB.Update("update friend_request set del = 1 where id = ?", fr.Id)
	if err != nil {
		log.Info.Println("update friend request err:", err)
		return 0
	}
	return affNum
}