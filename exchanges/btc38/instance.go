package btc38

import (
	"github.com/aQuaYi/ToDaMoon/exchanges"

	"time"
)

var btc38 *BTC38

//BTC38 包含了btc38.com的API所需的所有数据
type BTC38 struct {
	*API
	exchanges.TransRecordsDB
	exchanges.TradesCenter
	*exchanges.Account
}

//Instance 返回一个btc38的单例
func Instance() *BTC38 {
	//读取配置文件
	cfg := getConfig()

	//生成btc38实例
	btc38 = genBTC38By(cfg)

	//获取btc38各个coin的全局交易的最新数据到数据库，然后，发布最新全局交易数据订阅
	btc38.TradesCenter = exchanges.MakeTradesCenter(btc38.API, btc38.TransRecordsDB, time.Millisecond*100, time.Minute)

	return btc38
}

func genBTC38By(c *config) *BTC38 {
	a := NewAPI()
	tdb := exchanges.MakeTradesDBs(c.DBDir, c.Name, c.Markets)

	btc38 = &BTC38{API: a,
		TransRecordsDB: tdb,
	}

	return btc38
}

//Name 返回BTC38的name
func (b *BTC38) Name() string {
	return b.API.config.Name
}
