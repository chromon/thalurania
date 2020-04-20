package dao

import (
	"chalurania/comet/variable"
	"chalurania/service/db/conn"
	"chalurania/service/model"
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

func (u *UserDAO) AddUser(user model.User) (int64, error) {
	insertId, err := u.GoDB.Insert("insert into user values(?, ?, ?, ?, ?, ?, ?, ?)",
		nil, variable.IdWorker.GetId(), user.Nickname, user.Password, 0, "", time.Now(), time.Now())
	if err != nil {
		return 0, err
	}

	return insertId, nil
}
