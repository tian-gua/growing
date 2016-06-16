package utils

import (
	"database/sql"
	"fmt"
	"strings"
)


//默认数据库连接信息
var db string = "mysql"
var connStr string = ""

func Generate(tableName string) {
	//链接数据库
	db, err := sql.Open(db, connStr)
	if err != nil {
		panic(err)
	}

	//关闭链接
	defer db.Close()

	//查询第一条数据
	sqlString := "desc " + tableName
	rows, err := db.Query(sqlString)
	if err != nil {
		panic(err)
	}
	defer rows.Close()


	//获得表所有字段
	columns, err := rows.Columns()
	if err != nil {
		panic(err)
	}

	colNum := len(columns)

	//定义一个RawBytes切片接收所有字段的值
	values := make([]sql.RawBytes, colNum)
	//定义一个空接口切片用于封装RawBytes切片,Scan方法只能接收interface{}
	scans := make([]interface{}, colNum)
	//封装
	for i := range values {
		scans[i] = &values[i]
	}
	//用于存放所有字段
	fields := make([]string, 0)
	//存放所有字段类型
	fieldTypes := make([]string, 0)

	for rows.Next() {

		//读取所有的字段到 空结构体切片里
		err := rows.Scan(scans...)
		if err != nil {
			fmt.Println(err)
		}

		//存放每一条记录的第一个字段(表的字段名) 到fields切片里
		fields = append(fields, string(values[0]))
		ftype := string(values[1])
		switch true {
		case strings.HasPrefix(ftype, "varchar"):
			ftype = "string"
		case strings.HasPrefix(ftype, "int"):
			ftype = "int"
		case strings.HasPrefix(ftype, "decimal"):
			ftype = "int64"
		case strings.HasPrefix(ftype, "timestamp") || strings.HasPrefix(ftype, "datetime"):
			ftype = "time.Time"

		}
		fieldTypes = append(fieldTypes, ftype)
	}

	structString := "type " + tableName + " struct{\n"
	for i, v := range fields {
		structString += "\t" + v + "\t" + fieldTypes[i] + "\n"
	}
	structString += "}"
	fmt.Println(structString)

}

//设置数据库连接信息
func SetDBInfo(_db, _connStr string) {
	db = _db
	connStr = _connStr
}