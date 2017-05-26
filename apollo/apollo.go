package apollo

import (
	"ToDaMoon/exchanges/btc38"
	"ToDaMoon/util"
	"os"
	"sync"
)

const (
	// pidFile 用来存储程序pid代号的文件
	pidFile = "tdm.pid"
)

var once sync.Once

func init() {
	util.Init(pidFile)
}

//Launch 是阿波罗飞船的启动程序
func Launch() {
	//清理pid文件
	defer os.Remove(pidFile)
	done := util.WaitingKill()

	btc38.Run()

	//等待被kill
	<-done
}
