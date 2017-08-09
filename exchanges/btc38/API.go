package btc38

import (
	"fmt"
	"log"

	"sync"

	"github.com/aQuaYi/GoKit"
	"github.com/aQuaYi/ToDaMoon/exchanges"

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
var once sync.Once

//API 包含了btc38.com的所有API的wrapper
type API struct {
	*config
	exchanges.Net
}

//NewAPI 返回一个*API的单例
func NewAPI() *API {
	once.Do(
		func() {
			//读取配置文件
			cfg := getConfig()

			//根据配置生成*API实例
			api = &API{
				config: cfg,
				Net:    exchanges.NewNet(cfg.APISleepMS),
			}
		})

	return api
}

type config struct {
	Name       string
	DBDir      string
	IsLog      bool
	APISleepMS int //两次API访问的最小间隔时间，单位为毫秒
	ID         int
	IP         string
	PublicKey  string
	SecretKey  string
	Markets    map[string][]string //key是money，value是coins
}

func getConfig() *config {
	filename := exchanges.Config(name)

	cfg := new(config)
	if _, err := toml.DecodeFile(filename, cfg); err != nil {
		msg := fmt.Sprintf("无法加载%s/%s，并Decode到cfg变量: %s", GoKit.PWD(), filename, err)
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

	if c.APISleepMS < 10 {
		log.Fatalln("btc38的API访问间隔等待时间过短，请核查")
	}

	myIP, err := ExternalIP()
	if err != nil {
		log.Fatal("无法获取本机的外网IP:", err)
	}

	if myIP != c.IP {
		text := fmt.Sprintf("本机IP没在%s注册:%s", c.IP, myIP)
		log.Fatalf(text)
	}
}

//Name 返回API所在交易所的名字
func (a *API) Name() string {
	return a.config.Name
}
