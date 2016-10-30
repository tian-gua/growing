package gorm

import (
	"database/sql"
	"strings"
)

//根据数据库的table生成struct
//字段生成规则:
//id->Id
//user_name->UserName
func Generate(tableName string) (string, error) {

	//查询表结构信息
	sqlString := "desc " + tableName
	rows, err := gdb.Query(sqlString)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	//获得表所有字段
	columns, err := rows.Columns()
	if err != nil {
		return "", err
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
			return "", err
		}
		//存放每一条记录的第一个字段(表的字段名) 到fields切片里
		fields = append(fields, toCamelCase(string(values[0])))
		ftype := string(values[1])
		ftype = getDataType(ftype)
		fieldTypes = append(fieldTypes, ftype)
	}
	structString := "type " + strings.Title(tableName) + " struct{\n"
	for i, v := range fields {
		structString += "\t" + v + "\t" + fieldTypes[i] + "\t\t`field:\"" + unCamelCase(v) + "\"`" + "\n"
	}
	structString += "}"

	return structString, err
}
