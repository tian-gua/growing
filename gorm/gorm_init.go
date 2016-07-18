package gorm

import (
	"database/sql"
)


//初始化数据库连接
func InitDB(dialect, connStr string) {
	//链接数据库
	db, err := sql.Open(dialect, connStr)
	if err != nil {
		panic(err)
	}
	gdb = db
}