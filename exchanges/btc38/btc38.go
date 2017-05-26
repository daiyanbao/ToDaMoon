package btc38

import (
	"ToDaMoon/exchanges"
	"sync"
)

var once sync.Once

//Run 会启动btc38模块
func Run() exchanges.Exchanger {
	once.Do(build)
	return btc38
}

func build() {
	//生成一个btc38的实例
	btc38 = instance()

	//展开btc38的各项任务
	btc38.checkNewCoin()
	btc38.watching()
}
