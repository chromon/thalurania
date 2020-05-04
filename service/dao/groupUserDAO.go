package dao

import (
	"chalurania/service/db/conn"
	"chalurania/service/model"
	"time"
)

// 群组用户数据访问对象
type GroupUserDAO struct {
	GoDB *conn.GoDB
}

// 创建群组用户访问对象
func NewGroupUserDAO(goDB *conn.GoDB) *GroupUserDAO {
	return &GroupUserDAO{
		GoDB: goDB,
	}
}

// 创建新群组用户
func (gu *GroupUserDAO) AddGroupUser(g model.Group, u model.User) (int64, error) {
	insertId, err := gu.GoDB.Insert("insert into group_user values(?, ?, ?, ?, ?, ?, ?)",
		nil, g.GroupId, u.UserId, u.Username, "", time.Now(), time.Now())
	if err != nil {
		return 0, err
	}

	return insertId, nil
}