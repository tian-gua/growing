package gorm

import (
	"strings"
	"fmt"
)


//根据结构体生成查询sql
func ParseQuerySql(obj interface{}) string {
	var fieldList string = ""
	var conditionList string = ""
	//获得结构体反射的信息
	structInfo := GetStructInfo(obj)
	tName := structInfo.TableName
	for _, structFiled := range structInfo.FieldsMap {
		//拼接字段集合
		fieldList += structFiled.tableFieldName + ","
		//如果查询属性的值为零值得话 不写进where查询里
		if !isZero(structFiled.value) {
			//拼接条件语句
			conditionList += structFiled.tableFieldName + "=" + structFiled.stringValue + " and "
		}
	}
	//拼接sql
	//trim掉逗号和and
	sqlStr := fmt.Sprint("select %s from %s where %s", strings.TrimRight(fieldList, ","), tName, strings.TrimRight(conditionList, "and "))
	//trim掉where 因为 conditionList可能为空
	sqlStr = strings.Trim(sqlStr, " where")
	return sqlStr
}


//根据结构体生成查询sql
func ParseQueryAllSql(obj interface{}) string {
	var fieldList string = ""
	//获得结构体反射的信息
	structInfo := GetStructInfo(obj)
	tName := structInfo.TableName
	for _, structFiled := range structInfo.FieldsMap {
		//拼接字段集合
		fieldList += structFiled.tableFieldName + ","
	}
	//trim掉逗号
	sqlStr := fmt.Sprint("select %s from %s", strings.TrimRight(fieldList, ","), tName)
	return sqlStr
}


//根据结构体生成删除sql
func ParseDeleteByPrimaryKeySql(obj interface{}) string {
	//获得结构体反射的信息
	structInfo := GetStructInfo(obj)
	tName := structInfo.TableName
	//获得要删除的id
	id := structInfo.FieldsMap["id"].stringValue
	if "0" == id {
		panic(fmt.Errorf("id not fount"))
	}
	//拼sql
	sqlStr := fmt.Sprintf("delete from %s where id = %s", tName, id)
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

