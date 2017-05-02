package util

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
)

/*
Init 初始化的方法汇总，提供了以下功能:
1.利用CPU全部内核
2.保存当前程序pid号到pidFile文件中。关闭程序， 可在命令行使用
$ kill `cat $pidFile`
*/
func Init(pidFile string) {
	//利用全部CPU内核
	runtime.GOMAXPROCS(runtime.NumCPU())

	//保存当前程序的pid值
	if pid := syscall.Getpid(); pid != 1 {
		ioutil.WriteFile(pidFile, []byte(strconv.Itoa(pid)), 0777)
		fmt.Println("PID is", pid)
	}
}

// WaitingKill 会返回一个channel，当程序接收到kill信号时，会关闭此通道。
func WaitingKill() chan struct{} {
	closeSignal := make(chan os.Signal, 1)
	done := make(chan struct{})
	signal.Notify(closeSignal, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		log.Println("接收到关闭信号。", <-closeSignal)
		close(done)
	}()
	return done
}
