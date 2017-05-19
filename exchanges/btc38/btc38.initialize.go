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

	"github.com/go-ini/ini"
	"github.com/imkira/go-observer"
)

var notify Interface.Notify

// 以下代码实现了BTC38模块的单例特性
var btc38 *BTC38
var once sync.Once

//Instance 返回 btc38的一个单例
func Instance() *BTC38 {
	once.Do(buildBTC38)
	return btc38
}

func buildBTC38() {
	notify = pubu.New()
	notify.Info("开始生成BTC38的Instance")

	//读取btc38的配置
	cfg, err := ini.Load("./btc38.ini")
	if err != nil {
		msg := fmt.Sprintf("无法加载%s/btc38.ini: %s", util.PWD(), err)
		notify.Error(msg)
		log.Fatalf(msg)
	}

	//生成btc38的配置对象
	btc38Cfg := new(Config)
	if err := cfg.Section("btc38").MapTo(btc38Cfg); err != nil {
		log.Fatalln("无法读取btc38的配置", err)
	}
	btc38Cfg.Check()

	//启动线程安全的请求通道
	askChan := make(askChannel, 12)
	go start(askChan, time.Duration(btc38Cfg.MinAccessPeriodMS)*time.Millisecond)

	btc38 = &BTC38{Config: btc38Cfg,
		ask:      askChan,
		db:       map[string]db.DBM{},
		Property: map[string]observer.Property{},
	}

	btc38.makeDBs()
	if btc38Cfg.RecordHistory {
		btc38.makePropertys()
	}
	log.Println("btc38的单例初始化完成。")
	notify.Good("btc38的单例初始化完成。")
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
		o.db[coin], err = db.Connect(o.DBDir, o.Name, coin, "cny")
		if err != nil {
			text := fmt.Sprintf("无法连接%s\\%s的中%s的数据库。\n", o.DBDir, o.Name, coin)
			notify.Error(text)
			log.Fatalln(text)
		}

		if o.ShowDetail {
			maxTid, _ := o.db[coin].MaxTid()
			text := fmt.Sprintf("已经链接上了%s的数据库，其最大Tid是%d\n", coin, maxTid)
			log.Println(text)
		}
	}

	if o.ShowDetail {
		notify.Debug("已经连接了所有相关的数据库")
	}
}

func (o *BTC38) makePropertys() {
	wg := &sync.WaitGroup{}
	for _, coin := range o.Coins {
		wg.Add(1)
		go listeningTradeHistoryAndSave(o, coin, wg)
	}

	wg.Wait()
	text := fmt.Sprintln("已经创建了所有相关的监听属性")
	log.Println(text)
	if o.ShowDetail {
		go notify.Debug(text)
	}
}

func listeningTradeHistoryAndSave(o *BTC38, coin string, wg *sync.WaitGroup) {
	defer wg.Done()
	maxTid, err := o.db[coin].MaxTid()
	if err != nil {
		text := fmt.Sprintf("%s: 没有读取到相关的数据库的最大值", coin)
		notify.Error(text)
		log.Fatalln(text)
	}
	if o.ShowDetail {
		text := fmt.Sprintf("%s的%s的MaxTid是%d", o.Name, coin, maxTid)
		log.Println(text)
	}
	th, err := o.TradeHistory(coin, maxTid)
	if err != nil {
		log.Fatalf("无法获取%s的%s的历史交易数据。\n", o.Name, coin)
	}
	o.Property[coin] = observer.NewProperty(th)
	if o.ShowDetail {
		text := fmt.Sprintf("%s的%s: 已经创建了监听属性。", o.Name, coin)
		log.Println(text)
	}

	var thdb ec.Trades
	saveTime := time.Now()
	requestTime := time.Now()
	waitS := o.CoinPeriodS

	go func() {
		for {
			if th.Len() > 0 {
				maxTid = th[th.Len()-1].Tid
			}
			th, err = o.TradeHistory(coin, maxTid)
			if err != nil {
				text := fmt.Sprintf("请求%s的%s的历史交易数据失败, 5秒后重试。\n%s", o.Name, coin, err)
				notify.Error(text)
				log.Println(text)
				time.Sleep(time.Second * 5)
				continue
			}

			if th.Len() > 0 {
				o.Property[coin].Update(th)
				thdb = append(thdb, th...)
			}
			if thdb.Len() > 0 {
				if thdb.Len() > 10*10000 || time.Since(saveTime) > time.Minute*30 {
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

			if th.Len() < 50 { // 当th的长度较短时，是由于已经读取到最新的消息了。
				waitS = 300
			} else {
				waitS = o.CoinPeriodS
			}
			util.HoldOn(time.Duration(waitS)*time.Millisecond, &requestTime)
		}
	}()
}

func (o *BTC38) post(method string, v url.Values, result interface{}) (err error) {
	type Response struct {
		Result    bool  `json:"result"`
		ErrorCode int64 `json:"error_code"`
	}

	v.Set("key", o.PublicKey)
	nowTime := fmt.Sprint(time.Now().Unix())
	v.Set("time", nowTime)
	md5 := getMD5(nowTime)
	v.Set("md5", md5)

	encoded := v.Encode()
	path := apiURL + method

	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	headers["User-Agent"] = "Mozilla/4.0"

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

	err = ec.JSONDecode([]byte(ans.body), &result)
	if err != nil {
		str := err.Error()

		r := new(Response)
		err = ec.JSONDecode([]byte(ans.body), r)
		if err != nil {
			str = str + " AND " + err.Error()
			return errors.New(str)
		}

		return errors.New(str)
	}

	return nil
}
