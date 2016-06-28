package gorm

import (
	"database/sql"
	"errors"
	"reflect"
)

//数据库连接对象
var gdb *sql.DB = nil


//插入或者更新一条记录
//插入和更新取决于 id 字段是否为0
func Save(obj interface{}) error {

	//生成sql
	sqlStr := parseSaveSql(obj)

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

	//生成sql
	sqlStr := parseDeleteSql(obj)
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


//查询记录
func Query(obj, target interface{}) error {

	t := reflect.TypeOf(obj).Elem()
	targetVlaue := reflect.ValueOf(target).Elem()

	sqlStr := parseQuerySql(obj)

	//查询
	rows, err := gdb.Query(sqlStr)
	if err != nil {
		return err
	}

	//获得所有列
	columns, err := rows.Columns()
	if err != nil {
		return err
	}
	//获得列的数量
	colNum := len(columns)

	values := make([]sql.RawBytes, colNum)
	scans := make([]interface{}, colNum)
	//封装
	for i := range values {
		scans[i] = &values[i]
	}
	var index int = 0
	//遍历所有记录
	for rows.Next() {
		err := rows.Scan(scans...)
		if err != nil {
			return err
		}
		//根据反射来新建一个和记录对应的对象
		var newV = reflect.New(t).Elem()
		for i := 0; i < colNum; i++ {

			setValue(newV.Field(i), values[i])
		}
		targetVlaue = reflect.Append(targetVlaue, newV)
		index++
	}
	//更新target的值
	reflect.ValueOf(target).Elem().Set(targetVlaue.Slice(0, index))
	return nil

}






//关闭DB对象
func CloseDB() {
	gdb.Close()
}





////将值v设置到结构体s里
//func setValue(s, v interface{}, fieldNum int) {
//	//获得类型的信息
//	value := reflect.ValueOf(s).Elem()
//
//	//根据结构体s的字段类型来强转v
//	fieldType := value.Field(fieldNum).Kind()
//	switch fieldType {
//	case reflect.String:
//		value.Field(fieldNum).Set(string(v))
//	case reflect.Int:
//		value.Field(fieldNum).Set(int(v))
//	}
//}