package gorm

import (
	"reflect"
	"database/sql"
	"strconv"
	"time"
	"fmt"
)

//将v2的值赋给v1
func setValue(v1 reflect.Value, v2 sql.RawBytes) {
	switch t := v1.Interface().(type){
	case string:
		v1.Set(reflect.ValueOf(string(v2)))
	case int:
		num, _ := strconv.Atoi(string(v2))
		v1.Set(reflect.ValueOf(num))
	case time.Time:
		newTime, err := time.ParseInLocation("2006-01-02 15:04:05", string(v2), time.Local)
		if err != nil {
			panic(err)
		}
		v1.Set(reflect.ValueOf(newTime))
	default:
		fmt.Println("未处理的类型:%v", t)
	}
}
