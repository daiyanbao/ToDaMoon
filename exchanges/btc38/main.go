package btc38

import (
	"sync"

	"github.com/aQuaYi/ToDaMoon/apollo"
	"github.com/aQuaYi/pubu.im"
)

var onceMain sync.Once
var notify apollo.Notifier

// Start 会启动btc38模块
// 会执行tasks中的所有任务
func Start() {
	onceMain.Do(tasks)
}

func tasks() {
	// 使用pubu.im作为通知工具
	notify = pubu.New()

	// // 核心任务
	// 生成一个btc38的实例

	New()

	// 执行btc38的各项任务
	btc38.checkNewCoin()
	btc38.watching()
}
