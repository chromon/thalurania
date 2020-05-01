package dao

import (
	"chalurania/service/db/conn"
	"chalurania/service/model"
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