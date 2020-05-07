package dao

import (
	"chalurania/service/db/conn"
	"chalurania/service/log"
	"chalurania/service/model"
	"database/sql"
)

// 群组邀请数据访问对象
type GroupInviteDAO struct {
	GoDB *conn.GoDB
}

// 创建好友请求访问对象
func NewGroupInviteDAO(goDB *conn.GoDB) *GroupInviteDAO {
	return &GroupInviteDAO{
		GoDB: goDB,
	}
}

// 插入新好友请求
func (u *GroupInviteDAO) AddGroupInvite(g model.GroupInvite) (int64, error) {
	insertId, err := u.GoDB.Insert("insert into group_invite values(?, ?, ?, ?, ?)",
		nil, g.UserId, g.FriendId, g.GroupId, g.Del)
	if err != nil {
		return 0, err
	}

	return insertId, nil
}

// 查询收到的群组邀请数量
func (u *GroupInviteDAO) QueryGroupInviteCount(user model.User) int64 {
	// 查询
	var count int64
	err := u.GoDB.QueryRow("select count(*) from group_invite where friend_id=? and del=0", user.UserId).Scan(&count)
	if err != nil {
		log.Error.Println("query group invite count err:", err)
		return 0
	}
	return count
}

// 查询接收到的群组邀请
func (u *GroupInviteDAO) QueryGroupInvite(user model.User) (*sql.Rows, error) {
	// 查询
	row, err := u.GoDB.Query("select * from group_invite where friend_id=? and del=0", user.UserId)
	return row, err
}
