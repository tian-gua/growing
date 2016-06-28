package gorm

import (
	"reflect"
	"strings"
	"errors"
	"utils"
	"fmt"
	"time"
)


//根据结构体生成查询sql
func parseQuerySql(obj interface{}) string {

	var sqlStr string = ""
	var whereStr string = ""
	//获得类型的信息
	ov := reflect.ValueOf(obj)
	v := reflect.Indirect(ov)
	t := v.Type()

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

	return sqlStr

}




//根据结构体生成查询sql
func parseQueryAllSql(value reflect.Value) string {

	var sqlStr string = ""
	//获得反射信息
	v := reflect.Indirect(value)
	t := v.Type()

	tName := t.Name()
	//结构体首字母转换为小写
	//结构体首字母大写是为了供其他包访问,数据库则不用
	tName = strings.ToLower(tName)

	sqlStr = "select "

	//拼接需要查询的字段
	for i := 0; i < t.NumField(); i++ {
		fieldName := t.Field(i).Name
		//字符串反驼峰转换,例如 UserName 会变成 user_name
		fieldName, _ = utils.UnCamelCase(fieldName)
		sqlStr += fieldName + ","
	}
	//trim掉逗号和and
	sqlStr = strings.TrimRight(sqlStr, ",") + " from " + tName

	fmt.Println("[sql-gorm-" + utils.DateFormat(time.Now(), "yyyy-MM-dd HH:mm:ss") + "]:" + sqlStr)

	return sqlStr

}

//根据结构体生成插入或者更新sql
func parseSaveSql(obj interface{}) string {

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

	return sqlStr

}

//根据结构体生成删除sql
func parseDeleteSql(obj interface{}) string {

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

	fmt.Println("[sql-gorm-" + utils.DateFormat(time.Now(), "yyyy-MM-dd HH:mm:ss") + "]:" + sqlStr)

	return sqlStr
}
