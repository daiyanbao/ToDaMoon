package btc38

import (
	ec "ToDaMoon/exchanges"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func (a *API) md5(time string) string {
	md := fmt.Sprintf("%s_%d_%s_%s", a.PublicKey, a.ID, a.SecretKey, time)
	md5 := ec.MD5([]byte(md))
	return ec.HexEncodeToString(md5)
}

//Account 返回BTC38的账户信息
func (a *API) Account() (*ec.Account, error) {
	rawData, err := a.getMyBalanceRawData()
	if err != nil {
		msg := fmt.Sprintf("无法获取%s的MyBalance的RawDate:%s", a.Name, err)
		return nil, errors.New(msg)
	}

	m, err := handleMyBalanceRawData(rawData)
	if err != nil {
		msg := fmt.Sprintf("无法转换MyBalance的rawData(%s):%s", string(rawData), err)
		return nil, errors.New(msg)
	}

	return m.normalize(), nil
}

func (a *API) getMyBalanceRawData() ([]byte, error) {
	body := a.myBalanceBodyMaker()
	return a.Post(getMyBalanceURL, body)
}

func (a *API) myBalanceBodyMaker() io.Reader {
	v := url.Values{}
	v.Set("key", a.PublicKey)
	nowTime := fmt.Sprint(time.Now().Unix())
	v.Set("time", nowTime)
	md5 := a.md5(nowTime)
	v.Set("md5", md5)

	encoded := v.Encode()
	return strings.NewReader(encoded)
}

func handleMyBalanceRawData(rawData []byte) (myBalance, error) {
	resp := myBalance{}
	err := ec.JSONDecode(rawData, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

var transTypeMap = map[ec.TransType]int{
	ec.BUY:  1,
	ec.SELL: 2,
}

//Trans 下单交易
func (a *API) Trans(t ec.TransType, money, coin string, price, amount float64) (int, error) {
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

func handleTradeRawData(rawData []byte) (int, error) {
	r := string(rawData)

	if r[:5] == "succ|" {
		orderID, err := strconv.Atoi(r[5:])
		if err != nil {
			return 0, err
		}
		return orderID, nil
	}

	return 0, errors.New(r)
}

//CancelOrder 下单交易
func (a *API) CancelOrder(money, coin string, orderID int) (bool, error) {
	rawData, err := a.getCancelOrderRawData(money, coin, orderID)
	if err != nil {
		return false, err
	}

	return handleCancelOrderRawData(rawData)
}

func (a *API) getCancelOrderRawData(money, coin string, orderID int) ([]byte, error) {
	body := a.cancelOrderBodyMaker(money, coin, orderID)
	return a.Post(cancelOrderURL, body)
}

func (a *API) cancelOrderBodyMaker(money, coin string, orderID int) io.Reader {
	v := url.Values{}
	v.Set("key", a.PublicKey)
	nowTime := fmt.Sprint(time.Now().Unix())
	v.Set("time", nowTime)
	md5 := a.md5(nowTime)
	v.Set("md5", md5)

	v.Set("mk_type", money)
	v.Set("order_id", strconv.Itoa(orderID))
	v.Set("coinname", coin)
	encoded := v.Encode()

	return strings.NewReader(encoded)
}

func handleCancelOrderRawData(rawData []byte) (bool, error) {
	r := string(rawData)

	if r == "succ" {
		return true, nil
	}

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
func (a *API) MyOrders(money, coin string) ([]Order, error) {
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

func handleMyOrdersRawData(rawData []byte) ([]Order, error) {
	resp := []Order{}

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

//MyTrades 获取我的交易记录
func (a *API) MyTrades(money, coin string, page int) ([]MyTrade, error) {
	rawData, err := a.getMyTradesRawData(money, coin, page)
	if err != nil {
		return nil, err
	}

	return handleMyTradesRawData(rawData)
}

func (a *API) getMyTradesRawData(money, coin string, page int) ([]byte, error) {
	body := a.myTradesBodyMaker(money, coin, page)
	return a.Post(getMyTradeListURL, body)
}

func (a *API) myTradesBodyMaker(money, coin string, page int) io.Reader {
	v := url.Values{}
	v.Set("key", a.PublicKey)
	nowTime := fmt.Sprint(time.Now().Unix())
	v.Set("time", nowTime)
	md5 := a.md5(nowTime)
	v.Set("md5", md5)

	v.Set("mk_type", money)
	v.Set("coinname", coin)
	v.Set("page", strconv.Itoa(page))
	encoded := v.Encode()

	return strings.NewReader(encoded)
}

func handleMyTradesRawData(rawData []byte) ([]MyTrade, error) {
	resp := []MyTrade{}

	err := ec.JSONDecode(rawData, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
