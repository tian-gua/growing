package gorm

import (
	"reflect"
	"database/sql"
	"strconv"
	"time"
	"fmt"
)

//将v2的值赋给v1
func setRawData(v reflect.Value, rawData sql.RawBytes) {
	switch t := v.Interface().(type){
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
