package gorm

import (
	"reflect"
	"database/sql"
	"strconv"
)

//将v2的值赋给v1
func setValue(v1 reflect.Value, v2 sql.RawBytes) {

	kind := v1.Kind()
	switch kind {
	case reflect.String:
		v1.Set(reflect.ValueOf(string(v2)))
	case reflect.Int:
		num, _ := strconv.Atoi(string(v2))
		v1.Set(reflect.ValueOf(num))
	}

}
