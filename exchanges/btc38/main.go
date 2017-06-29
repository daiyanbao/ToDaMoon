package btc38

import (
	"sync"

	"github.com/aQuaYi/ToDaMoon/Interface"
	"github.com/aQuaYi/ToDaMoon/exchanges"
	"github.com/aQuaYi/ToDaMoon/pubu"
)

var once sync.Once
var notify Interface.Notifier

//Run 会启动btc38模块
func Run() exchanges.API {
	once.Do(build)

	return btc38
}

func build() {
	notify = pubu.New()
	//生成一个btc38的实例
	Instance()

	//执行btc38的各项任务
	btc38.checkNewCoin()
	btc38.watching()
}
