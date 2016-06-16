package gorm

import (
	"reflect"
	"strconv"
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