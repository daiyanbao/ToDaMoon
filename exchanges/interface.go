package exchanges

import (
	"ToDaMoon/util"
	"fmt"
	"time"
)

//API 交易所的标准接口
//每一个子交易所，都要求返回符合这个接口的子例
type API interface {
	//TODO: 为每个方法编写说明
	Name() string
	Ticker(money, coin string) (*Ticker, error)
	Depth(money, coin string) (*Depth, error)
	TransRecords(money, coin string, tid int64) (Trades, error)
	MyAccount() (*Account, error)
	Order(t OrderType, money, coin string, price, amount float64) (int64, error)
	CancelOrder(money, coin string, orderID int64) (bool, error)
	MyOrders(money, coin string) ([]Order, error)

	//MyTransRecords 返回在Date之后的成交记录
	MyTransRecords(money, coin string, date int64) (Trades, error)
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
func TestAPI(a API) string {
	result := ""

	fmt.Printf("==测试%s.Ticker()==\n", a.Name())
	t, err := a.Ticker("cny", "btc")
	if err != nil {
		msg := fmt.Sprintf("%s.Ticker Error:%s\n", a.Name(), err)
		result += msg
		fmt.Print(msg)
	} else {
		fmt.Printf("%s.Ticker of cny and btc\n%s\n", a.Name(), t)
	}
	btcPrice := t.Last

	fmt.Printf("==测试%s.Depth()==\n", a.Name())
	d, err := a.Depth("cny", "btc")
	if err != nil {
		msg := fmt.Sprintf("%s.Depth Error:%s\n", a.Name(), err)
		result += msg
		fmt.Print(msg)
	} else {
		fmt.Printf("%s.Depth of cny and btc\n%s\n", a.Name(), d)
	}

	fmt.Printf("==测试%s.TransRecords()==\n", a.Name())
	tr, err := a.TransRecords("cny", "btc", 1)
	if err != nil {
		msg := fmt.Sprintf("%s.TransRecords Error:%s\n", a.Name(), err)
		result += msg
		fmt.Print(msg)
	} else {
		fmt.Printf("%s.TransRecords of cny and btc Since Tid = 1\n", a.Name())
		fmt.Println(tr[:3])
		fmt.Println("... ... ...")
		fmt.Println(tr[len(tr)-2:])
	}

	fmt.Printf("==测试%s.MyAccount()==\n", a.Name())
	ma, err := a.MyAccount()
	if err != nil {
		msg := fmt.Sprintf("%s.MyAccount Error:%s\n", a.Name(), err)
		result += msg
		fmt.Print(msg)
	} else {
		fmt.Printf("%s.MyAccount\n%s\n", a.Name(), ma)
	}

	fmt.Printf("==测试%s.Order()===\n", a.Name())
	buyBTCPrice := btcPrice * 0.9
	buyMoney := 20.0
	buyAmount := buyMoney / buyBTCPrice
	orderID, err := a.Order(BUY, "cny", "btc", buyBTCPrice, buyAmount)
	if err != nil {
		msg := fmt.Sprintf(`%s.Order(BUY, "cny", "btc", %f,%f) Error:%s`, a.Name(), buyBTCPrice, buyAmount, err) + "\n"
		result += msg
		fmt.Print(msg)
	} else {
		fmt.Printf(`%s.Order(BUY, "cny", "btc", %f,%f) 下单成功，订单号%d`, a.Name(), buyBTCPrice, buyAmount, orderID)
		fmt.Println()
	}

	fmt.Printf("==测试%s.MyOrders()===\n", a.Name())
	orders, err := a.MyOrders("cny", "btc")
	if err != nil {
		msg := fmt.Sprintf(`%s.MyOrders("cny", "btc") Error:%s`, a.Name(), err) + "\n"
		result += msg
		fmt.Print(msg)
	} else {
		fmt.Printf(`%s.MyOrders("cny", "btc")的挂单如下`, a.Name())
		fmt.Printf("\n%v\n", orders)
	}

	fmt.Println("=====等待撤单=====")
	for i := 10; i > 0; i-- {
		fmt.Println(i)
		time.Sleep(time.Second)
	}

	fmt.Printf("==测试%s.CancelOrder()===\n", a.Name())
	if orderID == 0 {
		fmt.Println("orderID==0，无订单可取消")
	} else {
		canceled, err := a.CancelOrder("cny", "btc", orderID)
		if err != nil {
			msg := fmt.Sprintf(`%s.Order(BUY, "cny", "btc", %f,%f) Error:%s\n`, a.Name(), buyBTCPrice, buyAmount, err)
			result += msg
			fmt.Print(msg)
		} else if canceled {
			fmt.Printf(`%s.Order(BUY, "cny", "btc", %f,%f) 撤单功，订单号%d`, a.Name(), buyBTCPrice, buyAmount, orderID)
			fmt.Println()
		} else {
			fmt.Println("err is nil, but canceled is false.")
			time.Sleep(3 * time.Second)
		}
	}

	fmt.Printf("==测试%s.MyTransRecords()===\n", a.Name())
	money := "cny"
	coin := "doge"
	maxDate := int64(1)
	myRecords, err := a.MyTransRecords(money, coin, maxDate)
	if err != nil {
		msg := fmt.Sprintf(`%s.MyTransRecords("%s","%s", %d) Error:%s`, a.Name(), money, coin, maxDate, err)
		result += msg + "\n"
		fmt.Print(msg)
	} else {
		fmt.Printf("%s之后，%s的%s的成交记录\n", util.DateOf(maxDate), money, coin)
		fmt.Println(myRecords)
	}

	return result
}
