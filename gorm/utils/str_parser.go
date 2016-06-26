package utils

import (
	"reflect"
	"strconv"
	"regexp"
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
func ToCamelCase(str string) (string, error) {
	//判断参数是否为空
	if len(str) == 0 {
		return "", nil
	}
	//查找所有的_x 字符串,并替换成X
	reg, err := regexp.Compile("_([a-z])");
	if err != nil {
		return str, err
	}
	for {
		//找到匹配的字符串
		findStr := reg.FindString(str)
		if len(findStr) == 0 {
			break
		}
		//去掉匹配到的_x字符串中的_, 并将x转换成大写
		newStr := strings.ToUpper(strings.Trim(findStr, "_"))
		str = strings.Replace(str, findStr, strings.Title(newStr), 1)

	}
	return str, nil
}

//和ToCamelCase方法襄樊
//将ACd 转换成 A_cd
func UnCamelCase(str string) (string, error) {
	//判断参数是否为空
	if len(str) == 0 {
		return "", nil
	}
	//查找所有的大写字符
	//req, err := regexp.Compile("[A-Z]")
	//if err != nil {
	//	return str, err
	//}
	return "", nil

}