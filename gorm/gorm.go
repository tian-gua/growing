package gorm

import (
	"database/sql"
	"errors"
	"reflect"
	"strings"
)

//数据库连接对象
var gdb *sql.DB = nil


//插入或者更新一条记录
//插入和更新取决于 id 字段是否为0
func Save(obj interface{}) error {

	//用于判断是否为插入方法
	var isInsert bool = false
	//用于存放sql字段
	var sqlStr string = ""

	//获得类型的信息
	t := reflect.TypeOf(obj).Elem()
	v := reflect.ValueOf(obj).Elem()
	tName := t.Name()

	//取id得值判断是insert 还是 update
	id := Parse(v.FieldByName("id"))
	if "0" == id {
		isInsert = true
	}
	if isInsert {
		//拼sql
		sqlStr = "insert into " + tName + "("
		var valueStr string
		//获得所有字段名
		numField := t.NumField()
		for i := 0; i < numField; i++ {
			fieldName := t.Field(i).Name
			sqlStr += fieldName + ","
			value := Parse(v.FieldByName(fieldName))
			//如果遇到id字段,则用default代替id的值, 实现自动自增
			if "id" == t.Field(i).Name {
				valueStr += "default,"
			} else {
				valueStr += value + ","
			}

		}
		//去掉右边的逗号
		sqlStr = strings.TrimRight(sqlStr, ",")
		sqlStr += ") values("
		sqlStr += valueStr
		sqlStr = strings.TrimRight(sqlStr, ",")
		sqlStr += ")"
	} else {
		sqlStr = "update " + tName + " set "
		//获得所有字段名
		numField := t.NumField()
		for i := 0; i < numField; i++ {

			fieldName := t.Field(i).Name
			if "id" == fieldName {
				continue
			} else {
				sqlStr += fieldName + "=" + Parse(v.FieldByName(fieldName)) + ","

			}

		}
		sqlStr = strings.TrimRight(sqlStr, ",") + "where id = " + id

	}

	//执行sql
	result, err := gdb.Exec(sqlStr)
	if err != nil {
		return err
	}
	rownum, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rownum == 0 {
		return errors.New("no record insert")
	}

	return nil

}


//删除一条记录
func Delete(obj interface{}) error {
	//用于存放sql字段
	var sqlStr string = ""
	//获得类型的信息
	t := reflect.TypeOf(obj).Elem()
	v := reflect.ValueOf(obj).Elem()
	tName := t.Name()

	//获得要删除的id
	id := Parse(v.FieldByName("id"))
	if "0" == id {
		errors.New("id is null")
	}
	//拼sql
	sqlStr = "delete from " + tName + " where id = " + id
	//执行sql
	result, err := gdb.Exec(sqlStr)
	if err != nil {
		return err
	}
	rownum, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rownum == 0 {
		return errors.New("no record delete")
	}

	return nil
}

func Query() {

}
