package util

import (
	"time"
)

//HoldOn 使程序暂停一个duration。
//通常运用于API访问限制
func HoldOn(duration time.Duration, beginTime time.Time) time.Time {
	//无需判断，如果Sleep的时间为负，是不会Sleep的。
	time.Sleep(duration - time.Since(beginTime))
	return time.Now()
}

//DateOf 返回一个unix tiemstamp的格式化。
func DateOf(t int64) string {
	return time.Unix(t, 0).String()
}
