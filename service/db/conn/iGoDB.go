package conn

import "database/sql"

type IGoDB interface {
	// 查询数据
	Query(string, ...interface{}) (*sql.Rows, error)

	// 插入数据
	Insert(string, ...interface{}) (int64, error)

	// 更新数据
	Update(string, ...interface{}) (int64, error)

	// 删除数据
	Delete(string, ...interface{}) (int64, error)
}
