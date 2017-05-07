package btc38

import (
	"ToDaMoon/util"
	"io"
	"log"

	db "ToDaMoon/DataBase"

	observer "github.com/imkira/go-observer"
)

type askChannel chan ask

type ask struct {
	Type       askType
	Path       string
	Headers    map[string]string
	Body       io.Reader
	AnswerChan chan<- answer
}

type askType int

const (
	get askType = iota
	post
)

type answer struct {
	body []byte
	err  error
}

// Config contains all the ini settings
type Config struct {
	ShowDetail        bool
	MinAccessPeriodMS int    //两次API访问的最小间隔时间，单位为毫秒
	CoinPeriodS       int    //查询某一个币种最新交易记录的时间间隔，单位为秒
	IP                string //本机ip，btc38
	ID                int
	Name              string
	PublicKey         string
	SecretKey         string
	DBDir             string
	Coins             []string
}

// Check Config for setting mistakes
func (c *Config) Check() {
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
		log.Fatal("无法获取本机外部IP")
	}
	if myIP != c.IP {
		log.Fatal("本机外网IP地址没有在BTC38网注册")
	}
}

// BTC38 包含了模块所需的基本参数
type BTC38 struct {
	*Config
	ask      askChannel
	db       map[string]db.DBM
	Property map[string]observer.Property
}
