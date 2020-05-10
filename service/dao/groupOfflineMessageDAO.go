package dao

import (
	"chalurania/service/db/conn"
	"chalurania/service/model"
)

// 群组离线消息
type GroupOfflineMessageDAO struct {
	GoDB *conn.GoDB
}

// 创建群组离线消息访问对象
func NewGroupOfflineMessageDAO(goDB *conn.GoDB) *GroupOfflineMessageDAO {
	return &GroupOfflineMessageDAO{
		GoDB: goDB,
	}
}

// 添加群组离线消息
func (u *GroupOfflineMessageDAO) AddGroupOfflineMessage(g model.GroupOfflineMessage) (int64, error) {
	insertId, err := u.GoDB.Insert("insert into group_offline_message values(?, ?, ?, ?)",
		nil, g.UserId, g.GroupId, g.MessageId)
	if err != nil {
		return 0, err
	}

	return insertId, nil
}