package btc38

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/aQuaYi/ToDaMoon/exchanges"
)

// MyAccount 返回BTC38的账户信息
func (a *API) MyAccount() (*exchanges.Account, error) {
	rawData, err := a.myAccountRawData()
	if err != nil {
		msg := fmt.Sprintf("无法获取 %s的MyBalance的RawDate:%s", a.Name(), err)
		return nil, errors.New(msg)
	}

	if a.IsLog {
		log.Printf(`rawData %s.MyAccount()=%s`, a.Name(), string(rawData))
		log.Println()
	}

	m, err := handleMyAccountRawData(rawData)
	if err != nil {
		msg := fmt.Sprintf("无法转换MyBalance的rawData(%s):%s", string(rawData), err)
		return nil, errors.New(msg)
	}

	if a.IsLog {
		log.Printf(`After json.Unmarshal:%s.MyAccount()=%v`, a.Name(), m)
		log.Println()
	}

	return m.normalize(a.Markets["cny"])
}

func (a *API) myAccountRawData() ([]byte, error) {
	body := a.myAccountBody()
	return a.Post(myAccountURL, body)
}

func (a *API) myAccountBody() io.Reader {
	v := url.Values{}
	v.Set("key", a.PublicKey)
	nowTime := fmt.Sprint(time.Now().Unix())
	v.Set("time", nowTime)
	md5 := a.md5(nowTime)
	v.Set("md5", md5)

	encoded := v.Encode()
	return strings.NewReader(encoded)
}

func handleMyAccountRawData(rawData []byte) (myBalance, error) {
	resp := myBalance{}
	err := json.Unmarshal(rawData, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

var orderTypeMap = map[exchanges.OrderType]int{
	exchanges.BUY:  1,
	exchanges.SELL: 2,
}

// Order 下单交易
func (a *API) Order(t exchanges.OrderType, money, coin string, price, amount float64) (int64, error) {
	ot := orderTypeMap[t]
	rawData, err := a.orderRawData(ot, money, coin, price, amount)
	if err != nil {
		return 0, err
	}

	if a.IsLog {
		log.Printf(`rawData %s.Order(%s, %s, %s, %f, %f)=%s`, a.Name(), t, money, coin, price, amount, string(rawData))
		log.Println()
	}

	return handleOrderRawData(rawData)
}

func (a *API) orderRawData(ot int, money, coin string, price, amount float64) ([]byte, error) {
	//NOTICE: btc38的价格**可能**只支持6位有效数字。比如现在BTC的价格是18000多，如果买价是12345.6就可以下单，如果买价是12345.67，就会显示deciError1
	priceStr := price2Str(price)
	amountStr := fmt.Sprintf("%.6f", amount)

	body := a.orderBody(ot, money, coin, priceStr, amountStr)

	if a.IsLog {
		log.Printf(`Body of %s.Order(%d, %s, %s, %f, %f)=%v`, a.Name(), ot, money, coin, price, amount, body)
	}

	return a.Post(submitOrderURL, body)
}

func (a *API) orderBody(ot int, money, coin, price, amount string) io.Reader {
	v := url.Values{}
	v.Set("key", a.PublicKey)
	nowTime := fmt.Sprint(time.Now().Unix())
	v.Set("time", nowTime)
	md5 := a.md5(nowTime)
	v.Set("md5", md5)

	v.Set("type", fmt.Sprint(ot))
	v.Set("mk_type", money)
	v.Set("price", price)
	v.Set("amount", amount)
	v.Set("coinname", coin)
	encoded := v.Encode()

	if a.IsLog {
		log.Printf("encoded url.Values of %s.Order(%d, %s, %s, %s, %s)=%s\n", a.Name(), ot, money, coin, price, amount, encoded)
	}

	return strings.NewReader(encoded)
}

/*
http://www.btc38.com/help/document/2581.html
返回 / Return：
succ，挂单成功 / successful (submitOrder.php)
succ|123，挂单成功，123为您挂单的ID / successful, 123 is your order_id
overBalance，账户余额不足 / insuffient balance
其它返回表示不同的错误，情况太多，暂不罗列，这种错误发生的可能性不大
*/
func handleOrderRawData(rawData []byte) (int64, error) {
	r := string(rawData)

	//下单后，立马成交
	if r == "succ" {
		return 0, nil
	}

	//下单后，没有成交
	if r[:5] == "succ|" {
		orderID, err := strconv.Atoi(r[5:])
		if err != nil {
			return 0, err
		}
		return int64(orderID), nil
	}

	return 0, errors.New(r)
}

// CancelOrder 下单交易
func (a *API) CancelOrder(money, coin string, orderID int64) (bool, error) {
	rawData, err := a.cancelOrderRawData(money, coin, orderID)
	if err != nil {
		return false, err
	}

	if a.IsLog {
		log.Printf(`rawData %s.CancelOrder(%s, %s, %d)=%s`, a.Name(), money, coin, orderID, string(rawData))
	}

	return handleCancelOrderRawData(rawData)
}

func (a *API) cancelOrderRawData(money, coin string, orderID int64) ([]byte, error) {
	body := a.cancelOrderBody(money, coin, orderID)
	return a.Post(cancelOrderURL, body)
}

func (a *API) cancelOrderBody(money, coin string, orderID int64) io.Reader {
	v := url.Values{}
	v.Set("key", a.PublicKey)
	nowTime := fmt.Sprint(time.Now().Unix())
	v.Set("time", nowTime)
	md5 := a.md5(nowTime)
	v.Set("md5", md5)

	v.Set("mk_type", money)
	v.Set("order_id", fmt.Sprint(orderID))
	v.Set("coinname", coin)
	encoded := v.Encode()

	return strings.NewReader(encoded)
}

/*
http://www.btc38.com/help/document/2581.html
返回：
succ，撤单成功 / successful
overtime，该单不存在，或者已成交了 / order expired or traded
*/
func handleCancelOrderRawData(rawData []byte) (bool, error) {
	r := string(rawData)

	if r == "succ" {
		return true, nil
	}

	return false, errors.New(r)
}

// MyOrders 我所有还没有成交的挂单
func (a *API) MyOrders(money, coin string) ([]exchanges.Order, error) {
	rawData, err := a.myOrdersRawData(money, coin)
	if err != nil {
		return nil, err
	}

	if a.IsLog {
		log.Printf(`rawData %s.MyOrders(%s, %s)=%s`, a.Name(), money, coin, string(rawData))
	}

	return a.handleMyOrdersRawData(rawData, money)
}

func (a *API) myOrdersRawData(money, coin string) ([]byte, error) {
	body := a.myOrdersBody(money, coin)
	return a.Post(getOrderListURL, body)
}

func (a *API) myOrdersBody(money, coin string) io.Reader {
	v := url.Values{}
	v.Set("key", a.PublicKey)
	nowTime := fmt.Sprint(time.Now().Unix())
	v.Set("time", nowTime)
	md5 := a.md5(nowTime)
	v.Set("md5", md5)

	v.Set("mk_type", money)
	v.Set("coinname", coin)
	encoded := v.Encode()

	return strings.NewReader(encoded)
}

/*
http://www.btc38.com/help/document/2581.html
返回：
如果挂单为空，返回 "no_order"
如果挂单不为空，则返回比如 / if there is any order：
[{"order_id":"123", "order_type":"1", "order_coinname":"BTC", "order_amount":"23.232323", "order_price":"0.2929"}, {"order_id":"123", "order_type":"1", "order_coinname":"LTC","order_amount":"23.232323", "order_price":"0.2929"}]
*/
func (a *API) handleMyOrdersRawData(rawData []byte, money string) ([]exchanges.Order, error) {
	if string(rawData) == "no_order" {
		return []exchanges.Order{}, nil
	}

	resp := []order{}
	err := json.Unmarshal(rawData, &resp)
	if err != nil {
		msg := fmt.Sprintf("json.Unmarshal %s时，出错:%s", string(rawData), err)
		return nil, errors.New(msg)
	}

	if a.IsLog {
		log.Printf("\nBefore json.Unmarshal %s\nAfter json.Unmarshal %v\n", rawData, resp)
	}

	result := make([]exchanges.Order, len(resp))
	for i, v := range resp {
		r, err := v.normalize(money)
		if err != nil {
			msg := fmt.Sprintf("处理MyOrders的rawData失败，无法normalize:%s", err)
			return nil, errors.New(msg)
		}
		result[i] = *r
	}

	if a.IsLog {
		log.Printf("\nAfter normalize %v\n", result)
	}

	return result, nil
}

// MyTransRecords 获取我的交易记录
// tid参数应为数据库中，真实存在的tid
// 如果，不知道真实的tid，可以让tid=0来获取全部的交易数据。
func (a *API) MyTransRecords(money, coin string, tid int64) (exchanges.Trades, error) {
	next := nextECTrades(a, money, coin)
	res := exchanges.Trades{}
	done := false
	//res := exchanges.Trades{}

	for !done {
		temp, err := next()
		if err != nil {
			msg := fmt.Sprintf("获取next()失败:%s", err)
			return nil, errors.New(msg)
		}

		//当tid不是真实存在的tid时，
		//例如，tid=0时，会反复读取数据，直到返回的temp为空切片。
		if temp.Len() == 0 {
			break
		}
		//获取myTrade数据，len(mt)==0 return
		res, done = appendECTrades(res, temp, tid)
	}

	res.Sort()
	return res, nil
}

// nextECTrades
func nextECTrades(a *API, money, coin string) func() (exchanges.Trades, error) {
	page := 0

	return func() (exchanges.Trades, error) {
		mts, err := a.MyTradeList(money, coin, page)
		if err != nil {
			msg := fmt.Sprintf("nextPageList 无法获取%s.MyTradeList(%s, %s, %d)的数据: %s", a.Name(), money, coin, page, err)
			return nil, errors.New(msg)
		}

		res, err := myTrades2ECTrades(mts, a.ID)
		if err != nil {
			msg := fmt.Sprintf("无法把nextECTrades获取的%s转换成exchanges.Trades: %s", mts, err)
			return nil, errors.New(msg)
		}

		page++
		return res, nil
	}
}

func myTrades2ECTrades(mts []MyTrade, ID int) (exchanges.Trades, error) {
	res := make(exchanges.Trades, len(mts))
	for i, mt := range mts {
		et, err := mt.normalize(ID)
		if err != nil {
			msg := fmt.Sprintf("无法把%s转换成exchanges.Trade: %s", mt, err)
			return nil, errors.New(msg)
		}
		res[i] = et
	}

	//mts是降序，exchanges.Trades是升序，所以res要进行排序。
	res.Sort()
	return res, nil
}

func appendECTrades(res, temp exchanges.Trades, tid int64) (exchanges.Trades, bool) {
	temp.Sort()
	switch {
	case tid < temp[0].Tid:
		res = append(res, temp...)
		return res, false //没有找完
	case tid < temp[temp.Len()-1].Tid:
		res = append(res, temp.After(tid)...)
		return res, true //已经找完了。
	default:
		return res, true //已经找完了。
	}
}

// MyTradeList 按照btc38的API的格式，返回交易记录结果。
func (a *API) MyTradeList(money, coin string, page int) ([]MyTrade, error) {
	rawData, err := a.myTradeListRawData(money, coin, page)
	if err != nil {
		return nil, err
	}

	return a.handleMyTradeListRawData(rawData)
}

func (a *API) myTradeListRawData(money, coin string, page int) ([]byte, error) {
	body := a.myTradeListBody(money, coin, page)
	return a.Post(getMyTradeListURL, body)
}

func (a *API) myTradeListBody(money, coin string, page int) io.Reader {
	v := url.Values{}
	v.Set("key", a.PublicKey)
	nowTime := fmt.Sprint(time.Now().Unix())
	v.Set("time", nowTime)
	md5 := a.md5(nowTime)
	v.Set("md5", md5)

	v.Set("mk_type", money)
	v.Set("coinname", coin)
	v.Set("page", fmt.Sprint(page))
	encoded := v.Encode()

	return strings.NewReader(encoded)
}

func (a *API) handleMyTradeListRawData(rawData []byte) ([]MyTrade, error) {
	resp := []MyTrade{}
	err := json.Unmarshal(rawData, &resp)
	if err != nil {
		return nil, err
	}

	if a.IsLog {
		log.Printf("\n交易记录明细 %s\n", resp)
	}

	return resp, nil
}

// CheckMyTradeList 测试btc38.MyTradeList()
func CheckMyTradeList(a *API, money, coin string, page int) (result string) {
	method := fmt.Sprintf(`%s.MyTradeList("%s", "%s", %d)`, a.Name(), money, coin, page)

	fmt.Printf("==测试%s==\n", method)

	mtl, err := a.MyTradeList(money, coin, page)
	if err != nil {
		result = fmt.Sprintf("%s Error:%s\n", method, err)
		return
	}

	fmt.Printf("%s=", method)
	if len(mtl) < 5 {
		fmt.Println(mtl)
	} else {
		fmt.Println(mtl[:2])
		fmt.Println("... ... ... ...")
		fmt.Println(mtl[len(mtl)-2:])
	}

	return
}

func (a *API) md5(time string) string {
	md := fmt.Sprintf("%s_%d_%s_%s", a.PublicKey, a.ID, a.SecretKey, time)
	md5 := exchanges.MD5([]byte(md))
	return exchanges.HexEncodeToString(md5)
}
