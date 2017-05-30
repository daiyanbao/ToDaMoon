package util

import (
	"os"
)

//Exist 判断文件是否存在
func Exist(filename string) bool {
	if _, err := os.Stat(filename); err != nil {
		return false
	}
	return true
}
