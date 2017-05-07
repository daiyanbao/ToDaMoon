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

// 以下代码实现了okcoin模块的单例特性
var okcoin *OKCoin
var once sync.Once

// Instance make a singleton of okcoin.cn
func Instance(cfg *Config, notify Interface.Notify) *OKCoin {
	cfg.Check()
	once.Do(func() {
		askChan := make(askChannel, 12)
		okcoin = &OKCoin{Config: cfg,
			ask:      askChan,
			db:       map[string]db.DBM{},
			Property: map[string]observer.Property{},
			notify:   notify}
		okcoin.setErrors()
		go start(askChan, time.Duration(cfg.APIAccessPeriodMS)*time.Millisecond)

		okcoin.makeDBs()
		okcoin.makePropertys()
		notify.Info("单例初始化完成。")
	})

	return okcoin
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

func (o *OKCoin) makeDBs() {
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

func (o *OKCoin) makePropertys() {
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

func listeningTradeHistoryAndSave(o *OKCoin, coin string) {
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
	waitMS := o.ListenMS

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
				waitMS = o.ListenMS
			}
			util.HoldOn(time.Duration(waitMS)*time.Millisecond, &requestTime)
		}
	}()
}

func (o *OKCoin) post(method string, v url.Values, result interface{}) (err error) {
	type Response struct {
		Result    bool  `json:"result"`
		ErrorCode int64 `json:"error_code"`
	}

	v.Set("api_key", o.APIKey)
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

		if r.ErrorCode > 0 {
			s := fmt.Sprintln("失败原因:", o.restErrors[r.ErrorCode])
			return errors.New(s)
		}
		return errors.New(str)
	}

	return nil
}

// setErrors
func (o *OKCoin) setErrors() {
	o.restErrors = map[int64]string{
		10000: "必选参数不能为空",
		10001: "用户请求过于频繁",
		10002: "系统错误",
		10003: "未在请求限制列表中,稍后请重试",
		10004: "IP限制不能请求该资源",
		10005: "密钥不存在",
		10006: "用户不存在",
		10007: "签名不匹配",
		10008: "非法参数",
		10009: "订单不存在",
		10010: "余额不足",
		10011: "买卖的数量小于BTC/LTC最小买卖额度",
		10012: "当前网站暂时只支持btc_cny ltc_cny",
		10013: "此接口只支持https请求",
		10014: "下单价格不得≤0或≥1000000",
		10015: "下单价格与最新成交价偏差过大",
		10016: "币数量不足",
		10017: "API鉴权失败",
		10018: "借入不能小于最低限额[cny:100,btc:0.1,ltc:1]",
		10019: "页面没有同意借贷协议",
		10020: "费率不能大于1%",
		10021: "费率不能小于0.01%",
		10023: "获取最新成交价错误",
		10024: "可借金额不足",
		10025: "额度已满，暂时无法借款",
		10026: "借款(含预约借款)及保证金部分不能提出",
		10027: "修改敏感提币验证信息，24小时内不允许提现",
		10028: "提币金额已超过今日提币限额",
		10029: "账户有借款，请撤消借款或者还清借款后提币",
		10031: "存在BTC/LTC充值，该部分等值金额需6个网络确认后方能提出",
		10032: "未绑定手机或谷歌验证",
		10033: "服务费大于最大网络手续费",
		10034: "服务费小于最低网络手续费",
		10035: "可用BTC/LTC不足",
		10036: "提币数量小于最小提币数量",
		10037: "交易密码未设置",
		10040: "取消提币失败",
		10041: "提币地址未认证",
		10042: "交易密码错误",
		10043: "合约权益错误，提币失败",
		10044: "取消借款失败",
		10047: "当前为子账户，此功能未开放",
		10048: "提币信息不存在",
		10049: "小额委托（<0.5BTC)的未成交委托数量不得大于50个",
		10050: "重复撤单",
		10060: "您的提现功能被冻结，请联系客服!",
		10100: "账户被冻结",
		10101: "订单类型错误",
		10102: "不是本用户的订单",
		10103: "私密订单密钥错误",
		10104: "系统检测到您有可疑操作，暂时不可进行大宗交易!",
		10216: "非开放API",
		503:   "用户请求过于频繁(Http)",
	}
}
