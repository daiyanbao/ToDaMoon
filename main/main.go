package main

import (
	"ToDaMoon/pubu"
	"ToDaMoon/util"
	"fmt"
	"log"
	"os"

	"time"

	"github.com/go-ini/ini"
)

const (
	// pidFile 用来存储程序pid代号的文件
	pidFile = "tdm.pid"
)

func init() {
	util.Init(pidFile)
}

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
	//清理pid文件
	defer os.Remove(pidFile)

	if len(os.Args) > 1 && (os.Args[1] == "version" || os.Args[1] == "v") {
		fmt.Println("Version: ", Version+"."+BuildNumber)
		fmt.Println("Time:    ", BuildTime)
		fmt.Println("GitHash: ", GitHash)
		return
	}

	done := util.WaitingKill()
	//以上是程序的相关准备工作

	cfg, err := ini.Load("./settings.ini")
	if err != nil {
		log.Fatalln("无法加载当前目录下的settings.ini文件。", err)
	}
	pbKey := cfg.Section("pubu").Key("ToDaMoon").String()
	pbc := pubu.New(pbKey)
	pbc.Good("ToDaMoon成功启动。")

	//等待被kill
	<-done
	pbc.Good("10秒后，ToDaMoon关闭。")
	time.Sleep(time.Second * 10)
}
