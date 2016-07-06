package gorm

import (
	"database/sql"
)

//事务操作
//添加互斥锁 , 锁注释掉了,发现tx里面自带了锁
type Transaction struct {
	//m  *sync.Mutex
	tx *sql.Tx
}

//提交
func (this *Transaction) Commit() {
	this.tx.Commit()
	//this.m.Unlock()
}

//回滚
func (this *Transaction) RollBack() {
	this.tx.Rollback()
	//this.m.Unlock()
}

//开启事务
func Begin() (*Transaction, error) {
	transaction := &Transaction{}
	//transaction.m.Lock()
	tx, err := gdb.Begin()
	if err != nil {
		//如果开启事务出错,解锁返回错误
		//transaction.m.Unlock()
		return nil, err
	}
	transaction.tx = tx
	return transaction, nil
}


