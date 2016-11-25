package gorm

import (
	"fmt"
	"strings"
)

//根据结构体生成查询sql
func ParseSelectSql(obj interface{}, conditions ...string) (string, error) {
	var fieldList string = ""
	//获得结构体反射的信息
	structInfo, err := GetStructInfo(obj)
	if err != nil {
		return "", err
	}
	tName := structInfo.TableName
	for _, structFiled := range structInfo.FieldsMap {
		//拼接字段集合
		fieldList += structFiled.tableFieldName + ","
	}
	//trim掉逗号
	sqlStr := fmt.Sprintf("SELECT %s FROM %s", strings.TrimRight(fieldList, ","), tName)

	where := ""
	order := ""
	for _, v := range conditions {
		if len(v) != 0 {
			if strings.Contains(v, "=") {
				where = where + " AND " + v
			} else {
				order = order + v + ","
			}
		}
	}

	if len(where) > 0 {
		sqlStr = sqlStr + " WHERE 1=1" + where
	}
	if len(order) > 0 {
		sqlStr = sqlStr + " ORDER BY " + order
	}

	return strings.TrimSuffix(sqlStr, ","), nil
}

//根据结构体生成删除sql
func ParseDeleteByPrimaryKeySql(obj interface{}) (string, error) {
	//用于存放sql字段
	var sqlStr string = ""
	//获得结构体反射的信息
	structInfo, err := GetStructInfo(obj)
	if err != nil {
		return "", err
	}
	tName := structInfo.TableName
	//获得要删除的id
	if structField, ok := structInfo.FieldsMap["id"]; ok && !isZero(structField.value) {
		//拼sql
		sqlStr = fmt.Sprintf("DELETE FROM %s WHERE id = %s", tName, structField.stringValue)
	} else {
		panic(fmt.Errorf("id not found or value is zero"))
	}
	return sqlStr, nil
}

//生成insert into语句
//isSelective属性为ture的话，零值的字段不会被拼到sql语句中
func ParseInsertSql(obj interface{}, isSelective bool) (string, error) {

	/**获得结构体反射的信息*/
	structInfo, err := GetStructInfo(obj)
	if err != nil {
		return "", err
	}
	tName := structInfo.TableName

	var valueList string = ""
	var fieldList string = ""
	/**拼sql*/
	for _, structField := range structInfo.FieldsMap {
		if "id" == structField.tableFieldName {
			continue
		} else if !isZero(structField.value) || !isSelective {
			fieldList += structField.tableFieldName + ","
			valueList += structField.stringValue + ","
		}
	}

	//去掉右边的逗号
	sqlStr := fmt.Sprintf("INSERT INTO %s(id,%s)VALUES(default,%s)", tName, strings.TrimRight(fieldList, ","), strings.TrimRight(valueList, ","))
	return sqlStr, nil
}

//生成update语句
//isSelective属性为ture的话，零值的字段不会被拼到sql语句中
func ParseUpdateSql(obj interface{}, isSelective bool) (string, error) {

	/**获得结构体反射的信息*/
	structInfo, err := GetStructInfo(obj)
	if err != nil {
		return "", err
	}
	tName := structInfo.TableName
	id := structInfo.FieldsMap["id"].stringValue

	var kvList string = ""
	//拼sql
	for _, structField := range structInfo.FieldsMap {
		if "id" == structField.tableFieldName {
			continue
		} else {
			//如果属性为零值则不更新
			if !isZero(structField.value) || !isSelective {
				kvList += structField.tableFieldName + "=" + structField.stringValue + ","
			}
		}
	}

	sqlStr := fmt.Sprintf("UPDATE %s SET %s WHERE id=%s", tName, strings.TrimRight(kvList, ","), id)
	return sqlStr, nil
}