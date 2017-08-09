package exchanges

import (
	"fmt"
	"time"
)

// API 交易所的标准接口
// 每一个子交易所，都要求返回符合这个接口的子例
type API interface {
	// 交易所的名称
	Name() string

	// 反馈coin的交易指
	Ticker(money, coin string) (*Ticker, error)

	// 反馈市场中币种双方要价，已经排序过了
	// Asks[0].Price是最低的卖价
	// Bids[0].Price是最高的买价
	Depth(money, coin string) (*Depth, error)

	// 返回在Tid之后的一组全局交易记录，
	// 不同的交易所，返回的长度不一样。
	TransRecords(money, coin string, tid int64) (Trades, error)

	// 返回你在exchange的各个币的额度。
	MyAccount() (*Account, error)

	// 如果下单成功，会返回订单编号。
	Order(t OrderType, money, coin string, price, amount float64) (int64, error)

	CancelOrder(money, coin string, orderID int64) (bool, error)

	// 还没有成交的挂单
	MyOrders(money, coin string) ([]Order, error)

	// MyTransRecords 返回在Tid之后的成交记录
	// MyTransRecords 和 TransRecords 的逻辑是一致的。
	MyTransRecords(money, coin string, Tid int64) (Trades, error)
}

/*
以下test开头的函数，不是测试函数，是帮助具体交易所测试的函数。
*/

// CheckTicker 测试API.Ticker()
func CheckTicker(a API, money, coin string) (price float64, result string) {
	method := fmt.Sprintf(`%s.Ticker("%s","%s")`, a.Name(), money, coin)

	fmt.Printf("==测试%s==\n", method)

	t, err := a.Ticker(money, coin)
	if err != nil {
		result = fmt.Sprintf("%s Error:%s\n", method, err)
		return
	}

	fmt.Printf("%s\n%v\n", method, t)
	price = t.Last

	return
}

// CheckDepth 测试API.Depth()
func CheckDepth(a API, money, coin string) (result string) {
	method := fmt.Sprintf(`%s.Depth("%s","%s")`, a.Name(), money, coin)

	fmt.Printf("==测试%s==\n", method)

	d, err := a.Depth(money, coin)
	if err != nil {
		result = fmt.Sprintf("%s Error:%s\n", method, err)
		return
	}

	fmt.Printf("%s\n%s\n", method, d)
	return
}

// CheckTransRecords 测试API.TransRecords
func CheckTransRecords(a API, money, coin string, tid int64) (result string) {
	method := fmt.Sprintf(`%s.TransRecords("%s","%s", %d)`,
		a.Name(), money, coin, tid)

	fmt.Printf("==测试%s==\n", method)

	tr, err := a.TransRecords(money, coin, tid)
	if err != nil {
		result = fmt.Sprintf("%s Error:%s\n", method, err)
		return
	}

	fmt.Printf("%s\n", method)
	fmt.Println(tr[:2])
	fmt.Println("... ... ...")
	fmt.Println(tr[len(tr)-2:])

	return
}

// CheckMyAccount 测试API.MyAccount()
func CheckMyAccount(a API) (result string) {
	method := fmt.Sprintf(`%s.MyAccount()`, a.Name())

	fmt.Printf("==测试%s==\n", method)

	ma, err := a.MyAccount()
	if err != nil {
		result = fmt.Sprintf("%s Error:%s\n", method, err)
	}

	fmt.Printf("%s\n%s\n", method, ma)
	return
}

// CheckOrder 测试API.Order()
func CheckOrder(a API, ot OrderType, money, coin string, price, amount float64) (id int64, result string) {
	method := fmt.Sprintf(`%s.Order("%s", "%s", %f, %f)`,
		a.Name(), money, coin, price, amount)

	fmt.Printf("==测试%s===\n", method)

	id, err := a.Order(ot, money, coin, price, amount)
	if err != nil {
		result = fmt.Sprintf("%s Error:%s\n", method, err)
		return
	}

	fmt.Printf("%s下单成功，订单号%d", method, id)
	return
}

// CheckMyOrders 测试API.MyOrders()
func CheckMyOrders(a API, money, coin string) (result string) {
	method := fmt.Sprintf(`%s.MyOrders("%s", "%s")`, a.Name(), money, coin)

	fmt.Printf("==测试%s===\n", method)

	orders, err := a.MyOrders(money, coin)
	if err != nil {
		result = fmt.Sprintf("%s Error:%s\n", method, err)
		return
	}

	fmt.Printf("%s的挂单如下\n%v\n", method, orders)

	return
}

// CheckCancelOrder 测试API.CancelOrder()
func CheckCancelOrder(a API, money, coin string, id int64) (result string) {
	method := fmt.Sprintf(`%s.CancelOrder("%s", "%s", %d)`, a.Name(), money, coin, id)

	fmt.Printf("==测试%s==\n", method)

	if id == 0 {
		result = fmt.Sprintln("orderID==0，无订单可取消")
		return
	}

	fmt.Println("\t=====等待撤单=====")
	for i := 5; i > 0; i-- {
		fmt.Println(i)
		time.Sleep(time.Second)
	}

	canceled, err := a.CancelOrder(money, coin, id)
	if err != nil {
		result = fmt.Sprintf(`%s Error:%s\n`, method, err)
		return
	}

	if canceled {
		fmt.Printf("%s 撤单成功", method)
		return
	}

	result = fmt.Sprintln("err is nil, but canceled is false.")
	time.Sleep(3 * time.Second)

	return
}

// CheckMyTransRecords 测试API.MyTransRecords()
func CheckMyTransRecords(a API, money, coin string, sinceID int64) (result string) {
	method := fmt.Sprintf(`%s.MyTransRecords("%s", "%s", %d)`, a.Name(), money, coin, sinceID)

	fmt.Printf("\n===测试%s===\n", method)

	myRecords, err := a.MyTransRecords(money, coin, sinceID)
	if err != nil {
		result = fmt.Sprintf("%s Error:%s\n", method, err)
		return
	}

	fmt.Printf("%s=\n", method)
	fmt.Println(myRecords)

	return result
}
