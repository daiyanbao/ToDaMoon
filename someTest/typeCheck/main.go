package main

import (
	mit "ToDaMoon/someTest/typeCheck/myIntType"
	"fmt"
)

func main() {
	var m mit.MyInt = 1
	fmt.Printf("m's type is %T\n", m)
}
