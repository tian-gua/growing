package gorm

import (
	"strings"
	"regexp"
	"fmt"
)


//结构体类型和数据库类型转换
func getDataType(structType string) string {
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

//将字符串转为驼峰命名规则,并且首字母大写
//如: a_cd 转换成 ACd
func toCamelCase(str string) string {
	//判断参数是否为空
	if len(str) == 0 {
		return ""
	}
	//查找所有的_x 字符串,并替换成X
	reg, err := regexp.Compile("_([a-z])");
	if err != nil {
		fmt.Println(err)
		return str
	}
	for {
		//找到匹配的字符串
		findStr := reg.FindString(str)
		//匹配结束则退出循环
		if len(findStr) == 0 {
			break
		}
		//去掉匹配到的_x字符串中的_, 并将x转换成大写
		newStr := strings.ToUpper(strings.Trim(findStr, "_"))
		str = strings.Replace(str, findStr, newStr, -1)
	}
	return strings.Title(str)
}

//和ToCamelCase方法襄樊
//将ACd 转换成 a_cd
func unCamelCase(str string) string {
	//判断参数是否为空
	if len(str) == 0 {
		return ""
	}
	//查找所有的大写字符
	reg, err := regexp.Compile("[A-Z]")
	if err != nil {
		fmt.Println(err)
		return str
	}
	for {
		//找到匹配的字符串
		findStr := reg.FindString(str)
		//匹配结束则退出循环
		if len(findStr) == 0 {
			break
		}
		//将大写字符串变成 _ 加 对应小写
		newStr := "_" + strings.ToLower(findStr)
		str = strings.Replace(str, findStr, newStr, -1)
	}
	//如果转换完成自后 第一个 字符串 为 _ ,则删掉_
	return strings.TrimLeft(str, "_")
}

