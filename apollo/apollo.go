package apollo

import (
	"ToDaMoon/exchanges/btc38"
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

	//done等待程序结束的信号
	done := util.WaitingKill()

	//在这里启动各个交易所模块
	btc38.Run()

	//等待被kill
	<-done
}
