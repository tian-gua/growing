package gorm

import (
	"strings"
	"fmt"
)


//根据结构体生成查询sql
func parseQuerySql(obj interface{}) string {
	var sqlStr string = ""
	var whereStr string = ""
	//获得类型的信息
	structInfo := GetStructInfo(obj)
	tName := structInfo.TableName
	sqlStr = "select "
	for _, v := range structInfo.FieldsMap {
		sqlStr += v.tableFieldName + ","
		//如果查询属性的值为零值得话 不写进where查询里
		if !isZero(v.value) {
			whereStr += v.tableFieldName + "=" + v.stringValue + " and "
		}
	}
	//trim掉逗号和and
	sqlStr = strings.TrimRight(sqlStr, ",") + " from " + tName + " where " + strings.TrimRight(whereStr, "and ")
	//trim掉空格
	sqlStr = strings.TrimSpace(sqlStr)
	//trim掉where
	sqlStr = strings.Trim(sqlStr, "where")
	return sqlStr
}




//根据结构体生成查询sql
func parseQueryAllSql(obj interface{}) string {
	var sqlStr string = ""
	//获得反射信息
	structInfo := GetStructInfo(obj)
	tName := structInfo.TableName
	sqlStr = "select "
	for _, v := range structInfo.FieldsMap {
		sqlStr += v.tableFieldName + ","
	}
	//trim掉逗号和and
	sqlStr = strings.TrimRight(sqlStr, ",") + " from " + tName
	return sqlStr
}

//根据结构体生成插入或者更新sql
func parseSaveSql(obj interface{}) string {
	//用于判断是否为插入方法
	var isInsert bool = false
	//用于存放sql字段
	var sqlStr string = ""
	//获得类型的信息
	structInfo := GetStructInfo(obj)
	tName := structInfo.TableName
	//取id得值判断是insert 还是 update
	id := structInfo.FieldsMap["Id"].stringValue
	if "0" == id {
		isInsert = true
	}
	if isInsert {
		//拼sql
		sqlStr = "insert into " + tName + "("
		var valueStr string
		//拼sql
		for _, v := range structInfo.FieldsMap {
			sqlStr += v.tableFieldName + ","
			if "Id" == v.name {
				valueStr += "default,"
			} else {
				valueStr += v.stringValue + ","
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
		//拼sql
		for _, v := range structInfo.FieldsMap {
			if "Id" == v.name {
				continue
			} else {
				//如果属性为零值则不更新
				if !isZero(v.value) {
					sqlStr += v.tableFieldName + "=" + v.stringValue + ","
				}
			}
		}
		sqlStr = strings.TrimRight(sqlStr, ",") + " where id = " + id
	}
	return sqlStr
}

//根据结构体生成删除sql
func parseDeleteSql(obj interface{}) string {
	//用于存放sql字段
	var sqlStr string = ""
	//获得类型的信息
	structInfo := GetStructInfo(obj)
	tName := structInfo.TableName
	//获得要删除的id
	id := structInfo.FieldsMap["Id"].stringValue
	if "0" == id {
		panic(fmt.Errorf("id not fount"))
	}
	//拼sql
	sqlStr = "delete from " + tName + " where id = " + id
	return sqlStr
}
