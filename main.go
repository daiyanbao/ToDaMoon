package main

import (
	"ToDaMoon/pubu"
	"fmt"
	"log"
	"os"

	"time"

	"github.com/go-ini/ini"
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
	cfg, err := ini.Load("./settings.ini")
	if err != nil {
		log.Fatalln("无法加载当前目录下的settings.ini文件。", err)
	}
	pbKey := cfg.Section("pubu").Key("ToDaMoon").String()
	pb := pubu.New(pbKey)
	imOK := pubu.Msg(time.Now().String())
	imOK.Success(
		"ToDaMoon成功启动",
		"详情请点击这里",
		"https://github.com/aQuaYi/ToDaMoon",
	)
	pb.Send(imOK)
}
