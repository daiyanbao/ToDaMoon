package exchanges

//API 交易所的标准接口
//每一个子交易所，都要求返回符合这个接口的子例
type API interface {
	Name() string
	Ticker(money, coin string) (*Ticker, error)
	Depth(money, coin string) (*Depth, error)
	TransRecords(money, coin string, tid int64) (Trades, error)
	Account() (Account, error)
	Trans(t TransType, money, coin string, price, amount float64) (int, error)
	CancelOrder(money, coin string, orderID int) (bool, error)
	MyOrders(money, coin string) ([]Order, error)
	MyTransRecords(money, coin string, tid int64) (Trades, error)
}

//TransType 指定了交易的类型
type TransType string

const (
	//BUY 是使用money换coin的过程
	BUY TransType = "buy"
	//SELL 是使用coin换money的过程
	SELL TransType = "sell"
)

//TestAPI 用于测试通用API接口的功能
func TestAPI(a API) {

}
