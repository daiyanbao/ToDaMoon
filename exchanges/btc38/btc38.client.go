package btc38

import (
	"ToDaMoon/Interface"
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
	ShowDetail                     bool
	APIAccessPeriodMS              int //两次API访问的最小间隔时间
	MinListenMS                    int
	Name, APIKey, SecretKey, DBDir string
	Coins                          []string
}

// Check Config for setting mistakes
func (c *Config) Check() {
	if len(c.APIKey) != 36 {
		log.Fatalln("Settings.ini -> okcoin.cn -> APIKey: 长度应为36位")
	}

	if len(c.SecretKey) != 32 {
		log.Fatalln("Settings.ini -> okcoin.cn -> SecretKey: 长度应为32位")
	}

	if c.APIAccessPeriodMS < 10 {
		log.Fatalln("Settings.ini -> okcoin.cn -> WaitMillisecond: 等待时间过短，请核查")
	}
}

// OKCoin 包含了模块所需的基本参数
type OKCoin struct {
	*Config
	ask        askChannel
	restErrors map[int64]string
	db         map[string]db.DBM
	Property   map[string]observer.Property
	notify     Interface.Notify
}
