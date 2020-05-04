package dao

import (
	"chalurania/service/db/conn"
	"chalurania/service/log"
	"chalurania/service/model"
	"time"
)

// 群组数据访问对象
type GroupDAO struct {
	GoDB *conn.GoDB
}

// 创建群组访问对象
func NewGroupDAO(goDB *conn.GoDB) *GroupDAO {
	return &GroupDAO{
		GoDB: goDB,
	}
}

// 创建新群组
func (u *GroupDAO) AddGroup(g model.Group) (int64, error) {
	insertId, err := u.GoDB.Insert("insert into im_group values(?, ?, ?, ?, ?, ?, ?, ?, ?)",
		nil, g.GroupId, g.Name, g.Introduction, g.UserCount, g.Type, g.Extra, time.Now(), time.Now())
	if err != nil {
		return 0, err
	}

	return insertId, nil
}

// 通过群组 id 查询
func (u *GroupDAO) QueryGroupById(g model.Group) (bool, *model.Group) {
	// 查询
	row := u.GoDB.QueryRow("select * from im_group where id=?", g.Id)

	err := row.Scan(&g.Id, &g.GroupId, &g.Name, &g.Introduction, &g.UserCount, &g.Type, &g.Extra, &g.CreateTime, &g.UpdateTime)
	if err != nil {
		log.Error.Println("query group by id err:", err)
		return false, nil
	}

	return true, &g
}

// 通过群组 id 查询
func (u *GroupDAO) QueryGroupByGroupId(g model.Group) (bool, *model.Group) {
	// 查询
	row := u.GoDB.QueryRow("select * from im_group where group_id=?", g.GroupId)

	err := row.Scan(&g.Id, &g.GroupId, &g.Name, &g.Introduction, &g.UserCount, &g.Type, &g.Extra, &g.CreateTime, &g.UpdateTime)
	if err != nil {
		log.Error.Println("query group by group id err:", err)
		return false, nil
	}

	return true, &g
}