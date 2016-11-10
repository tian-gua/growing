package gorm

import (
	"testing"
	"errors"
)

func TestTsHook(t *testing.T) {
	InitDB("mysql", "root:root@tcp(127.0.0.1:3306)/practice")
	t.Log("测试回滚")
	TsHook(hook)
	t.Log("测试提交")
	TsHook(hook2)
}

func hook(ts *Transaction) error {
	return errors.New("aaa")
}

func hook2(ts *Transaction) error {
	return nil
}
