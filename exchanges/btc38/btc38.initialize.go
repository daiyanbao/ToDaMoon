package btc38

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"
	"sync"
	"time"

	db "ToDaMoon/DataBase"
	"ToDaMoon/Interface"
	ec "ToDaMoon/exchanges"
	"ToDaMoon/pubu"
	"ToDaMoon/util"

	"github.com/imkira/go-observer"
)

var notify Interface.Notify

func init() {
	notify = pubu.New()
}

// 以下代码实现了BTC38模块的单例特性
var btc38 *BTC38
var once sync.Once

// Instance make a singleton of btc38.com
func Instance(cfg *Config, notify Interface.Notify) *BTC38 {
	cfg.Check()
	once.Do(func() {
		askChan := make(askChannel, 12)
		btc38 = &BTC38{Config: cfg,
			ask:      askChan,
			db:       map[string]db.DBM{},
			Property: map[string]observer.Property{},
		}
		go start(askChan, time.Duration(cfg.MinAccessPeriodMS)*time.Millisecond)

		btc38.makeDBs()
		btc38.makePropertys()
		notify.Info("单例初始化完成。")
	})

	return btc38
}

func start(askChan askChannel, waitTime time.Duration) {
	beginTime := time.Now()
	for ask := range askChan {
		switch ask.Type {
		case get:
			data, err := ec.Get(ask.Path)
			ask.AnswerChan <- answer{body: data, err: err}
		case post:
			data, err := ec.Post(ask.Path, ask.Headers, ask.Body)
			ask.AnswerChan <- answer{body: data, err: err}
		default:
			log.Println("Wrong ask type.")
		}
		util.HoldOn(waitTime, &beginTime)
	}
}

func (o *BTC38) makeDBs() {
	for _, coin := range o.Coins {
		var err error
		o.db[coin], err = db.New(o.DBDir, o.Name, coin, "cny")
		if err != nil {
			text := fmt.Sprintf("无法创建OKCoin的中%s的数据库。\n", coin)
			notify.Error(text)
			log.Fatalln(text)
		}

		if o.ShowDetail {
			maxTid, _ := o.db[coin].MaxTid()
			text := fmt.Sprintf("已经链接上了%s的数据库，其最大Tid是%d\n", coin, maxTid)
			notify.Info(text)
		}
	}

	if o.ShowDetail {
		notify.Debug("已经创建了相关的数据库")
	}
}

func (o *BTC38) makePropertys() {
	for _, coin := range o.Coins {
		if o.ShowDetail {
			text := fmt.Sprintf("%s: 要开始创建监听属性了。", coin)
			notify.Debug(text)
		}
		listeningTradeHistoryAndSave(o, coin)
		if o.ShowDetail {
			text := fmt.Sprintf("%s: 已经创建了相关的监听属性。", coin)
			notify.Debug(text)
		}
	}

	if o.ShowDetail {
		text := fmt.Sprintln("已经创建了所有相关的监听属性")
		notify.Debug(text)
	}
}

func listeningTradeHistoryAndSave(o *BTC38, coin string) {
	maxTid, err := o.db[coin].MaxTid()
	if err != nil {
		text := fmt.Sprintf("%s: 没有读取到相关的数据库的最大值", coin)
		notify.Error(text)
		log.Fatalln(text)
	}
	if o.ShowDetail {
		text := fmt.Sprintf("OKCoin的%s的MaxTid是%d\n", coin, maxTid)
		notify.Debug(text)
	}
	th, err := o.TradeHistory(coin, maxTid)
	if err != nil {
		return
	}
	o.Property[coin] = observer.NewProperty(th)
	if o.ShowDetail {
		text := fmt.Sprintf("%s: 已经创建了监听属性。", coin)
		notify.Debug(text)
	}
	var thdb ec.Trades
	saveTime := time.Now()
	requestTime := time.Now()
	waitMS := o.CoinPeriodS

	go func() {
		for {
			if th.Len() > 0 {
				maxTid = th[th.Len()-1].Tid
			}
			th, err = o.TradeHistory(coin, maxTid)
			if err != nil {
				text := fmt.Sprintf("请求OKCoin.cn的%s的历史交易数据失败\n%s", coin, err)
				notify.Error(text)
				log.Println(text)
				time.Sleep(time.Second * 2)
				continue
			}

			if th.Len() > 0 {
				o.Property[coin].Update(th)
				thdb = append(thdb, th...)
			}
			if thdb.Len() > 0 {
				if thdb.Len() > 20*10000 || time.Since(saveTime) > time.Hour {
					if err := o.db[coin].Insert(thdb); err != nil {
						text := fmt.Sprintf("往%s的%s的数据库插入数据出错:%s\n", o.Name, coin, err)
						notify.Error(text)
						log.Fatalln(text)
					}
					date := thdb[thdb.Len()-1].Date
					text := fmt.Sprintf("%s的**%s数据库**的最新日期为%s", o.Name, coin, util.DateOf(date))
					notify.Info(text)
					saveTime = time.Now()
					thdb = ec.Trades{}
				}
			}

			if th.Len() < 100 { // 当th的长度较短时，是由于已经读取到最新的消息了。
				waitMS = 1000 * 60
			} else {
				waitMS = o.CoinPeriodS
			}
			util.HoldOn(time.Duration(waitMS)*time.Millisecond, &requestTime)
		}
	}()
}

func (o *BTC38) post(method string, v url.Values, result interface{}) (err error) {
	type Response struct {
		Result    bool  `json:"result"`
		ErrorCode int64 `json:"error_code"`
	}

	v.Set("api_key", o.PublicKey)
	hasher := ec.MD5([]byte(v.Encode() + "&secret_key=" + o.SecretKey))
	v.Set("sign", strings.ToUpper(ec.HexEncodeToString(hasher)))
	encoded := v.Encode()
	path := apiURL + method

	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	ansChan := make(chan answer)
	o.ask <- ask{Type: post,
		Path:       ec.Path(path, v),
		Headers:    headers,
		Body:       strings.NewReader(encoded),
		AnswerChan: ansChan}
	ans := <-ansChan

	if ans.err != nil {
		return ans.err
	}

	// err = ec.JSONDecode([]byte(ans.body), &result)
	// if err != nil {
	// 	if o.ShowDetail {
	// 		log.Println(string(ans.body))
	// 	}
	// 	return err
	// }

	err = ec.JSONDecode([]byte(ans.body), &result)
	if err != nil {
		str := err.Error()

		r := new(Response)
		err = ec.JSONDecode([]byte(ans.body), r)
		if err != nil {
			str = str + " AND " + err.Error()
			return errors.New(str)
		}

		// if r.ErrorCode > 0 {
		// 	s := fmt.Sprintln("失败原因:", o.restErrors[r.ErrorCode])
		// 	return errors.New(s)
		// }
		return errors.New(str)
	}

	return nil
}
