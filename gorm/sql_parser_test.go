package gorm

import (
	_ "github.com/go-sql-driver/mysql"
	"testing"
)

type userAgent struct {
	Id   int    `field:"id"`
	Name string `field:"name"`
	Age  int    `field:"age"`
	Sex  int    `field:"sex"`
}

func (u *userAgent) GetTableName() string {
	return "tb_user_agent"
}

func Test_ParseSql(t *testing.T) {
	t.Log(ParseDeleteByPrimaryKeySql(&userAgent{Id: 1, Age: 2}))
	t.Log(ParseQueryAllSql(&userAgent{}))
	t.Log(ParseQuerySql(&userAgent{Id: 1, Name: "aaa", Age: 22, Sex: 1}))
	t.Log(ParseSaveSql(&userAgent{Id: 1, Name: "aaa", Sex: 1}))
	t.Log(ParseSaveSql(&userAgent{Name: "aaa", Age: 22}))
	t.Log(ParseSaveSql(&userAgent{Id: 1, Name: "阿斯蒂芬", Age: 44, Sex: 1}))
}