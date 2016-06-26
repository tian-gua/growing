package utils

import (
	"reflect"
	"strconv"
	"regexp"
	"fmt"
	"utils"
	"strings"
)

//将接受的反射的值 转为 字符换类型
func Parse(v reflect.Value) string {

	var result string
	//根据值得类型转换字符串
	switch v.Kind() {
	case reflect.String:
		result = "'" + v.String() + "'"
	case reflect.Int:
		result = strconv.FormatInt(v.Int(), 10)
	}
	return result

}


//判断是否是零值
func IsZero(v reflect.Value) bool {
	//获得值得类型
	kind := v.Kind()
	switch kind {
	case reflect.String:
		if "" == v.String() {
			return true
		}
		return false
	case reflect.Int:
		if 0 == v.Int() {
			return true
		}
		return false
	}
	return false
}


//将字符串转为驼峰命名规则,并且首字母大写
//如: a_cd -> ACd
func ToCamelCase(str string) string {
	//判断参数是否为空
	if utils.IsEmpty(&str) {
		return ""
	}
	//获得第一个字符串并转换成大写
	bytes := []byte(str)
	firstWord := string(bytes[0])
	newFirstWord := strings.ToUpper(firstWord)
	str = strings.Replace(str, firstWord, newFirstWord, 1)
	//*str = string(bytes)
	//查找所有的_x 字符串,并替换成X
	reg, err := regexp.Compile("_([a-z])");
	if err != nil {
		fmt.Println(err)
	}
	//去掉匹配到的_x字符串中的_, 并将x转换成大写
	newStr := strings.ToUpper(strings.Trim(reg.FindString(str), "_"))
	//将x替换成X
	return reg.ReplaceAllString(str, newStr)
}