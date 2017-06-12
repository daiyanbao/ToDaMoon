package exchanges

import (
	"ToDaMoon/util"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/imkira/go-observer"
)

//TradeSubject 会发布最新的交易数据。
//会通过observer.Property来更新最新的交易数据。
//updateCycleCh 可以修改Property的更新周期
type TradeSubject struct {
	observer.Property
	UpdateCycleCh chan<- time.Duration
}

//TradesSubject 是exchnge的trades监听属性
type TradesSubject map[string]map[string]TradeSubject

//MakeSubjectes 返回了所尊coin的最新订阅消息
func MakeSubjectes(e Exchanger, tdbs TradesDBs, checkCycle, updateCycle time.Duration) TradesSubject {
	log.Printf("开始创建%s的监听属性", e.Name())
	tp := make(TradesSubject)
	wg := &sync.WaitGroup{}
	for money, coinDBs := range tdbs {
		tp[money] = make(map[string]TradeSubject)
		for coin, db := range coinDBs {
			wg.Add(1)
			tp[money][coin] = makePropertyAndSaveToDB(e, money, coin, db, checkCycle, updateCycle, wg)
		}
	}

	wg.Wait()
	text := fmt.Sprintf("已经创建了%s所有相关的监听属性", e.Name())
	log.Println(text)

	return tp
}

func makePropertyAndSaveToDB(e Exchanger, money, coin string, db *TradesDB, checkCycle, updateCycle time.Duration, wg *sync.WaitGroup) TradeSubject {
	th := Trades{}
	p := observer.NewProperty(th)
	ch := updatePropertyAndSaveToDB(e, money, coin, p, db, checkCycle, updateCycle)
	log.Printf("已经创建了%s-%s-%s的监听属性", e.Name(), money, coin)

	wg.Done()
	return TradeSubject{
		Property:      p,
		UpdateCycleCh: ch,
	}
}

func updatePropertyAndSaveToDB(e Exchanger, money, coin string, p observer.Property, db *TradesDB, checkCycle, updateCycle time.Duration) chan<- time.Duration {
	maxTid, err := db.MaxTid()
	if err != nil {
		msg := fmt.Sprintf("updatePropertyAndSaveToDB(): 无法获取%s数据库的MaxTid: %s", db.Name(), err)
		log.Fatalln(msg)
	}
	name := fmt.Sprintf("%s的%s的%s的Trade的Property", e.Name(), money, coin)
	waitCh, wait := util.WaitFunc(checkCycle, name)

	go func() {
		for {
			th, err := e.Trades(money, coin, maxTid)
			if err != nil {
				msg := fmt.Sprintf("updatePropertyAndSaveToDB(): 获取%s交易所的%s市场的%s的历史交易数据失败：%s\n。。。。5秒后重试", e.Name(), money, coin, err)
				log.Println(msg)
				time.Sleep(time.Second * 5)
				continue
			}

			if len(th) > 0 {

				p.Update(th)

				if err := db.Insert(th); err != nil {
					msg := fmt.Sprintf("插入%s交易所的%s市场的%s的历史交易数据失败：%s", e.Name(), money, coin, err)
					log.Fatalln(msg)
				}

				maxTid = th[len(th)-1].Tid
			}
			wait()
		}
	}()

	waitCh <- updateCycle

	return waitCh
}

//ChangeUpdateCycleTo 修改了Property的更新周期
func (t TradesSubject) ChangeUpdateCycleTo(duration time.Duration) {
	for _, coins := range t {
		for _, ts := range coins {
			ts.UpdateCycleCh <- duration
		}
	}
}
