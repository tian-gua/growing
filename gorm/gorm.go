package gorm

import (
	"database/sql"
	"errors"
	"reflect"
	"strings"
	"strconv"
	"fmt"
	"utils"
	"time"
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
	//结构体首字母转换为小写
	//结构体首字母大写是为了供其他包访问,数据库则不用
	tName = strings.ToLower(tName)

	//取id得值判断是insert 还是 update
	id := utils.Parse(v.FieldByName("Id"))
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
			unCamelFieldName, _ := utils.UnCamelCase(fieldName)
			sqlStr += unCamelFieldName + ","
			value := utils.Parse(v.FieldByName(fieldName))
			//如果遇到id字段,则用default代替id的值, 实现自动自增
			if "Id" == t.Field(i).Name {
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
			unCamelFieldName, _ := utils.UnCamelCase(fieldName)
			if "Id" == fieldName {
				continue
			} else {

				sqlStr += unCamelFieldName + "=" + utils.Parse(v.FieldByName(fieldName)) + ","

			}

		}
		sqlStr = strings.TrimRight(sqlStr, ",") + " where id = " + id

	}

	fmt.Println("[sql-gorm-" + utils.DateFormat(time.Now(), "yyyy-MM-dd HH:mm:ss") + "]:" + sqlStr)

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
	//结构体首字母转换为小写
	//结构体首字母大写是为了供其他包访问,数据库则不用
	tName = strings.ToLower(tName)

	//获得要删除的id
	id := utils.Parse(v.FieldByName("Id"))
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


//查询记录
func Query(obj interface{}, target interface{}) error {

	var sqlStr string = ""
	var whereStr string = ""
	//获得类型的信息
	t := reflect.TypeOf(obj).Elem()
	v := reflect.ValueOf(obj).Elem()
	targetVlaue := reflect.ValueOf(target).Elem()

	tName := t.Name()
	//结构体首字母转换为小写
	//结构体首字母大写是为了供其他包访问,数据库则不用
	tName = strings.ToLower(tName)

	sqlStr = "select "

	//拼接需要查询的字段
	for i := 0; i < t.NumField(); i++ {
		fieldName := t.Field(i).Name
		fieldValue := v.FieldByName(fieldName)
		//字符串反驼峰转换,例如 UserName 会变成 user_name
		fieldName, _ = utils.UnCamelCase(fieldName)
		sqlStr += fieldName + ","
		//如果查询属性的值为零值得话 不写进where查询里
		if !utils.IsZero(fieldValue) {
			whereStr += fieldName + "=" + utils.Parse(fieldValue) + " and "
		}

	}
	//trim掉逗号和and
	sqlStr = strings.TrimRight(sqlStr, ",") + " from " + tName + " where " + strings.TrimRight(whereStr, "and ")
	//trim掉空格
	sqlStr = strings.TrimSpace(sqlStr)
	//trim掉where
	sqlStr = strings.Trim(sqlStr, "where")

	fmt.Println("[sql-gorm-" + utils.DateFormat(time.Now(), "yyyy-MM-dd HH:mm:ss") + "]:" + sqlStr)

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


//将v2的值赋给v1
func setValue(v1 reflect.Value, v2 sql.RawBytes) {

	kind := v1.Kind()
	switch kind {
	case reflect.String:
		v1.Set(reflect.ValueOf(string(v2)))
	case reflect.Int:
		num, _ := strconv.Atoi(string(v2))
		v1.Set(reflect.ValueOf(num))
	}

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