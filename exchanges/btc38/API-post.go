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

var transTypeMap = map[ec.OrderType]int{
	ec.BUY:  1,
	ec.SELL: 2,
}

//Order 下单交易
//TODO: 修改order的子函数的名称
func (a *API) Order(t ec.OrderType, money, coin string, price, amount float64) (int64, error) {
	ot := transTypeMap[t]
	//TODO: 修改trade为trans
	rawData, err := a.getTradeRawData(ot, money, coin, price, amount)
	if err != nil {
		return 0, err
	}

	return handleTradeRawData(rawData)
}

func (a *API) getTradeRawData(ot int, money, coin string, price, amount float64) ([]byte, error) {
	body := a.tradeBodyMaker(ot, money, coin, price, amount)
	return a.Post(submitOrderURL, body)
}

func (a *API) tradeBodyMaker(ot int, money, coin string, price, amount float64) io.Reader {
	v := url.Values{}
	v.Set("key", a.PublicKey)
	nowTime := fmt.Sprint(time.Now().Unix())
	v.Set("time", nowTime)
	md5 := a.md5(nowTime)
	v.Set("md5", md5)

	v.Set("type", fmt.Sprint(ot))
	v.Set("coinname", coin)
	v.Set("mk_type", money)
	v.Set("price", strconv.FormatFloat(price, 'f', -1, 64))
	v.Set("amount", strconv.FormatFloat(amount, 'f', -1, 64))
	encoded := v.Encode()

	return strings.NewReader(encoded)
}

func handleTradeRawData(rawData []byte) (int64, error) {
	r := string(rawData)

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
	rawData, err := a.getCancelOrderRawData(money, coin, orderID)
	if err != nil {
		return false, err
	}

	return handleCancelOrderRawData(rawData)
}

func (a *API) getCancelOrderRawData(money, coin string, orderID int64) ([]byte, error) {
	body := a.cancelOrderBodyMaker(money, coin, orderID)
	return a.Post(cancelOrderURL, body)
}

func (a *API) cancelOrderBodyMaker(money, coin string, orderID int64) io.Reader {
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

func handleCancelOrderRawData(rawData []byte) (bool, error) {
	r := string(rawData)

	if r == "succ" {
		return true, nil
	}
	//TODO: 最好能够给出cancel时候，遇到的各种情况，比如，订单已经成交等等。
	return false, errors.New(r)
}

//TODO: 取消orderType
type orderType int

const (
	//BUY 是使用money换coin的过程
	BUY orderType = 1 //不用iota是因为btc38的api指定了数字1为买入
	//SELL 是使用coin换money的过程
	SELL orderType = 2
)

//Order 是订单
//TODO: 转换成exchanges的标准形式
type Order struct {
	ID   int    `json:"id,string"`
	Coin string `json:"coinname"`
	//TODO: 取消orderType
	OrderType orderType `json:"type,string"`
	Amount    float64   `json:"amount,string"`
	Price     float64   `json:"price,string"`
	Time      string    `json:"time"`
}

//MyOrders 获取我所有的挂单
func (a *API) MyOrders(money, coin string) ([]ec.Order, error) {
	rawData, err := a.getMyOrdersRawData(money, coin)
	if err != nil {
		return nil, err
	}

	return handleMyOrdersRawData(rawData)
}

func (a *API) getMyOrdersRawData(money, coin string) ([]byte, error) {
	body := a.myOrdersBodyMaker(money, coin)
	return a.Post(getOrderListURL, body)
}

func (a *API) myOrdersBodyMaker(money, coin string) io.Reader {
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

func handleMyOrdersRawData(rawData []byte) ([]ec.Order, error) {
	resp := []ec.Order{}

	fmt.Println("ec.Order还是空的")

	err := ec.JSONDecode(rawData, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

//MyTrade 是我的交易记录
//TODO: 转换成exchanges的标准结果
type MyTrade struct {
	ID       int     `json:"id,string"`
	BuyerID  string  `json:"buyer_id"`
	SellerID string  `json:"seller_id"`
	Volume   float64 `json:"volume,string"`
	Price    float64 `json:"price,string"`
	Coin     string  `json:"coinname"`
	Time     string  `json:"time"`
}

//MyTransRecords 获取我的交易记录
//TODO: 修改子函数的名称
//FIXME: 这个函数的方法，还没有统一。
func (a *API) MyTransRecords(money, coin string, page int64) (ec.Trades, error) {
	rawData, err := a.getMyTradesRawData(money, coin, page)
	if err != nil {
		return nil, err
	}

	return handleMyTradesRawData(rawData)
}

func (a *API) getMyTradesRawData(money, coin string, page int64) ([]byte, error) {
	body := a.myTradesBodyMaker(money, coin, page)
	return a.Post(getMyTradeListURL, body)
}

func (a *API) myTradesBodyMaker(money, coin string, page int64) io.Reader {
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

func handleMyTradesRawData(rawData []byte) (ec.Trades, error) {
	resp := ec.Trades{}
	//TODO: 这个函数是没有完成的。
	fmt.Println("这个函数是没有完成的")
	err := ec.JSONDecode(rawData, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (a *API) md5(time string) string {
	md := fmt.Sprintf("%s_%d_%s_%s", a.PublicKey, a.ID, a.SecretKey, time)
	md5 := ec.MD5([]byte(md))
	return ec.HexEncodeToString(md5)
}
