package util

import (
	"testing"
	"time"
)

func Test_WaitFunc(t *testing.T) {
	t.Log("这个测试很花时间。。。")

	beginTime := time.Now()
	checkCycle := time.Millisecond * 100

	waitCh, wait := WaitFunc(checkCycle, "Test_WaitFunc")

	wait()
	t.Log(time.Now())
	waitTime1 := checkCycle
	if time.Since(beginTime) < waitTime1 || waitTime1+checkCycle < time.Since(beginTime) {
		t.Error("wait()在最小时间前结束了。")
	}

	updateCycle2 := time.Millisecond * 500
	waitCh <- updateCycle2
	wait()
	t.Log(time.Now())
	waitTime2 := waitTime1 + updateCycle2
	if time.Since(beginTime) < waitTime2 || waitTime2+checkCycle < time.Since(beginTime) {
		t.Error("wait()没能在updateCycle结束前，修改为更大的updateCycle")
	}

	updateCycle3 := time.Millisecond * 200
	go func() {
		time.Sleep(updateCycle3 / 2)
		waitCh <- updateCycle3
	}()
	wait()
	t.Log(time.Now())
	waitTime3 := waitTime2 + updateCycle3
	if time.Since(beginTime) < waitTime3 || waitTime3+checkCycle < time.Since(beginTime) {
		t.Error("wait()没能在updateCycle结束前，修改为更小的updateCycle")
	}
}
