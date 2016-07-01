package gutils

import "strings"


//结构体类型和数据库类型转换
func GetDBType(structType string) string {
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
