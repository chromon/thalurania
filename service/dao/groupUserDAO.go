package dao

import (
	"chalurania/service/db/conn"
	"chalurania/service/log"
	"chalurania/service/model"
	"database/sql"
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

// 查询用户是否是群组成员
func (gu *GroupUserDAO) QueryGroupUserById(user model.User, group model.Group) (bool, *model.GroupUser) {
	// 查询
	row := gu.GoDB.QueryRow("select * from group_user where group_id =? && user_id=?", group.GroupId, user.UserId)

	var g model.GroupUser
	err := row.Scan(&g.Id, &g.GroupId, &g.UserId, &g.Label, &g.Extra, &g.CreateTime, &g.UpdateTime)
	if err != nil {
		return false, nil
	}

	return true, &g
}

// 查询收到的群组邀请数量
func (gu *GroupUserDAO) QueryGroupUserCount(group model.Group) int64 {
	// 查询
	var count int64
	err := gu.GoDB.QueryRow("select count(*) from group_user where group_id=?", group.GroupId).Scan(&count)
	if err != nil {
		log.Error.Println("query group user count err:", err)
		return 0
	}
	return count
}

// 查询群组成员列表
func (gu *GroupUserDAO) QueryGroupUsers(group model.Group) (*sql.Rows, error) {
	// 查询
	row, err := gu.GoDB.Query("select * from group_user where group_id=?", group.GroupId)
	return row, err
}

// 查询用户所在的群组数量
func (gu *GroupUserDAO) QueryGroupCountByUser(user model.User) int64 {
	// 查询
	var count int64
	err := gu.GoDB.QueryRow("select count(*) from group_user where user_id=?", user.UserId).Scan(&count)
	if err != nil {
		log.Error.Println("query group count by user err:", err)
		return 0
	}
	return count
}

// 查询用户所在的群组
func (gu *GroupUserDAO) QueryGroupByUser(user model.User) (*sql.Rows, error) {
	// 查询
	row, err := gu.GoDB.Query("select * from group_user where user_id=?", user.UserId)
	return row, err
}