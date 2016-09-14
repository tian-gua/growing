package gorm

import (
	"strings"
	"regexp"
	"fmt"
)


//结构体类型和数据库类型转换
func getDBType(structType string) string {
	dbtype := "string"
	switch true {
	case strings.HasPrefix(structType, "varchar"):
		dbtype = "string"
	case strings.HasPrefix(structType, "int"):
		dbtype = "int"
	case strings.HasPrefix(structType, "decimal"):
		dbtype = "int64"
	case strings.HasPrefix(structType, "timestamp") || strings.HasPrefix(structType, "datetime"):
		dbtype = "time.Time"
	}
	return dbtype
}

//根据结构体名获取表名
//例如：getTableName("userAgent") = user_agent
func getTableName(structName string) string {
	//判断参数是否为空
	if len(structName) == 0 {
		return ""
	}
	//查找所有的大写字符
	reg, err := regexp.Compile("[A-Z]")
	if err != nil {
		fmt.Println(err)
		return structName
	}
	for {
		//找到匹配的字符串
		findStr := reg.FindString(structName)
		//匹配结束则退出循环
		if len(findStr) == 0 {
			break
		}
		//将大写字符串变成 _ 加 对应小写
		newStr := "_" + strings.ToLower(findStr)
		structName = strings.Replace(structName, findStr, newStr, -1)
	}
	//如果转换完成自后 第一个 字符串 为 _ ,则删掉_
	return strings.TrimLeft(structName, "_")
}