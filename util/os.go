package util

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

//PWD like shell cmd pwd
//PWD 返回当前目录路径
func PWD() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	index := strings.LastIndex(path, string(os.PathSeparator))
	ret := path[:index]
	return ret
}
