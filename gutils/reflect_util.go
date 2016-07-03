package gutils

import (
	"reflect"
	"strconv"
)

//判断是否是零值
func IsZero(v reflect.Value) bool {
	//获得值得类型
	kind := reflect.Indirect(v).Kind()
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
	return true
}


//将接受的反射的值 转为 字符换类型
func ParseValueToString(v reflect.Value) string {
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

