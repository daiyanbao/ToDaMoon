package btc38

var btc38 *BTC38

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

	//生成btc38实例

	//连接btc38的全局交易数据库

	//连接btc38的本人交易数据库

	//获取btc38各个coin的全局交易的最新数据到数据库，然后，发布最新全局交易数据订阅

	return nil
}
