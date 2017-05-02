package util

import (
	"time"
)

//HoldOn 使程序暂停一个duration。
//通常运用于API访问限制
func HoldOn(duration time.Duration, beginTime *time.Time) {
	//无需判断，如果Sleep的时间为负，是不会Sleep的。
	time.Sleep(duration - time.Since(*beginTime))
	*beginTime = time.Now()
}
