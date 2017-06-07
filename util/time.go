package util

import (
	"log"
	"time"
)

//DateOf 返回一个unix tiemstamp的格式化。
func DateOf(t int64) string {
	return time.Unix(t, 0).String()
}

//WaitFunc 返回一个可以使用通道修改休眠时间的等待函数
//updaetCycle的修改会立即生效，不用等到此updateCycle结束。
//例如：当100秒的updateCycle时间已过51秒时，把updateCycle修改为50秒，程序会立刻结束。不用等到100秒才结束。
//例如：100秒的updateCycle在结束前，把updateCycle修改为200秒，程序会立即生效。
//checkCycle是检查是否到期的时间段，也是最小等待时间段。
//没有修改updateCycle的时候，程序的等待时间是checkCycle×int(updateCycle/checkCycle+1)
func WaitFunc(checkCycle time.Duration, name string) (chan<- time.Duration, func()) {
	cycleCh := make(chan time.Duration, 3)
	beginTime := time.Now()
	updateCycle := checkCycle

	return cycleCh, func() {
		for beginTime.Add(updateCycle).After(time.Now()) {
			select {
			case updateCycle = <-cycleCh:
				if updateCycle <= checkCycle {
					log.Println("WARNING: updateCycle<=checkCycle，程序会按照checkCycle来等待。")
				}
				log.Printf("%s的wait的updateCycle已经修改为%s", name, updateCycle)
			default:
			}
			time.Sleep(checkCycle)
		}
		beginTime = time.Now()
	}
}

//SleepFunc 返回一个等待sleep函数， 使程序暂停一个duration。
//SleepFunc是以上WaitFunc的简化版本，通常运用于API访问限制
//利用闭包，把beginTime变量包裹到了sleep函数内。
func SleepFunc(duration time.Duration) func() {
	beginTime := time.Now()
	return func() {
		//无需判断，如果Sleep的时间为负，则不会Sleep。
		time.Sleep(duration - time.Since(beginTime))
		beginTime = time.Now()
	}
}
