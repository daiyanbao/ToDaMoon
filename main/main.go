package main

import (
	ec "ToDaMoon/exchanges"
	"ToDaMoon/exchanges/btc38"

	"ToDaMoon/apollo"
	"ToDaMoon/pubu"
	"ToDaMoon/util"
	"fmt"
	"log"
	"os"
	"time"
)

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

	if len(os.Args) > 1 && (os.Args[1] == "version" || os.Args[1] == "v") {
		fmt.Println("Version: ", Version+"."+BuildNumber)
		fmt.Println("Time:    ", BuildTime)
		fmt.Println("GitHash: ", GitHash)
		return
	}

	log.Println("Launch, Version ", Version+"."+BuildNumber)
	defer log.Println("Landing, Version ", Version+"."+BuildNumber)

	apollo.Launch()

	//TODO: 清除以下程序

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

	//fmt.Println(b38.Balance())

	//fmt.Println(b38.TradeHistory("btc", 1))
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
}
