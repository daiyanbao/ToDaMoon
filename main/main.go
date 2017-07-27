package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aQuaYi/GoKit"
	"github.com/aQuaYi/ToDaMoon/apollo"
)

const (
	// pidFile 用来存储程序pid代号的文件
	pidFile = "tdm.pid"
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

func init() {
	//获取并保存pid
	GoKit.Init(pidFile)
}

func main() {
	//清理pid文件
	defer os.Remove(pidFile)

	ver := Version + "." + BuildNumber
	if len(os.Args) > 1 && (os.Args[1] == "version" || os.Args[1] == "v") {
		fmt.Println("Version: ", ver)
		fmt.Println("Time:    ", BuildTime)
		fmt.Println("GitHash: ", GitHash)
		return
	}

	log.Println("======= LAUNCH, Version ", ver, "=======")
	defer log.Println("======= LANDED, Version ", ver, "=======")

	apollo.Launch()
}
