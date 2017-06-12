package btc38

import (
	"ToDaMoon/exchanges"
	"ToDaMoon/util"
	"fmt"
	"log"
	"time"

	"github.com/BurntSushi/toml"
)

var btc38 *BTC38
var name = "btc38"

//BTC38 包含了btc38.com的API所需的所有数据
type BTC38 struct {
	*config
	*exchanges.Net
	exchanges.TradesDBs
	exchanges.TradesSubject
	*exchanges.Account
}

type config struct {
	ShowDetail        bool
	RecordHistory     bool   //使用数据库记录历史交易记录
	MinAccessPeriodMS int    //两次API访问的最小间隔时间，单位为毫秒
	CoinPeriodS       int    //查询某一个币种最新交易记录的时间间隔，单位为秒
	IP                string //本机的公网ip，btc38要求访问API的公网ip在网站上注册过。
	ID                int
	Name              string
	PublicKey         string
	SecretKey         string
	DBDir             string
	Markets           map[string][]string //key是money，value是coins
}

//instance 返回一个btc38的单例
func instance() *BTC38 {
	//读取配置文件
	cfg := getConfig()

	//生成btc38实例
	btc38 = generateBy(cfg)

	//获取btc38各个coin的全局交易的最新数据到数据库，然后，发布最新全局交易数据订阅
	btc38.TradesSubject = exchanges.MakeSubjectes(btc38, btc38.TradesDBs, time.Millisecond*100, time.Minute)

	return btc38
}

func getConfig() *config {
	filename := getConfigFileName()

	cfg := new(config)
	if _, err := toml.DecodeFile(filename, cfg); err != nil {
		msg := fmt.Sprintf("无法加载%s/%s，并Decode到cfg变量: %s", util.PWD(), filename, err)
		log.Fatalf(msg)
	}

	cfg.check()
	return cfg
}

func getConfigFileName() string {
	return fmt.Sprintf("%s.toml", name)
}

// Check Config for setting mistakes
func (c *config) check() {
	if len(c.PublicKey) != 32 {
		log.Fatalln("btc38的PublicKey长度应为32位")
	}

	if len(c.SecretKey) != 64 {
		log.Fatalln("btc38的SecretKey长度应为64位")
	}

	if c.MinAccessPeriodMS < 10 {
		log.Fatalln("btc38的API访问间隔等待时间过短，请核查")
	}

	myIP, err := util.ExternalIP()
	if err != nil {
		log.Fatal("无法获取本机的外网IP")
	}
	if myIP != c.IP {
		text := fmt.Sprintf("本机外网IP地址%s没有在BTC38网注册", myIP)
		notify.Bad(text)
		log.Fatalf(text)
	}
}

func generateBy(c *config) *BTC38 {
	n := &exchanges.Net{
		Header: genHeader(),
	}
	n.Start(c.MinAccessPeriodMS)

	tdb := exchanges.MakeTradesDBs(c.DBDir, c.Name, c.Markets)

	btc38 = &BTC38{config: c,
		Net:       n,
		TradesDBs: tdb,
	}
	return btc38
}

func genHeader() map[string]string {
	header := make(map[string]string)
	header["Content-Type"] = "application/x-www-form-urlencoded"
	header["User-Agent"] = "Mozilla/4.0"
	return header
}

//Name 返回BTC38的name
func (b *BTC38) Name() string {
	return b.config.Name
}
