package btc38

import (
	ec "ToDaMoon/exchanges"
	"errors"
	"fmt"
	"io"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"
)

//MyAccount 返回BTC38的账户信息
func (a *API) MyAccount() (*ec.Account, error) {
	rawData, err := a.myAccountRawData()
	if err != nil {
		msg := fmt.Sprintf("无法获取%s的MyBalance的RawDate:%s", a.Name(), err)
		return nil, errors.New(msg)
	}

	if a.ShowDetail {
		log.Printf(`rawData btc38.MyAccount()=%s`, string(rawData))
	}

	m, err := handleMyAccountRawData(rawData)
	if err != nil {
		msg := fmt.Sprintf("无法转换MyBalance的rawData(%s):%s", string(rawData), err)
		return nil, errors.New(msg)
	}

	if a.ShowDetail {
		log.Printf(`After JSONDecode: btc38.MyAccount()=%v`, m)
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
	err := ec.JSONDecode(rawData, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

var orderTypeMap = map[ec.OrderType]int{
	ec.BUY:  1,
	ec.SELL: 2,
}

//Order 下单交易
func (a *API) Order(t ec.OrderType, money, coin string, price, amount float64) (int64, error) {
	ot := orderTypeMap[t]
	rawData, err := a.orderRawData(ot, money, coin, price, amount)
	if err != nil {
		return 0, err
	}

	if a.ShowDetail {
		log.Printf(`rawData btc38.Order(%s, %s, %s, %f, %f)=%s`, t, money, coin, price, amount, string(rawData))
		log.Println()
	}

	return handleOrderRawData(rawData)
}

func (a *API) orderRawData(ot int, money, coin string, price, amount float64) ([]byte, error) {
	//NOTICE: btc38的价格**可能**只支持6位有效数字。比如现在BTC的价格是18000多，如果买价是12345.6就可以下单，如果买价是12345.67，就会显示deciError1
	priceStr := price2Str(price)
	amountStr := fmt.Sprintf("%.6f", amount)

	body := a.orderBody(ot, money, coin, priceStr, amountStr)

	if a.ShowDetail {
		log.Printf(`Body of btc38.Order(%d, %s, %s, %f, %f)=%v`, ot, money, coin, price, amount, body)
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

	if a.ShowDetail {
		log.Printf("encoded url.Values of btc38.Order(%d, %s, %s, %s, %s)=%s\n", ot, money, coin, price, amount, encoded)
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

//CancelOrder 下单交易
func (a *API) CancelOrder(money, coin string, orderID int64) (bool, error) {
	rawData, err := a.cancelOrderRawData(money, coin, orderID)
	if err != nil {
		return false, err
	}

	if a.ShowDetail {
		log.Printf(`rawData btc38.CancelOrder(%s, %s, %d)=%s`, money, coin, orderID, string(rawData))
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

//MyOrders 我所有还没有成交的挂单
func (a *API) MyOrders(money, coin string) ([]ec.Order, error) {
	rawData, err := a.myOrdersRawData(money, coin)
	if err != nil {
		return nil, err
	}

	if a.ShowDetail {
		log.Printf(`rawData btc38.MyOrders(%s, %s)=%s`, money, coin, string(rawData))
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
func (a *API) handleMyOrdersRawData(rawData []byte, money string) ([]ec.Order, error) {
	if string(rawData) == "no_order" {
		return []ec.Order{}, nil
	}

	resp := []order{}
	err := ec.JSONDecode(rawData, &resp)
	if err != nil {
		msg := fmt.Sprintf("JSONDecode %s时，出错:%s", string(rawData), err)
		return nil, errors.New(msg)
	}

	if a.ShowDetail {
		log.Printf("\nBefore JSONDecode %s\nAfter JSONDecode %v\n", rawData, resp)
	}

	result := make([]ec.Order, len(resp))
	for i, v := range resp {
		r, err := v.normalize(money)
		if err != nil {
			msg := fmt.Sprintf("处理MyOrders的rawData失败，无法normalize:%s", err)
			return nil, errors.New(msg)
		}
		result[i] = *r
	}

	if a.ShowDetail {
		log.Printf("\nAfter normalize %v\n", result)
	}

	return result, nil
}

//MyTransRecords 获取我的交易记录
func (a *API) MyTransRecords(money, coin string, Tid int64) (ec.Trades, error) {
	//TODO: 完成这个方法

	done := false
	//res := ec.Trades{}

	for !done {
		//获取myTrade数据，len(mt)==0 return
		//转换成标准的ec.Trades 格式 tempECTrades

		//找出所有符合条件的tempECTrades2
		//len(tempECTrades2)==0 , return

		//res = append(res, tempECTrades2)
		done = true

		//len(tempECTrades2) < len(tempECTrades), return

	}
	//return

	mtl, err := a.MyTradeList(money, coin, 1)
	if err != nil {
		msg := fmt.Sprintf("MyTradeList(%s, %s)获取失败:%s", money, coin, err)
		return nil, errors.New(msg)
	}

	res := ec.Trades{}
	for _, mt := range mtl {
		t, err := mt.normalize(a.ID)
		if err != nil {
			return nil, err
		}
		res = append(res, t)
	}

	return res, nil
}

//nextPageList
func nextECTrades(a *API, money, coin string) func() (ec.Trades, error) {
	page := 1 //REVIEW: how about page = 0

	return func() (ec.Trades, error) {
		mts, err := a.MyTradeList(money, coin, page)
		if err != nil {
			msg := fmt.Sprintf("nextPageList 无法获取%s.MyTradeList(%s, %s, %d)的数据: %s", a.Name(), money, coin, page, err)
			return nil, errors.New(msg)
		}
		//TODO: 删除此处内容
		fmt.Println(mts)

		page++
		return nil, nil
	}
}

//FIXME: 把MyTransRecord抽象完成。
func (a *API) myTrades2ECTrades(mts []MyTrade) (ec.Trades, error) {
	res := make(ec.Trades, len(mts))
	for i, mt := range mts {
		et, err := mt.normalize(a.ID)
		if err != nil {
			msg := fmt.Sprintf("无法把%s转换成ec.Trade: %s", mt, err)
			return nil, errors.New(msg)
		}
		res[i] = et
	}

	//mts是降序，ec.Trades是升序，所以res要reverse一下
	res.Sort()

	return res, nil
}

//MyTradeList 按照btc38的API的格式，返回交易记录结果。
//TODO: 当page很大的时候，会返回一个空切片，还是报错。
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
	err := ec.JSONDecode(rawData, &resp)
	if err != nil {
		return nil, err
	}

	if a.ShowDetail {
		log.Printf("\n前3条交易记录明细 %s\n", resp[:3])
	}

	return resp, nil
}

func (a *API) md5(time string) string {
	md := fmt.Sprintf("%s_%d_%s_%s", a.PublicKey, a.ID, a.SecretKey, time)
	md5 := ec.MD5([]byte(md))
	return ec.HexEncodeToString(md5)
}
