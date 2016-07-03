package gorm

import "errors"

var (
	//获得结构体信息错误
	GET_STRUCTINFO_ERROR = errors.New("get struct info error")
	//没有找到id
	NOT_FOUND_ID = errors.New("not found id")
)