package main

import (
	ec "ToDaMoon/exchanges"
	"ToDaMoon/exchanges/btc38"
	"fmt"
)

func main() {
	a := btc38.NewAPI()
	fmt.Printf("===开始===%s的API测试===\n", a.Name())

	price, result := ec.TestTicker(a, "cny", "btc")

	result += ec.TestDepth(a, "cny", "btc")

	result += ec.TestTransRecords(a, "cny", "btc", 1)

	result += ec.TestMyAccount(a)

	price = price * 0.8
	amount := 20 / price
	orderID, res := ec.TestOrder(a, ec.BUY, "cny", "btc", price, amount)
	result += res

	result += ec.TestMyOrders(a, "cny", "btc")

	result += ec.TestCancelOrder(a, "cny", "btc", orderID)

	result += ec.TestMyTransRecords(a, "cny", "doge", 1)

	result += btc38.TestAllTicker(a, "cny")

	result += btc38.TestMyTradeList(a, "cny", "doge", 0)

	fmt.Printf("===结束===%s的API测试===", a.Name())
	if result == "" {
		fmt.Print("\n全部通过\n")
	} else {
		fmt.Print("\n××失败××\n")
		fmt.Print(result)
	}
}
