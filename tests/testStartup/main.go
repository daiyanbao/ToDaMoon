package main

/*
本程序是为了配合开机启动ToDaMoon的脚本而编写的程序。

程序会产生一个filename命名文件，包含以下内容
程序执行的时间
程序的pwd
执行程序的用户
用户的username
用户的home
用户的uid
用户的gid
*/
import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"time"
)

var filename = "testStartup.txt"

func main() {
	content := time.Now().String()
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	content += "\n" + wd + "\n"

	usr, err := user.Current()
	if err != nil {
		fmt.Println(err)
	}

	content += fmt.Sprintln("Name     :", usr.Name)
	content += fmt.Sprintln("username :", usr.Username)
	content += fmt.Sprintln("home     :", usr.HomeDir)
	content += fmt.Sprintln("UID      :", usr.Uid)
	content += fmt.Sprintln("GID      :", usr.Gid)

	err = ioutil.WriteFile(filename, []byte(content), 0667)
	if err != nil {
		log.Fatal(err)
	}
}
