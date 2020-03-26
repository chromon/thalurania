package dao

import (
	"chalurania/service/db/conn"
	"time"
)

// 用户数据访问对象
type UserDAO struct {
	GoDB *conn.GoDB
}

// 创建用户访问对象
func NewUserDAO(goDB *conn.GoDB) *UserDAO {
	return &UserDAO{
		GoDB: goDB,
	}
}

func (u *UserDAO) AddUser() (int64, error) {
	insertId, err := u.GoDB.Insert("insert into user values(?, ?, ?, ?, ?, ?, ?)", nil, 101, "ellery", 1, "xxx", time.Now(), time.Now())
	if err != nil {
		return 0, err
	}

	return insertId, nil
}