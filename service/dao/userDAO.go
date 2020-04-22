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

// 插入新用户
func (u *UserDAO) AddUser(user model.User) (int64, error) {
	insertId, err := u.GoDB.Insert("insert into user values(?, ?, ?, ?, ?, ?, ?, ?)",
		nil, variable.IdWorker.GetId(), user.Nickname, user.Password, 0, "", time.Now(), time.Now())
	if err != nil {
		return 0, err
	}

	return insertId, nil
}

// 查询判断用户名，密码是否正确
func (u *UserDAO) QueryUserByNamePass(user model.User) (bool, *model.User) {
	// 查询
	row, _ := u.GoDB.QueryRow("select * from user where nickname=? && password=?", user.Nickname, user.Password)
	defer func() {
		if err := row.Close(); err != nil {
			panic(err)
		}
	}()

	err := row.Scan(&user.Id, &user.UserId, &user.Nickname, &user.Password, &user.Gender, &user.Extra, &user.CreateTime, &user.UpdateTime)
	if err != nil {
		return false, nil
	}

	return true, &user
}