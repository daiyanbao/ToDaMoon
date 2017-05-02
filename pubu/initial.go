package pubu

import (
	"ToDaMoon/Interface"
	"ToDaMoon/util"
	"sync"
	"time"
)

var pbc *client
var once sync.Once

//New 返回一个单例的*client客户端
func New(hook string) Interface.Notify {
	once.Do(
		func() {
			icChan := make(incomingChan, 12)
			pbc = &client{
				hook:   hook,
				icChan: icChan,
			}
			go start(icChan)
			pbc.Info("零信的初始化工作完成。")
		})

	return pbc
}

func start(icChan incomingChan) {
	beginTime := time.Now()
	for m := range icChan {
		pbc.send(m)
		//由于pubu.im有API访问次数限制，“每个接入调用限制为每秒 10 次”。所以，需要暂停一下。
		util.HoldOn(time.Millisecond*100, &beginTime)
	}
}
