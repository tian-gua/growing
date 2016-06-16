package gorm

import "database/sql"


//初始化数据库连接
func InitDB(db, connStr string) {
	//链接数据库
	_db, err := sql.Open(db, connStr)
	if err != nil {
		panic(err)
	}
	gdb = _db

}
