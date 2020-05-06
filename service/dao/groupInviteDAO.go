package dao

import (
	"chalurania/service/db/conn"
	"chalurania/service/model"
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
