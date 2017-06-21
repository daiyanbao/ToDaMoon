package exchanges

import (
	"ToDaMoon/util"
	"fmt"
	"time"
)

//API 交易所的标准接口
//每一个子交易所，都要求返回符合这个接口的子例
type API interface {
	//交易所的名称
	Name() string

	//反馈coin的交易指
	Ticker(money, coin string) (*Ticker, error)

	//TODO: 对返回结果进行排序
	//反馈市场中币种双方要价，已经排序过了
	//Asks[0]是最低的卖价
	//Bids[0]是最高的买价
	Depth(money, coin string) (*Depth, error)

	//返回在Tid之后的一组全局交易记录，
	//不同的交易所，返回的长度不一样。
	TransRecords(money, coin string, tid int64) (Trades, error)

	//返回你在exchange的各个币的额度。
	MyAccount() (*Account, error)

	//如果下单成功，会返回订单编号。
	Order(t OrderType, money, coin string, price, amount float64) (int64, error)

	//REVIEW: 想想返回bool是否可以改进。
	CancelOrder(money, coin string, orderID int64) (bool, error)

	//还没有成交的挂单
	MyOrders(money, coin string) ([]Order, error)

	//MyTransRecords 返回在Tid之后的成交记录
	//MyTransRecords 和 TransRecords 的逻辑是一致的。
	MyTransRecords(money, coin string, Tid int64) (Trades, error)
}

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
	maxTid := int64(1)
	myRecords, err := a.MyTransRecords(money, coin, maxTid)
	if err != nil {
		msg := fmt.Sprintf(`%s.MyTransRecords("%s","%s", %d) Error:%s`, a.Name(), money, coin, maxTid, err)
		result += msg + "\n"
		fmt.Print(msg)
	} else {
		fmt.Printf("%s之后，%s的%s的成交记录\n", util.DateOf(maxTid), money, coin)
		//TODO: 要求输出全部的
		fmt.Println(myRecords)
	}

	return result
}
