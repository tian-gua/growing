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
	sqlStr := fmt.Sprintf("select %s from %s where %s", strings.TrimRight(fieldList, ","), tName, strings.TrimRight(conditionList, "and "))
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
	sqlStr := fmt.Sprintf("select %s from %s", strings.TrimRight(fieldList, ","), tName)
	return sqlStr
}


//根据结构体生成删除sql
func ParseDeleteByPrimaryKeySql(obj interface{}) string {
	//用于存放sql字段
	var sqlStr string = ""
	//获得结构体反射的信息
	structInfo := GetStructInfo(obj)
	tName := structInfo.TableName
	//获得要删除的id
	if structField, ok := structInfo.FieldsMap["id"]; ok && !isZero(structField.value){
		//拼sql
		sqlStr = fmt.Sprintf("delete from %s where id = %s", tName, structField.stringValue)
	} else {
		panic(fmt.Errorf("id not found or value is zero"))
	}
	return sqlStr
}


//根据结构体生成插入或者更新sql
func ParseSaveSql(obj interface{}) string {
	//用于判断是否为保存还是更新
	var isInsert bool = false
	//用于存放sql字段
	var sqlStr string = ""
	//获得类型的信息
	structInfo := GetStructInfo(obj)
	tName := structInfo.TableName
	//取id得值判断是insert 还是 update
	id := structInfo.FieldsMap["id"].stringValue
	if "0" == id {
		isInsert = true
	}
	if isInsert {
		var valueList string = ""
		var fieldList string = ""
		//拼sql
		for _, structField := range structInfo.FieldsMap {
			fieldList += structField.tableFieldName + ","
			if "id" == structField.tableFieldName {
				valueList += "default,"
			} else {
				valueList += structField.stringValue + ","
			}
		}
		//去掉右边的逗号
		sqlStr = fmt.Sprintf("insert into %s(%s)values(%s)", tName, strings.TrimRight(fieldList, ","), strings.TrimRight(valueList, ","))
	} else {
		var kvList string = ""
		//拼sql
		for _, structField := range structInfo.FieldsMap {
			if "id" == structField.tableFieldName {
				continue
			} else {
				//如果属性为零值则不更新
				if !isZero(structField.value) {
					kvList += structField.tableFieldName + "=" + structField.stringValue + ","
				}
			}
		}
		sqlStr = fmt.Sprintf("update %s set %s where id = %s", tName, strings.TrimRight(kvList, ","), id)
	}
	return sqlStr
}

