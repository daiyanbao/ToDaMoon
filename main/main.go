package main

import (
	ec "ToDaMoon/exchanges"
	"ToDaMoon/exchanges/btc38"
	"ToDaMoon/pubu"
	"ToDaMoon/util"
	"fmt"
	"log"
	"os"
	"time"
)

const (
	// pidFile 用来存储程序pid代号的文件
	pidFile = "tdm.pid"
)

func init() {
	util.Init(pidFile)
}

var (
	//Version 版本号
	Version string
	//BuildTime 编译时间
	BuildTime string
	//GitHash 当前的Git Hash码
	GitHash string
	//BuildNumber 编译次数
	BuildNumber string
)

func main() {
	//清理pid文件
	defer os.Remove(pidFile)

	if len(os.Args) > 1 && (os.Args[1] == "version" || os.Args[1] == "v") {
		fmt.Println("Version: ", Version+"."+BuildNumber)
		fmt.Println("Time:    ", BuildTime)
		fmt.Println("GitHash: ", GitHash)
		return
	}

	log.Println("Version: ", Version+"."+BuildNumber)

	done := util.WaitingKill()
	//以上是程序的相关准备工作
	pubuClient := pubu.New()
	b38 := btc38.Instance()
	fmt.Println(b38.Name)
	pubuClient.Good("ToDaMoon成功启动。")
	fmt.Println("ToDaMoon启动成功。")
	time.Sleep(time.Second)
	pubuClient.Good("我就看看能不能")
	fmt.Println(b38.Ticker("btc"))
	fmt.Println(b38.AllCoins())

	fmt.Println(b38.Balance())

	fmt.Println(b38.TradeHistory("btc", 1))
	b38btcStream := b38.Property["btc"].Observe()
	go func() {
		for {
			p := b38btcStream.WaitNext().(ec.Trades)
			fmt.Println("BTC\t\t", util.DateOf(p[len(p)-1].Date))
		}
	}()

	b38ltcStream := b38.Property["ltc"].Observe()
	go func() {
		for {
			p := b38ltcStream.WaitNext().(ec.Trades)
			fmt.Println("\tLTC\t", util.DateOf(p[len(p)-1].Date))
		}
	}()

	b38dogeStream := b38.Property["doge"].Observe()
	go func() {
		for {
			p := b38dogeStream.WaitNext().(ec.Trades)
			fmt.Println("\t\tdoge ", util.DateOf(p[len(p)-1].Date))
		}
	}()
	//等待被kill
	<-done
	pubuClient.Good("3秒后，ToDaMoon关闭。")
	time.Sleep(time.Second * 3)
}
