package apollo

import (
	"ToDaMoon/exchanges/btc38"
	"ToDaMoon/util"
)

const ()

//Launch 是阿波罗飞船的启动程序
func Launch() {

	//done等待程序结束的信号
	done := util.WaitingKill()

	//在这里启动各个交易所模块
	//TODO: 更换为control
	btc38.Run()

	//等待被kill
	<-done
}
