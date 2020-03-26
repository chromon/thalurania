package conn

import (
	"chalurania/service/log"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

// 数据库连接
type GoDB struct {
	*sql.DB
}

func NewDB(driverName, dataSource string) (*GoDB, error) {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return nil, err
	}

	// 设置数据库最大连接数
	db.SetConnMaxLifetime(100)
	// 设置上数据库最大闲置连接数，连接池中的保持连接的最大连接数
	db.SetMaxIdleConns(10)

	goDB := &GoDB{DB: db}
	return goDB, nil
}

// 查询数据
func (g *GoDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	result, err := g.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// 插入数据
func (g *GoDB) Insert(sql string, args ...interface{}) (int64, error) {
	// 开启事务
	tx, err := g.DB.Begin()
	if err != nil {
		log.Error.Println("Transaction begin err:", err)
		return 0, err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(sql)
	if err != nil {
		log.Error.Println("Statement prepare err:", err)
		return 0, err
	}
	defer stmt.Close()

	// 执行
	result, err := stmt.Exec(args...)
	if err != nil {
		log.Error.Println("Statement exec err:", err)
		return 0, err
	}

	// 提交事务
	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	// 返回插入自增 Id
	lastId, _ := result.LastInsertId()
	return lastId, nil
}

// 更新数据
func (g *GoDB) Update(sql string, args ...interface{}) (int64, error) {
	tx, err := g.DB.Begin()
	if err != nil {
		log.Error.Println("Transaction begin err:", err)
		return 0, err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(sql)
	if err != nil {
		log.Error.Println("Statement prepare err:", err)
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(args...)
	if err != nil {
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	affectedNum, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return affectedNum, nil
}

// 删除数据
func (g *GoDB) Delete(sql string, args ...interface{}) (int64, error) {
	tx, err := g.DB.Begin()
	if err != nil {
		log.Error.Println("Transaction begin err:", err)
		return 0, err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(sql)
	if err != nil {
		log.Error.Println("Statement prepare err:", err)
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(args...)
	if err != nil {
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	affectedNum, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return affectedNum, nil
}
