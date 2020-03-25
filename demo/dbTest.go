package main

import (
	"chalurania/db/conn"
	"chalurania/service/log"
	"strings"
)

const (
	userName = "root"
	password = "root"
	ip       = "127.0.0.1"
	port     = "3306"
	dbName   = "thalurania"
)

func main() {

	path := strings.Join([]string{userName, ":", password, "@tcp(", ip, ":", port, ")/", dbName, "?charset=utf8"}, "")

	db, err := conn.NewDB("mysql", path)
	if err != nil {
		return
	}
	defer db.Close()

	// 插入
	//insertId, err := db.Insert("insert into user values(?, ?, ?, ?, ?, ?, ?)", nil, 101, "ellery", 1, "xxx", time.Now(), time.Now())
	//if err != nil {
	//	return
	//}
	//
	//log.Info.Println("insert Id：", insertId)

	// 查询
	//row, _ := db.Query("select id, nickname from user where id=?", 1)
	//defer row.Close()
	//
	//for row.Next() {
	//	id := 0
	//	nickname := ""
	//	err = row.Scan(&id, &nickname)
	//	if err != nil {
	//		return
	//	}
	//	log.Info.Println("id", id, "nickname", nickname)
	//}

	// 更新
	//affNum, err := db.Update("update user set user_id = ? where id = ?", 102, 1)
	//if err != nil {
	//	return
	//}
	//
	//log.Info.Println("affNum:", affNum)

	// 删除
	affNum, err := db.Delete("delete from user where id = ?", 1)
	if err != nil {
		return
	}
	log.Info.Println("affNum:", affNum)
}