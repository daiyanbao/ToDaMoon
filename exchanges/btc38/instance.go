package btc38

import (
	"ToDaMoon/util"
	"fmt"
	"log"

	"github.com/go-ini/ini"
)

var btc38 *BTC38
var name = "btc38"

//BTC38 包含了btc38.com的API所需的所有数据
type BTC38 struct {
}

type config struct {
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

//instance 返回一个btc38的单例
func instance() *BTC38 {
	//读取配置文件

	//检查配置文件
	c := getConfig()

	//生成btc38实例
	btc38 = generate(c)
	//连接btc38的全局交易数据库

	//连接btc38的本人交易数据库

	//获取btc38各个coin的全局交易的最新数据到数据库，然后，发布最新全局交易数据订阅

	return btc38
}

func getConfig() *config {
	//TODO: 把配置文件改成toml格式的

	//获取cfg文件
	iniName := getIniFileName()
	iniFile, err := ini.Load(iniName)
	if err != nil {
		msg := fmt.Sprintf("无法加载%s/%s.ini: %s", util.PWD(), name, err)
		log.Fatalf(msg)
	}

	//生成配置对象
	cfg := new(config)
	if err := iniFile.Section(name).MapTo(cfg); err != nil {
		msg := fmt.Sprintf("无法Map设置的参数内容到%s的配置对象", name)
		log.Fatalln(msg, err)
	}

	cfg.check()
	return cfg
}

func getIniFileName() string {
	return fmt.Sprintf("./%s.ini", name)
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
		log.Fatal("无法获取本机外部IP")
	}
	if myIP != c.IP {
		text := fmt.Sprintf("本机外网IP地址%s没有在BTC38网注册", myIP)
		notify.Bad(text)
		log.Fatalf(text)
	}
}

func generate(c *config) *BTC38 {

}

