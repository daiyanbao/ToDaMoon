package main

import (
	ec "ToDaMoon/exchanges"
	"ToDaMoon/exchanges/btc38"
	"fmt"
)

func main() {
	a := btc38.NewAPI()
	fmt.Printf("===开始===%s的API测试===\n", a.Name())
	//FIXME:
	result := ec.TestAPI(a)

	fmt.Printf("===结束===%s的API测试===", a.Name())
	if result == "" {
		fmt.Print("全部通过\n")
	} else {
		fmt.Print("××失败××\n")
		fmt.Print(result)
	}
}
