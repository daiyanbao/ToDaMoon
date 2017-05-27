package main

import (
	"ToDaMoon/apollo"
	"fmt"
	"log"
	"os"
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

}
