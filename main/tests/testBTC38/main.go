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

	fmt.Printf("\n==测试%s.AllTicker()==\n", a.Name())
	at, err := a.AllTicker("cny")
	if err != nil {
		msg := fmt.Sprintf("%s.AllTicker Error:%s\n", a.Name(), err)
		result += msg
		fmt.Print(msg)
	} else {
		fmt.Printf("%s.AllTicker of cny and btc\n%s\n", a.Name(), at)
	}

	fmt.Printf("\n==测试%s.MyTradeList()==\n", a.Name())
	mtl, err := a.MyTradeList("cny", "doge", 1)
	if err != nil {
		msg := fmt.Sprintf("%s.MyTradeList Error:%s\n", a.Name(), err)
		result += msg
		fmt.Print(msg)
	} else {
		fmt.Printf("%s.MyTradeList(cny, doge, 1)=\n%v\n", a.Name(), mtl)
	}

	fmt.Printf("===结束===%s的API测试===", a.Name())
	if result == "" {
		fmt.Print("全部通过\n")
	} else {
		fmt.Print("××失败××\n")
		fmt.Print(result)
	}
}
