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
	insertId, err := u.GoDB.Insert("insert into user values(?, ?, ?, ?, ?, ?, ?, ?, ?)",
		nil, variable.IdWorker.GetId(), user.Username, "", user.Password, 0, "", time.Now(), time.Now())
	if err != nil {
		return 0, err
	}

	return insertId, nil
}

// 查询判断用户名，密码是否正确
func (u *UserDAO) QueryUserByNamePass(user model.User) (bool, *model.User) {
	// 查询
	row := u.GoDB.QueryRow("select * from user where username=? && password=?", user.Username, user.Password)

	err := row.Scan(&user.Id, &user.UserId, &user.Username, &user.Nickname, &user.Password, &user.Gender, &user.Extra, &user.CreateTime, &user.UpdateTime)
	if err != nil {
		//log.Info.Println("check user info err:", err)
		return false, nil
	}

	return true, &user
}

// 通过用户名查询用户
func (u *UserDAO) QueryUserByName(user model.User) (bool, *model.User) {
	// 查询
	row := u.GoDB.QueryRow("select * from user where username=?", user.Username)

	err := row.Scan(&user.Id, &user.UserId, &user.Username, &user.Nickname, &user.Password, &user.Gender, &user.Extra, &user.CreateTime, &user.UpdateTime)
	if err != nil {
		return false, nil
	}

	return true, &user
}

// 通过用户 id 查询用户
func (u *UserDAO) QueryUserById(user model.User) (bool, *model.User) {
	// 查询
	row := u.GoDB.QueryRow("select * from user where user_id=?", user.UserId)

	err := row.Scan(&user.Id, &user.UserId, &user.Username, &user.Nickname, &user.Password, &user.Gender, &user.Extra, &user.CreateTime, &user.UpdateTime)
	if err != nil {
		return false, nil
	}

	return true, &user
}