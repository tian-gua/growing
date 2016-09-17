package gorm

import (
	"reflect"
	"strconv"
)


//将接收的 值反射 转换成字符串类型
func parseValueToDBString(v reflect.Value) string {
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
func isZero(v reflect.Value) bool {
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

//获得空切片元素的类型
func getEmptySliceValue(slice reflect.Value) reflect.Value {
	t := slice.Type()
	//给切片元素开辟一个空间
	vSlice := reflect.MakeSlice(t, 1, 1)
	//获得 切片元素 的反射信息
	element := vSlice.Slice(0, 1).Index(0)
	return element
}