package exchanges

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/imkira/go-observer"
)

//TradesProperty 是exchange的trades监听属性
type TradesProperty map[string]map[string]observer.Property

func makeProperties(e Exchanger, tdbs TradesDBs, sleepTime time.Duration) TradesProperty {
	tp := make(TradesProperty)
	wg := &sync.WaitGroup{}
	for money, coinDBs := range tdbs {
		tp[money] = make(map[string]observer.Property)
		for coin, db := range coinDBs {
			wg.Add(1)
			tp[money][coin] = makePropertyAndSaveToDB(e, money, coin, db, sleepTime, wg)
		}
	}

	wg.Wait()
	text := fmt.Sprintln("已经创建了所有相关的监听属性")
	log.Println(text)

	return tp
}

func makePropertyAndSaveToDB(e Exchanger, money, coin string, db *TradesDB, sleepTime time.Duration, wg *sync.WaitGroup) observer.Property {
	th := Trades{}
	p := observer.NewProperty(th)
	go updatePropertyAndSaveToDB(e, money, coin, p, db, sleepTime)
	wg.Done()
	return p
}

func updatePropertyAndSaveToDB(e Exchanger, money, coin string, p observer.Property, db *TradesDB, sleepTime time.Duration) {
	maxTid, err := db.MaxTid()
	if err != nil {
		msg := fmt.Sprintf("updatePropertyAndSaveToDB(): 无法获取%s数据库的MaxTid: %s", db.Name(), err)
		log.Fatalln(msg)
	}

	for {
		th, err := e.Trades(money, coin, maxTid)
		if err != nil {
			msg := fmt.Sprintf("updatePropertyAndSaveToDB(): 获取%s交易所的%s市场的%s的历史交易数据失败：%s\n。。。。5秒后重试", e.Name, money, coin, err)
			log.Println(msg)
			time.Sleep(time.Second * 5)
			continue
		}

		if len(th) == 0 {
			time.Sleep(sleepTime)
			continue
		}

		p.Update(th)

		if err := db.Insert(th); err != nil {
			msg := fmt.Sprintf("插入%s交易所的%s市场的%s的历史交易数据失败：%s", e.Name, money, coin, err)
			log.Fatalln(msg)
		}

		time.Sleep(sleepTime)
	}
}
