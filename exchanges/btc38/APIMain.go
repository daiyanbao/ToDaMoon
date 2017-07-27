package btc38

import (
	"fmt"
	"log"

	ec "github.com/aQuaYi/ToDaMoon/exchanges"
	"github.com/aQuaYi/ToDaMoon/pubu"
	"github.com/aQuaYi/ToDaMoon/util"

	"sync"

	"github.com/BurntSushi/toml"
)

const (
	name              = "btc38"
	baseURL           = "http://api.btc38.com/v1/"
	tickerURL         = baseURL + "ticker.php"
	depthURL          = baseURL + "depth.php"
	transRecordsURL   = baseURL + "trades.php"
	myAccountURL      = baseURL + "getMyBalance.php"
	submitOrderURL    = baseURL + "submitOrder.php"
	cancelOrderURL    = baseURL + "cancelOrder.php"
	getOrderListURL   = baseURL + "getOrderList.php"
	getMyTradeListURL = baseURL + "getMyTradeList.php"
)

var api *API
var onceAPI sync.Once

//API 包含了btc38.com的所有API的wrapper
type API struct {
	*config
	*ec.Net
}

//NewAPI 返回一个*API的单例
func NewAPI() *API {

	onceAPI.Do(
		func() {
			// 使用pubu.im作为通知工具
			notify = pubu.New()

			//读取配置文件
			cfg := getConfig()

			//根据配置生成*API实例
			api = genAPI(cfg)
		})

	return api
}

//TODO: 把config分解成
// type config struct {
// 	APICfg
// 	DebugCfg
// 	DBCfg
// 	Markets
// }

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

func getConfig() *config {
	filename := ec.GetConfigFilename(name)

	cfg := new(config)
	if _, err := toml.DecodeFile(filename, cfg); err != nil {
		msg := fmt.Sprintf("无法加载%s/%s，并Decode到cfg变量: %s", util.PWD(), filename, err)
		log.Fatalf(msg)
	}

	cfg.check()
	return cfg
}

// check Config for setting mistakes
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
		log.Fatal("无法获取本机的外网IP:", err)
	}

	if myIP != c.IP {
		text := fmt.Sprintf("本机外网IP地址%s没有在BTC38网注册", myIP)
		log.Fatalf(text)
	}
}

func genAPI(c *config) *API {
	n := &ec.Net{
		Header: genHeader(),
	}
	n.Start(c.MinAccessPeriodMS)

	a := &API{config: c,
		Net: n,
	}

	return a
}

func genHeader() map[string]string {
	header := make(map[string]string)
	header["Content-Type"] = "application/x-www-form-urlencoded"
	header["User-Agent"] = "Mozilla/4.0"
	return header
}

//Name 返回API所在交易所的名字
func (a *API) Name() string {
	return a.config.Name
}