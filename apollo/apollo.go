package apollo

import (
	"ToDaMoon/util"
	"os"
)

const (
	// pidFile 用来存储程序pid代号的文件
	pidFile = "tdm.pid"
)

func init() {
	util.Init(pidFile)
}

//Launch 是阿波罗飞船的启动程序
func Launch() {
	//清理pid文件
	defer os.Remove(pidFile)
	done := util.WaitingKill()

	//等待被kill
	<-done
}
