package btc38

import (
	"ToDaMoon/Interface"
	"ToDaMoon/exchanges"
	"ToDaMoon/pubu"
	"fmt"
	"sync"
	"time"
)

var once sync.Once
var notify Interface.Notify

//Run 会启动btc38模块
func Run() exchanges.Exchanger {
	notify = pubu.New()
	once.Do(build)

	//以下是测试内容
	fmt.Println("=============================================================")
	b3Ticker, err := btc38.Ticker("btc", "cny")
	if err != nil {
		fmt.Println("BTC38.com BTC Ticker Error:", err)
	} else {
		fmt.Println("BTC38.com BTC Ticker", b3Ticker)
	}

	b3All, err := btc38.allTicker("cny")
	if err != nil {
		fmt.Println("无法获取btc38的cny市场的全部币的ticker")
	} else {
		fmt.Println("BTC38.com All Coins's Ticker:")
		for k, v := range b3All {
			fmt.Println(k, *v)
		}
	}

	b3Depth, err := btc38.Depth("btc", "cny")
	if err != nil {
		fmt.Println("无法获取btc38的cny市场的btc的depth")
	} else {
		fmt.Println("BTC38.com btc depth:")
		fmt.Println(b3Depth)

	}

	fmt.Println("=============================================================")
	b3Trades, err := btc38.Trades("btc", "cny", 0)
	if err != nil {
		fmt.Println("无法获取btc38的cny市场的btc的最新交易记录")
	} else {
		fmt.Println("BTC38.com btc 最新的交易记录:")
		fmt.Println(b3Trades)
	}

	fmt.Println("=============================================================")
	b3TradesSince1, err := btc38.Trades("btc", "cny", 1)
	if err != nil {
		fmt.Println("无法获取btc38的cny市场的btc的从1开始的交易记录")
	} else {
		fmt.Println("BTC38.com btc 从1开始的交易记录:")
		fmt.Println(b3TradesSince1)
	}

	fmt.Println("=============================================================")
	b3Balance, err := btc38.Balance()
	if err != nil {
		fmt.Println("无法获取btc38的账户信息")
	} else {
		fmt.Println("BTC38.com 的账户信息:")
		fmt.Println(b3Balance)
	}

	fmt.Println("=============================================================")
	for i := 20; i <= 40; i += 10 {
		orderID, err := btc38.Trade(BUY, "btc", "cny", 10000, float64(i)/10000)
		if err != nil {
			fmt.Println("无法在btc38.com下单买btc", err)
		} else {
			fmt.Println("BTC38.com下单买btc后的orderID是:")
			fmt.Println(i, orderID)
		}
		time.Sleep(time.Second)
	}

	fmt.Println("=============================================================")
	orderID, err := btc38.Trade(BUY, "btc", "cny", 10000, 90.0/10000)
	if err != nil {
		fmt.Println("无法在btc38.com下单买btc", err)
	} else {
		fmt.Println("BTC38.com下单买btc后的orderID是:")
		fmt.Println(orderID)
	}

	fmt.Println("=====等待撤单=====")
	for i := 10; i > 0; i-- {
		fmt.Println(i)
		time.Sleep(time.Second)
	}

	canceled, err := btc38.CancelOrder("btc", "cny", orderID)
	if err != nil {
		fmt.Println("撤销订单失败：", err)
	} else {
		fmt.Println("以下订单，已被撤销：", orderID, canceled)
	}

	fmt.Println("==============查看我的订单====================")
	myOrders, err := btc38.getMyOrders("btc", "cny")
	if err != nil {
		fmt.Println("无法获取我的订单", err)
	} else {
		for _, o := range myOrders {
			fmt.Println(o)
		}
	}

	//以上是测试内容

	return btc38
}

func build() {
	//生成一个btc38的实例
	btc38 = instance()

	//执行btc38的各项任务
	btc38.checkNewCoin()
	btc38.watching()
}
