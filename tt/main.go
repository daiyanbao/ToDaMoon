package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello World")
	s := Strs{
		"aaa",
		"bbb",
	}
	fmt.Println(s.Len())
}

type Strs []string

func (s Strs) Len() int {
	return len(s)
}
