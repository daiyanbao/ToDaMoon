package btc38

import (
	"fmt"
	"log"

	"time"

	"sync"

	"github.com/aQuaYi/GoKit"
	"github.com/aQuaYi/ToDaMoon/exchanges"
	observer "github.com/imkira/go-observer"
)

// watching 观察最新的TransOrders数据
func (b *BTC38) watching() {
	if !b.IsLog {
		return
	}
	notify.Info("BTC38 的 watching 开始工作了～～～")

	// REVIEW: 观察能否修改时间周期。
	done := make(chan struct{}, 3)

	sum := 1
	for money, coins := range b.TradesCenter {
		c := 0
		for coin, property := range coins {
			tsStream := property.Observe()
			logTransOrders(b.Name(), money, coin, tsStream, property.UpdateCycleCh, done)
			c++
		}
		sum += c
	}

	time.Sleep(5 * time.Second)
	log.Println("现在来修改获取时间")
	changeUpdateCycle(b, sum, done)

}

func logTransOrders(name, money, coin string, stream observer.Stream, updateCycleCH chan<- time.Duration, done chan<- struct{}) {
	var once sync.Once
	format := fmt.Sprintf("%s-%s-%s的最新数据是，", name, money, coin) + "TID:%d, Date:%s"

	go func() {
		for {
			ts, ok := stream.WaitNext().(exchanges.Trades)
			if ok && ts.Len() > 0 {
				log.Printf(format, ts.LastTID(), GoKit.DateOf(ts.LastDate()))

			}

			// 为了其他的coin尽快收集
			if ts.Len() < 30 {
				updateCycleCH <- 5 * time.Minute
				once.Do(func() {
					done <- struct{}{}
					msg := "又收集到一个最新数据\n"
					msg += fmt.Sprintf(format, ts.LastTID(), GoKit.DateOf(ts.LastDate()))
					notify.Good(msg)
				})
			}
		}
	}()
}

func changeUpdateCycle(b *BTC38, sum int, done <-chan struct{}) {
	go func() {
		for {
			for _, coins := range b.TradesCenter {
				for _, property := range coins {
					property.UpdateCycleCh <- time.Duration(sum*b.Database.MinUpdateCycleMS) * time.Millisecond
				}
			}
			<-done
			sum--
			if sum == 1 {
				notify.Good("BTC38所有的Coin，都已经收集完全了。")
			}
		}
	}()
}
