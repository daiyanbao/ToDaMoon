package main

import (
	"fmt"

	ec "github.com/aQuaYi/exchanges"
	"github.com/aQuaYi/exchanges/btc38"
)

func main() {
	a := btc38.NewAPI()
	fmt.Printf("===开始===%s的API测试===\n", a.Name())

	price, result := ec.CheckTicker(a, "cny", "btc")

	result += ec.CheckDepth(a, "cny", "btc")

	result += ec.CheckTransRecords(a, "cny", "btc", 1)

	result += ec.CheckMyAccount(a)

	price = price * 0.8
	amount := 20 / price
	orderID, res := ec.CheckOrder(a, ec.BUY, "cny", "btc", price, amount)
	result += res

	result += ec.CheckMyOrders(a, "cny", "btc")

	result += ec.CheckCancelOrder(a, "cny", "btc", orderID)

	result += ec.CheckMyTransRecords(a, "cny", "doge", 1)

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
