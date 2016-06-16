package gorm

import (
	"database/sql"
	"reflect"
	"fmt"
	"strings"
)

//数据库连接对象
var gdb *sql.DB = nil

func Update() {

}

//保存一条记录
func Insert(obj interface{}) {

	//获得类型的信息
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	tName := t.Name()
	//拼sql
	sqlStr := "insert into " + tName + "("
	var valueStr string
	//获得所有字段名
	numField := t.NumField()
	for i := 0; i < numField; i++ {
		sqlStr += t.Field(i).Name + ","
		value := Parse(v.Field(i))
		valueStr += value + ","

	}
	//去掉右边的逗号
	sqlStr = strings.TrimRight(sqlStr, ",")
	sqlStr += ") values("
	sqlStr += valueStr
	sqlStr = strings.TrimRight(sqlStr, ",")
	sqlStr += ")"


	fmt.Println(sqlStr)
}

func Delete() {

}

func Query() {

}

func init() {

}
