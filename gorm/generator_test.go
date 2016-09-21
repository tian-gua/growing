package gorm

import (
	_ "github.com/go-sql-driver/mysql"
	"testing"
)

func Test_Generate(t *testing.T) {
	InitDB("mysql", "root:root@tcp(127.0.0.1:3306)/practice")
	str, err := Generate("practice")
	if err != nil {
		t.Log(err)
	}
	t.Log("\n" + str)
}
