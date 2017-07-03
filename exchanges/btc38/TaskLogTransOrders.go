package btc38

import (
	"fmt"
	"log"

	ec "github.com/aQuaYi/ToDaMoon/exchanges"
	"github.com/aQuaYi/ToDaMoon/util"
	observer "github.com/imkira/go-observer"
)

// watching 观察最新的TransOrders数据
func (b *BTC38) watching() {
	if !b.ShowDetail {
		return
	}

	for money, coins := range b.TradesCenter {
		for coin, property := range coins {
			tsStream := property.Observe()
			logTransOrders(b.Name(), money, coin, tsStream)
		}
	}

}

func logTransOrders(name, money, coin string, stream observer.Stream) {
	format := fmt.Sprintf("%s-%s-%s的最新数据是，", name, money, coin) + "TID:%d, Date:%s"
	go func() {
		for {
			ts, ok := stream.WaitNext().(ec.Trades)
			if ok && ts.Len() > 0 {
				log.Printf(format, ts.LastTID(), util.DateOf(ts.LastDate()))
			}
		}
	}()
}
