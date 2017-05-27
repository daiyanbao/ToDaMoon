package apollo

import (
	"ToDaMoon/exchanges/btc38"
	"ToDaMoon/util"
	"fmt"
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

	//在这里启动各个交易所模块
	b3 := btc38.Run()
	b3Ticker, err := b3.Ticker("btc", "cny")
	if err != nil {
		fmt.Println("BTC38.com BTC Ticker Error:", err)
	} else {
		fmt.Println("BTC38.com BTC Ticker", b3Ticker)
	}
	//等待被kill
	<-done
}
