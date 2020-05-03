package dao

import (
	"chalurania/service/db/conn"
	"chalurania/service/log"
	"chalurania/service/model"
	"database/sql"
)

// 消息数据访问对象
type MessageDAO struct {
	GoDB *conn.GoDB
}

// 创建消息访问对象
func NewMessageDAO(goDB *conn.GoDB) *MessageDAO {
	return &MessageDAO{
		GoDB: goDB,
	}
}

// 插入新消息
func (u *MessageDAO) AddMessage(m model.Message) (int64, error) {
	insertId, err := u.GoDB.Insert("insert into message values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		nil, m.Seq, m.Content, m.MessageTypeId, m.SenderType, m.SenderId, m.ReceiverType, m.ReceiverId, "",
		m.SendTime, m.Status, m.CreateTime, m.UpdateTime)
	if err != nil {
		return 0, err
	}

	return insertId, nil
}

// 查询未读消息数量
func (u *MessageDAO) QueryOfflineMsgCount(user, friend model.User) int64 {
	// 查询
	var count int64
	err := u.GoDB.QueryRow("select count(*) from message where sender_id=? and receiver_id=? and status=1", user.UserId, friend.UserId).Scan(&count)
	if err != nil {
		log.Error.Println("query offline message count err:", err)
		return 0
	}
	return count
}

// 查询未读消息
func (u *MessageDAO) QueryOfflineMessage(user, friend model.User) (*sql.Rows, error) {
	// 查询
	rows, err := u.GoDB.Query("select * from message where sender_id=? and receiver_id=? and status=1", user.UserId, friend.UserId)
	return rows, err
}