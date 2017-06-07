package pubu

import (
	"ToDaMoon/Interface"
	"ToDaMoon/util"
	"log"
	"sync"
	"time"

	"github.com/go-ini/ini"
)

var pbc *client
var once sync.Once

//New 返回一个单例的*client客户端
func New() Interface.Notifier {
	once.Do(
		func() {
			cfg, err := ini.Load("./pubu.ini")
			if err != nil {
				log.Fatalf("无法加载%s/pubu.ini: %s", util.PWD(), err)
			}
			hook := cfg.Section("pubu").Key("hook").String()

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
	sleep := util.SleepFunc(time.Millisecond * 100)
	for m := range icChan {
		pbc.send(m)

		//由于pubu.im有API访问次数限制，“每个接入调用限制为每秒 10 次”。所以，需要暂停一下。
		sleep()
	}
}
