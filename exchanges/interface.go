package exchanges

//Exchanger 交易所的标准接口
//每一个子交易所，都要求返回符合这个接口的子例
type Exchanger interface {
	Name() string
	Ticker(money, coin string) (*Ticker, error)
	Trades(money, coin string, tid int64) (Trades, error)
}
