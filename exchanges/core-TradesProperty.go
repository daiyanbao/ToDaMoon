package exchanges

import (
	"ToDaMoon/util"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/imkira/go-observer"
)

//TradePublish 会发布最新的交易数据。
//会通过observer.Property来更新最新的交易数据。
//updateCycleCh 可以修改Property的更新周期
type TradePublish struct {
	observer.Property
	UpdateCycleCh chan<- time.Duration
}

//TradesCenter 包含了exchnge的所有trades监听属性
type TradesCenter map[string]map[string]TradePublish

//MakeTradesCenter 返回了所有coin的最新订阅属性
func MakeTradesCenter(a API, tdbs TransRecordsDB, checkCycle, updateCycle time.Duration) TradesCenter {
	log.Printf("开始创建%s的监听属性", a.Name())
	tp := make(TradesCenter)
	wg := &sync.WaitGroup{}
	for money, coinDBs := range tdbs {
		tp[money] = make(map[string]TradePublish)
		for coin, db := range coinDBs {
			wg.Add(1)
			tp[money][coin] = makePropertyAndSaveToDB(a, money, coin, db, checkCycle, wg)
			tp[money][coin].UpdateCycleCh <- updateCycle
		}
	}

	wg.Wait()
	text := fmt.Sprintf("已经创建了%s所有相关的监听属性", a.Name())
	log.Println(text)

	return tp
}

func makePropertyAndSaveToDB(a API, money, coin string, db *TradesDB, checkCycle time.Duration, wg *sync.WaitGroup) TradePublish {
	th := Trades{}
	p := observer.NewProperty(th)
	ch := updatePropertyAndSaveToDB(a, money, coin, p, db, checkCycle)
	log.Printf("已经创建了%s-%s-%s的监听属性", a.Name(), money, coin)

	wg.Done()
	return TradePublish{
		Property:      p,
		UpdateCycleCh: ch,
	}
}

func updatePropertyAndSaveToDB(a API, money, coin string, p observer.Property, db *TradesDB, checkCycle time.Duration) chan<- time.Duration {
	maxTid, err := db.MaxTid()
	if err != nil {
		msg := fmt.Sprintf("updatePropertyAndSaveToDB(): 无法获取%s数据库的MaxTid: %s", db.Name(), err)
		log.Fatalln(msg)
	}
	name := fmt.Sprintf("%s的%s的%s的Trade的Property", a.Name(), money, coin)
	waitCh, wait := util.WaitFunc(checkCycle, name)

	go func() {
		for {
			th, err := a.TransRecords(money, coin, maxTid)
			if err != nil {
				msg := fmt.Sprintf("updatePropertyAndSaveToDB(): 获取%s交易所的%s市场的%s的历史交易数据失败：%s\n。。。。5秒后重试", a.Name(), money, coin, err)
				log.Println(msg)
				time.Sleep(time.Second * 5)
				continue
			}

			if len(th) > 0 {

				p.Update(th)

				if err := db.Insert(th); err != nil {
					msg := fmt.Sprintf("插入%s交易所的%s市场的%s的历史交易数据失败：%s", a.Name(), money, coin, err)
					log.Fatalln(msg)
				}

				maxTid = th[len(th)-1].Tid
			}
			wait()
		}
	}()

	return waitCh
}

//ChangeUpdateCycleTo 修改了Property的更新周期
func (t TradesCenter) ChangeUpdateCycleTo(duration time.Duration) {
	for _, coins := range t {
		for _, ts := range coins {
			ts.UpdateCycleCh <- duration
		}
	}
}
