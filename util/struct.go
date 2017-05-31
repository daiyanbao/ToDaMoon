package util

import "fmt"

//IsTypeEqual 判断两个interface{}类型的变量的原始类型是否相等
//TODO: 找到一种更好的判断方法
func IsTypeEqual(a, b interface{}) bool {
	aType := fmt.Sprintf("%T", a)
	bType := fmt.Sprintf("%T", b)
	if aType == bType {
		return true
	}
	return false
}
