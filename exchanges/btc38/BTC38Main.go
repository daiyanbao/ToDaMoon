package btc38

import (
	"github.com/aQuaYi/ToDaMoon/exchanges"

	"sync"
	"time"
)

var btc38 *BTC38
var onceBTC38 sync.Once

// BTC38 包含了btc38.com的API所需的所有数据
type BTC38 struct {
	*config
	exchanges.TransRecordsDB
	exchanges.TradesCenter
	*exchanges.Account
}

// New 返回一个btc38的单例
func New() *BTC38 {
	onceBTC38.Do(buildBTC38)
	return btc38
}

// buildBTC38 构建了全局变量btc38
func buildBTC38() {
	// TODO: 还是感觉这几个步骤，太丑了。太丑了。

	//产生api的实例
	a := NewAPI()
	cfg := getConfig()
	// 产生btc38的实例
	btc38 = &BTC38{config: cfg}

	//配置tradesDB
	btc38.TransRecordsDB = exchanges.MakeTradesDBs(cfg.Database.Dir, cfg.Name, a.Markets)

	//获取btc38各个coin的全局交易的最新数据到数据库，然后，发布最新全局交易数据订阅
	btc38.TradesCenter = exchanges.MakeTradesCenter(a, btc38.TransRecordsDB, time.Millisecond*100, time.Minute)
}

//Name 返回BTC38的name
func (b *BTC38) Name() string {
	return b.config.Name
}
