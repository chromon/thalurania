package main

import (
	"chalurania/service/db/conn"
	"chalurania/service/log"
	"strings"
	"testing"
	"time"
)

const (
	userName = "root"
	password = "root"
	ip       = "127.0.0.1"
	port     = "3306"
	dbName   = "thalurania"
)

func TestInsert(t *testing.T) {
	dataSource := strings.Join([]string{userName, ":", password, "@tcp(", ip, ":", port, ")/", dbName, "?charset=utf8"}, "")

	db, err := conn.NewDB("mysql", dataSource)
	if err != nil {
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			panic(err)
		}
	}()

	// 插入
	insertId, err := db.Insert("insert into user values(?, ?, ?, ?, ?, ?, ?)", nil, 101, "ellery", 1, "xxx", time.Now(), time.Now())
	if err != nil {
		return
	}
	log.Info.Println("insert Id：", insertId)
}

func TestQuery(t *testing.T) {
	dataSource := strings.Join([]string{userName, ":", password, "@tcp(", ip, ":", port, ")/", dbName, "?charset=utf8"}, "")

	db, err := conn.NewDB("mysql", dataSource)
	if err != nil {
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			panic(err)
		}
	}()

	// 查询
	row, _ := db.Query("select id, nickname from user where id=?", 1)
	defer func() {
		if err := row.Close(); err != nil {
			panic(err)
		}
	}()

	for row.Next() {
		id := 0
		nickname := ""
		err = row.Scan(&id, &nickname)
		if err != nil {
			return
		}
		log.Info.Println("id", id, "nickname", nickname)
	}
}

func TestUpdate(t *testing.T) {
	dataSource := strings.Join([]string{userName, ":", password, "@tcp(", ip, ":", port, ")/", dbName, "?charset=utf8"}, "")

	db, err := conn.NewDB("mysql", dataSource)
	if err != nil {
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			panic(err)
		}
	}()

	// 更新
	affNum, err := db.Update("update user set user_id = ? where id = ?", 102, 1)
	if err != nil {
		return
	}

	log.Info.Println("affNum:", affNum)
}

func TestDelete(t *testing.T) {
	dataSource := strings.Join([]string{userName, ":", password, "@tcp(", ip, ":", port, ")/", dbName, "?charset=utf8"}, "")

	db, err := conn.NewDB("mysql", dataSource)
	if err != nil {
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			panic(err)
		}
	}()

	// 删除
	affNum, err := db.Delete("delete from user where id = ?", 1)
	if err != nil {
		return
	}
	log.Info.Println("affNum:", affNum)
}
