package util

import "reflect"

//IsTypeEqual 判断两个interface{}类型的变量的原始类型是否相等
func IsTypeEqual(a, b interface{}) bool {
	at := reflect.TypeOf(a).Kind()
	bt := reflect.TypeOf(b).Kind()

	if at == bt {
		return true
	}
	return false
}
