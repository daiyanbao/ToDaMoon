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

func makePropertys(e Exchanger, tdbs TradesDBs) TradesProperty {
	wg := &sync.WaitGroup{}
	for money, coinDBs := range tdbs {
		for coin, db := range coinDBs {
			wg.Add(1)
			go listeningTradeHistoryAndSave(e, money, coin, db, wg)
		}
	}
	wg.Wait()
	text := fmt.Sprintln("已经创建了所有相关的监听属性")
	log.Println(text)
	if o.ShowDetail {
		go notify.Debug(text)
	}
}

func updateDB(e Exchanger, money, coin string, db *TradesDB) {
	maxTid, err := db.MaxTid()
	if err != nil {
		msg := fmt.Sprintf("updataDB(): 无法获取%s数据库的MaxTid: %s", db.Name(), err)
		log.Fatalln(msg)
	}

	unUpdated := true
	for unUpdated {
		th, err := e.Trades(money, coin, maxTid)
		if err != nil {
			msg := fmt.Sprintf("获取%s交易所的%s市场的%s的历史交易数据失败：%s\n。。。。5秒后重试", e.Name, money, coin, err)
			log.Println(msg)
			time.Sleep(time.Second * 5)
			continue
		}

		if len(th) == 0 { //已经获取到最新的Trades数据了。
			return
		}

		if th[len(th)-1].Date+60 > time.Now().Unix() { //已经更新到一分钟以内的数据了。认为已经更新到最新的数据了。
			unUpdated = false
		}

		maxTid = th[len(th)-1].Tid
		if err := db.Insert(th); err != nil {
			msg := fmt.Sprintf("插入%s交易所的%s市场的%s的历史交易数据失败：%s", e.Name, money, coin, err)
			log.Fatalln(msg)
		}
	}
}

func updateDBAndProperty(e Exchanger, money, coin string, db *TradesDB) observer.Property {
	maxTid, err := db.MaxTid()
	if err != nil {
		msg := fmt.Sprintf("updataDBAndProperty(): 无法获取%s数据库的MaxTid: %s", db.Name(), err)
		log.Fatalln(msg)
	}

	//NOTICE: 本来准备让p发布的所有参数都是非零的。
	//但是，发现，就算获取了最新的数据，也有可能是长度为0的，我就放弃了。
	p := observer.NewProperty(Trades{})

	var thdb Trades
	saveTime := time.Now()

	go func() {
		for {
			th, err := e.Trades(money, coin, maxTid)
			if err != nil {
				msg := fmt.Sprintf("获取%s交易所的%s市场的%s的历史交易数据失败：%s\n。。。。5秒后重试", e.Name, money, coin, err)
				log.Println(msg)
				time.Sleep(time.Second * 5)
				continue
			}

			if len(th) > 0 {
				p.Update(th)
				//FIXME: 此处应该是通过channel发送到另外一个goroutin，让那个goroutin来管理数据什么时候存入数据库。
				//那个goroutine不订阅这个property的原因是因为，我也不知道
				thdb = append(thdb, th...)
			}
		}
	}()
	return p
}

// func listeningTradeHistoryAndSave(e Exchanger, money, coin string, db *TradesDB, wg *sync.WaitGroup) {
// 	defer wg.Done()

// 	maxTid, err := db.MaxTid()
// 	if err != nil {

// 	}
// 	if o.ShowDetail {
// 		text := fmt.Sprintf("%s的%s的MaxTid是%d", o.Name, coin, maxTid)
// 		log.Println(text)
// 	}
// 	th, err := o.TradeHistory(coin, maxTid)
// 	if err != nil {
// 		log.Fatalf("无法获取%s的%s的历史交易数据。\n", o.Name, coin)
// 	}
// 	o.Property[coin] = observer.NewProperty(th)
// 	if o.ShowDetail {
// 		text := fmt.Sprintf("%s的%s: 已经创建了监听属性。", o.Name, coin)
// 		log.Println(text)
// 	}

// 	var thdb ec.Trades
// 	saveTime := time.Now()
// 	requestTime := time.Now()
// 	waitS := o.CoinPeriodS

// 	go func() {
// 		for {
// 			if th.Len() > 0 {
// 				maxTid = th[th.Len()-1].Tid
// 			}
// 			th, err = o.TradeHistory(coin, maxTid)
// 			if err != nil {
// 				text := fmt.Sprintf("请求%s的%s的历史交易数据失败, 5秒后重试。\n%s", o.Name, coin, err)
// 				notify.Error(text)
// 				log.Println(text)
// 				time.Sleep(time.Second * 5)
// 				continue
// 			}

// 			if th.Len() > 0 {
// 				o.Property[coin].Update(th)
// 				thdb = append(thdb, th...)
// 			}
// 			if thdb.Len() > 0 {
// 				if thdb.Len() > 10*10000 || time.Since(saveTime) > time.Minute*30 {
// 					if err := o.db[coin].Insert(thdb); err != nil {
// 						text := fmt.Sprintf("往%s的%s的数据库插入数据出错:%s\n", o.Name, coin, err)
// 						notify.Error(text)
// 						log.Fatalln(text)
// 					}
// 					date := thdb[thdb.Len()-1].Date
// 					text := fmt.Sprintf("%s的**%s数据库**的最新日期为%s", o.Name, coin, util.DateOf(date))
// 					notify.Info(text)
// 					saveTime = time.Now()
// 					thdb = ec.Trades{}
// 				}
// 			}

// 			if th.Len() < 50 { // 当th的长度较短时，是由于已经读取到最新的消息了。
// 				waitS = 300
// 			} else {
// 				waitS = o.CoinPeriodS
// 			}
// 			util.HoldOn(time.Duration(waitS)*time.Millisecond, &requestTime)
// 		}
// 	}()
// }
