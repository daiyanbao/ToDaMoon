package btc38

import (
	"ToDaMoon/util"
	"io"
	"log"

	db "ToDaMoon/DataBase"
	ec "ToDaMoon/exchanges"

	"fmt"

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
	RecordHistory     bool   //使用数据库记录历史交易记录
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
		text := fmt.Sprintf("本机外网IP地址%s没有在BTC38网注册", myIP)
		notify.Bad(text)
		log.Fatalf(text)
	}
}

// BTC38 包含了模块所需的基本参数
type BTC38 struct {
	*Config
	ask      askChannel
	db       map[string]db.DBM
	Property map[string]observer.Property
}

func getMD5(time string) string {
	md := fmt.Sprintf("%s_%d_%s_%s", btc38.PublicKey, btc38.ID, btc38.SecretKey, time)
	md5 := ec.MD5([]byte(md))
	return ec.HexEncodeToString(md5)
}
