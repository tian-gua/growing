package gorm

import (
	"database/sql"
	"fmt"
)

//事务操作
//添加互斥锁 , 锁注释掉了,发现tx里面自带了锁
type Transaction struct {
	tx *sql.Tx
}

//nil, 不适用事务时,可以用到此变量
var nilTs *Transaction

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

//需要保证事务执行的函数的类型
type TsFunc func(*Transaction) error

//执行用户需要保证事务的函数
func TsHook(tsf TsFunc) error {

	/**获取tx*/
	tx, err := Begin()
	if err != nil {
		return err
	}

	/**如果没有发生异常则提交,否则回滚*/
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("事务过程中发生异常:[%s],即将回滚...\n", err)
			tx.RollBack()
			fmt.Println("事务回滚成功!")
		} else {
			fmt.Println("事务过程执行完毕,即将提交...")
			tx.Commit()
			fmt.Println("事务提交成功!")
		}
	}()

	/**执行hook,如果发生错误就panic,交给defer处理*/
	err = tsf(tx)
	if err != nil {
		panic(err)
	}
	return nil
}

//获得statement,有事务和非事务2种情况
func getStatement(sqlStr string, gtx *Transaction) (*sql.Stmt, error) {
	//校验是否初始化
	if !isInit {
		panic("no db init")
	}

	var stmt *sql.Stmt
	var err error
	//判断是否在事务中执行
	if gtx != nil {
		stmt, err = gtx.tx.Prepare(sqlStr)
		if err != nil {
			return nil, err
		}
	} else {
		//从sql.DB里获得stmt
		stmt, err = gdb.Prepare(sqlStr)
		if err != nil {
			return nil, err
		}
	}
	return stmt, err
}