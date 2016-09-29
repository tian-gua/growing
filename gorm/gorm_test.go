package gorm

import (
	_ "github.com/go-sql-driver/mysql"
	"testing"
	"time"
)

type practice struct {
	Id         int       `field:"id"`
	CreateTime time.Time `field:"create_time"`
	Name       string    `field:"name"`
}

func Test_gorm_query(t *testing.T) {
	InitDB("mysql", "root:root@tcp(127.0.0.1:3306)/practice")
	p := new([]practice)
	err := Query(&practice{Id: 3}, p)
	if err != nil {
		panic(err)
	}
	t.Log(p)

	p = new([]practice)
	err = QueryAll(p)
	if err != nil {
		panic(err)
	}
	t.Log(p)

	p = new([]practice)
	err = CustomQuery("select * from practice where id = 3", p)
	if err != nil {
		panic(err)
	}
	t.Log(p)
}

func Test_gorm_delete(t *testing.T) {
	InitDB("mysql", "root:root@tcp(127.0.0.1:3306)/practice")
	count, err := Delete(&practice{Id: 6})
	if err != nil {
		panic(err)
	}
	t.Log(count)
}

func Test_gorm_update(t *testing.T) {
	InitDB("mysql", "root:root@tcp(127.0.0.1:3306)/practice")
	id, err := Save(&practice{Id: 2, Name: "saf", CreateTime: time.Now()})
	if err != nil {
		panic(err)
	}
	t.Log(id)
}

func Test_gorm_save(t *testing.T) {
	InitDB("mysql", "root:root@tcp(127.0.0.1:3306)/practice")
	id, err := Save(&practice{Name: "aa", CreateTime: time.Now()})
	if err != nil {
		panic(err)
	}
	t.Log(id)
}
