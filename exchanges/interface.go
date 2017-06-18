package exchanges

//API 交易所的标准接口
//每一个子交易所，都要求返回符合这个接口的子例
type API interface {
	Name() string
	Ticker(money, coin string) (*Ticker, error)
	Depth(money, coin string) (*Depth, error)
	TransRecords(money, coin string, tid int64) (Trades, error)
	MyAccount() (*Account, error)
	Order(t OrderType, money, coin string, price, amount float64) (int64, error)
	CancelOrder(money, coin string, orderID int64) (bool, error)
	MyOrders(money, coin string) ([]Order, error)
	MyTransRecords(money, coin string, tid int64) (Trades, error)
}

//TODO: 把以下内容，移入struct.go

//OrderType 指定了交易的类型
type OrderType string

const (
	//BUY 是使用money换coin的过程
	BUY OrderType = "buy"
	//SELL 是使用coin换money的过程
	SELL OrderType = "sell"
)

//TestAPI 用于测试通用API接口的功能
func TestAPI(a API) {

}
