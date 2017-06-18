package main

import (
	ec "ToDaMoon/exchanges"
	"ToDaMoon/exchanges/btc38"
	"fmt"
)

func main() {
	a := btc38.NewAPI()
	fmt.Printf("===开始===%s的API测试===\n", a.Name())
	result := ec.TestAPI(a)

	fmt.Printf("==测试%s.AllTicker()==\n", a.Name())
	at, err := a.AllTicker("cny")
	if err != nil {
		msg := fmt.Sprintf("%s.AllTicker Error:%s\n", a.Name(), err)
		result += msg
		fmt.Print(msg)
	} else {
		fmt.Printf("%s.AllTicker of cny and btc\n%s\n", a.Name(), at)
	}

	fmt.Printf("===结束===%s的API测试===", a.Name())
	if result == "" {
		fmt.Print("全部通过\n")
	} else {
		fmt.Print("××失败××\n")
		fmt.Print(result)
	}
}

/*
//testAPI 会逐个测试网站的API
func testAPI(btc38 *btc38.API) {

	fmt.Println("==测试btc38.Ticker()==")
	b3Ticker, err := btc38.Ticker("cny", "btc")
	if err != nil {
		fmt.Println("BTC38.com BTC Ticker Error:", err)
	} else {
		fmt.Println("BTC38.com BTC Ticker", b3Ticker)
	}

	fmt.Println("==测试btc38.allTicker()==")
	b3All, err := btc38.AllTicker("cny")
	if err != nil {
		fmt.Println("无法获取btc38的cny市场的全部币的ticker")
	} else {
		fmt.Println("BTC38.com All Coins's Ticker:")
		for k, v := range b3All {
			fmt.Println(k, *v)
		}
	}

	b3Depth, err := btc38.Depth("cny", "btc")
	if err != nil {
		fmt.Println("无法获取btc38的cny市场的btc的depth")
	} else {
		fmt.Println("BTC38.com btc depth:")
		fmt.Println(b3Depth)

	}

	fmt.Println("=============================================================")
	b3Trades, err := btc38.Trades("cny", "btc", 0)
	if err != nil {
		fmt.Println("无法获取btc38的cny市场的btc的最新交易记录")
	} else {
		fmt.Println("BTC38.com btc 最新的交易记录:")
		fmt.Println(b3Trades)
	}

	fmt.Println("=============================================================")
	b3TradesSince1, err := btc38.Trades("cny", "btc", 1)
	if err != nil {
		fmt.Println("无法获取btc38的cny市场的btc的从1开始的交易记录")
	} else {
		fmt.Println("BTC38.com btc 从1开始的交易记录:")
		for _, t := range b3TradesSince1 {
			fmt.Println(*t)
		}
	}

	fmt.Println("=============================================================")
	b3Balance, err := btc38.Account()
	if err != nil {
		fmt.Println("无法获取btc38的账户信息")
	} else {
		fmt.Println("BTC38.com 的账户信息:")
		fmt.Println(b3Balance)
	}

	fmt.Println("=============================================================")
	// for i := 20; i <= 40; i += 10 {
	// 	orderID, err := btc38.Trade(BUY, "cny", "btc", 10000, float64(i)/10000)
	// 	if err != nil {
	// 		fmt.Println("无法在btc38.com下单买btc", err)
	// 	} else {
	// 		fmt.Println("BTC38.com下单买btc后的orderID是:")
	// 		fmt.Println(i, orderID)
	// 	}
	// 	time.Sleep(time.Second)
	// }

	fmt.Println("=============================================================")
	orderID, err := btc38.Trade(btc38.BUY, "cny", "btc", 10000, 90.0/10000)
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

	canceled, err := btc38.CancelOrder("cny", "btc", orderID)
	if err != nil {
		fmt.Println("撤销订单失败：", err)
	} else {
		fmt.Println("以下订单，已被撤销：", orderID, canceled)
	}

	fmt.Println("==============查看我的订单====================")
	myOrders, err := btc38.getMyOrders("cny", "btc")
	if err != nil {
		fmt.Println("无法获取我的订单", err)
	} else {
		for _, o := range myOrders {
			fmt.Println(o)
		}
	}

	fmt.Println("=======查看我的交易记录==========")
	myDogeTrades, err := btc38.getMyTrades("cny", "doge", 0)
	if err != nil {
		fmt.Println("无法获取的doge交易记录。", err)
	} else {
		for _, t := range myDogeTrades {
			fmt.Println(t)
		}
	}
}
*/
