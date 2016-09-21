package gorm

import (
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"testing"
)

type userAgent struct {
	Id   int `field:"id"`
	Name string `field:"name"`
	Age  int `field:"age"`
	Sex  int `field:"sex"`
}

func Test_ParseSql(t *testing.T) {
	fmt.Println(ParseDeleteByPrimaryKeySql(&userAgent{Id:1, Age:2}))
	fmt.Println(ParseQueryAllSql(&userAgent{}))
	fmt.Println(ParseQuerySql(&userAgent{Id:1, Name:"aaa", Age:22, Sex:1}))
	fmt.Println(ParseSaveSql(&userAgent{Id:1, Name:"aaa", Sex:1}))
	fmt.Println(ParseSaveSql(&userAgent{Name:"aaa", Age:22}))
	fmt.Println(ParseSaveSql(&userAgent{Id:1, Name:"阿斯蒂芬", Age:44, Sex:1}))
}