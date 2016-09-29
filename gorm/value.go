package gorm

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

//将v2的值赋给v1
func setRawData(v reflect.Value, rawData sql.RawBytes) {
	switch t := v.Interface().(type) {
	case string:
		v.Set(reflect.ValueOf(string(rawData)))
	case int:
		num, _ := strconv.Atoi(string(rawData))
		v.Set(reflect.ValueOf(num))
	case time.Time:
		newTime, err := time.ParseInLocation("2006-01-02 15:04:05", string(rawData), time.Local)
		if err != nil {
			panic(err)
		}
		v.Set(reflect.ValueOf(newTime))
	default:
		fmt.Println("未处理的类型:%v", t)
	}
}

//将接收的 值反射 转换成字符串类型
func parseValueToDBString(v reflect.Value) string {
	//如果v是指针，则取指向的值
	v = reflect.Indirect(v)
	//根据值得类型转换字符串
	switch t := v.Interface().(type) {
	case string:
		return "'" + v.String() + "'"
	case int:
		return strconv.FormatInt(v.Int(), 10)
	case int64:
		return strconv.FormatInt(v.Int(), 10)
	case time.Time:
		return "'" + t.Format("2006-01-02 15:04:05") + "'"
	}
	return ""
}

//判断是否是零值
func isZero(v reflect.Value) bool {
	//如果v是指针，则取指向的值
	v = reflect.Indirect(v)
	switch t := v.Interface().(type) {
	case string:
		return "" == v.String()
	case int:
		return 0 == int(t)
	case int64:
		return 0 == int64(t)
	case time.Time:
		return 0 == t.Second()
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
