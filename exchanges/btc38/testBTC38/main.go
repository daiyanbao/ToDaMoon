package main

import (
	"fmt"

	"github.com/aQuaYi/ToDaMoon/exchanges"
	"github.com/aQuaYi/ToDaMoon/exchanges/btc38"
)

func main() {
	a := btc38.NewAPI()
	fmt.Printf("===开始===%s的API测试===\n", a.Name())

	price, result := exchanges.CheckTicker(a, "cny", "btc")

	result += exchanges.CheckDepth(a, "cny", "btc")

	result += exchanges.CheckTransRecords(a, "cny", "btc", 1)

	result += exchanges.CheckMyAccount(a)

	price = price * 0.8
	amount := 20 / price
	orderID, res := exchanges.CheckOrder(a, exchanges.BUY, "cny", "btc", price, amount)
	result += res

	result += exchanges.CheckMyOrders(a, "cny", "btc")

	result += exchanges.CheckCancelOrder(a, "cny", "btc", orderID)

	result += exchanges.CheckMyTransRecords(a, "cny", "doge", 1)

	result += btc38.CheckAllTicker(a, "cny")

	result += btc38.CheckMyTradeList(a, "cny", "doge", 0)

	fmt.Printf("===结束===%s的API测试===", a.Name())
	if result == "" {
		fmt.Print("\n全部通过\n")
	} else {
		fmt.Print("\n××失败××\n")
		fmt.Print(result)
	}
}
