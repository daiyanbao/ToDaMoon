package apollo

import (
	"github.com/aQuaYi/ToDaMoon/exchanges/btc38"
	"github.com/aQuaYi/ToDaMoon/util"
)

const ()

//Launch 是阿波罗飞船的启动程序
func Launch() {

	//done等待程序结束的信号
	done := util.WaitingKill()

	//在这里启动各个交易所模块
	//TODO: 更换为controller()
	btc38.Start()

	//等待被kill
	<-done
}
