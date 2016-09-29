package gutils

//判断字符串指针的值是否为空
func IsEmpty(str *string) bool {
	if str == nil || len(*str) == 0 {
		return true
	}
	return false
}

//判断字符串指针的值是否不为空
func IsNotEmpty(str *string) bool {
	if str == nil || len(*str) == 0 {
		return false
	}
	return true
}
