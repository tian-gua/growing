package gorm

import (
	"database/sql"
)

//事务操作
//添加互斥锁 , 锁注释掉了,发现tx里面自带了锁
type Transaction struct {
	tx *sql.Tx
}

//提交
func (this *Transaction) Commit() {
	err := this.tx.Commit()
	if err != nil {
		panic(err)
	}
}

//回滚
func (this *Transaction) RollBack() {
	err := this.tx.Rollback()
	if err != nil {
		panic(err)
	}
}

//开启事务
func Begin() (*Transaction, error) {
	transaction := &Transaction{}
	tx, err := gdb.Begin()
	if err != nil {
		return nil, err
	}
	transaction.tx = tx
	return transaction, nil
}


