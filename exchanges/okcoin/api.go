package okcoin

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

const (
	baseURL           = "http://api.btc38.com/v1/"
	tickerURL         = baseURL + "ticker.php"
	depthURL          = baseURL + "depth.php"
	tradesURL         = baseURL + "trades.php"
	getMyBalanceURL   = baseURL + "getMyBalance.php"
	submitOrderURL    = baseURL + "submitOrder.php"
	cancelOrderURL    = baseURL + "cancelOrder.php"
	getOrderListURL   = baseURL + "getOrderList.php"
	getMyTradeListURL = baseURL + "getMyTradeList.php"
)

//Ticker 可以返回coin的ticker信息
func (o *OKCoin) Ticker(money, coin string) (*ec.Ticker, error) {
	rawData, err := o.getTickerRawData(money, coin)
	if err != nil {
		return nil, err
	}

	resp := TickerResponse{}
	err = ec.JSONDecode(rawData, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Ticker.normalize(), nil
}

//AllTicker 返回money市场中全部coin的ticker
func (o *OKCoin) allTicker(money string) (map[string]*ec.Ticker, error) {
	rawData, err := o.getTickerRawData("all", money)
	if err != nil {
		return nil, err
	}

	resp := make(map[string]TickerResponse)
	err = ec.JSONDecode(rawData, &resp)
	if err != nil {
		return nil, err
	}

	result := make(map[string]*ec.Ticker)
	for k, v := range resp {
		result[k] = v.Ticker.normalize()
	}
	return result, nil
}

func (o *OKCoin) getTickerRawData(money, coin string) ([]byte, error) {
	path := tickerURLMaker(money, coin)
	return o.Get(path)
}

func tickerURLMaker(money, coin string) string {
	return urlMaker(tickerURL, money, coin)
}

func urlMaker(URL string, money, coin string) string {
	v := url.Values{}
	v.Set("c", coin)
	v.Set("mk_type", money)

	return ec.Path(URL, v)
}

//Depth 是反馈市场深度信息
func (o *OKCoin) Depth(money, coin string) (*ec.Depth, error) {
	rawData, err := o.getDepthRawData(money, coin)
	if err != nil {
		return nil, err
	}

	resp := ec.Depth{}
	err = ec.JSONDecode(rawData, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (o *OKCoin) getDepthRawData(money, coin string) ([]byte, error) {
	path := depthURLMaker(money, coin)
	return o.Get(path)
}

func depthURLMaker(money, coin string) string {
	return urlMaker(depthURL, money, coin)
}

//Trades 返回市场的交易记录
//当tid<=0时，返回最新的30条记录
func (o *OKCoin) Trades(money, coin string, tid int64) (ec.Trades, error) {
	rawData, err := o.getTradesRawData(money, coin, tid)
	if err != nil {
		return nil, err
	}

	resp := ec.Trades{}
	err = ec.JSONDecode(rawData, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (o *OKCoin) getTradesRawData(money, coin string, tid int64) ([]byte, error) {
	path := tradesURLMaker(money, coin, tid)
	return o.Get(path)
}

func tradesURLMaker(money, coin string, tid int64) string {
	path := urlMaker(tradesURL, money, coin)
	if tid <= 0 {
		return path
	}
	postfix := fmt.Sprintf("&tid=%d", tid)
	return path + postfix
}

//Balance 返回市场的交易记录
//TODO: 把返回的数据修改成ec.Balance
func (o *OKCoin) Balance() (MyBalance, error) {
	rawData, err := o.getMyBalanceRawData()
	if err != nil {
		return nil, err
	}

	return handleMyBalanceRawData(rawData)
}

func (o *OKCoin) getMyBalanceRawData() ([]byte, error) {
	body := o.myBalanceBodyMaker()
	return o.Post(getMyBalanceURL, body)
}

func (o *OKCoin) myBalanceBodyMaker() io.Reader {
	v := url.Values{}
	v.Set("key", o.PublicKey)
	nowTime := fmt.Sprint(time.Now().Unix())
	v.Set("time", nowTime)
	md5 := o.md5(nowTime)
	v.Set("md5", md5)

	encoded := v.Encode()
	return strings.NewReader(encoded)
}

func (o *OKCoin) md5(time string) string {
	md := fmt.Sprintf("%s_%s_%s", o.PublicKey, o.SecretKey, time)
	md5 := ec.MD5([]byte(md))
	return ec.HexEncodeToString(md5)
}

func handleMyBalanceRawData(rawData []byte) (MyBalance, error) {
	resp := MyBalance{}
	err := ec.JSONDecode(rawData, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

type orderType int

const (
	//BUY 是使用money换coin的过程
	BUY orderType = 1 //不用iota是因为btc38的api指定了数字1为买入
	//SELL 是使用coin换money的过程
	SELL orderType = 2
)

//Trade 下单交易
func (o *OKCoin) Trade(ot orderType, money, coin string, price, amount float64) (int, error) {
	rawData, err := o.getTradeRawData(ot, money, coin, price, amount)
	if err != nil {
		return 0, err
	}

	return handleTradeRawData(rawData)
}

func (o *OKCoin) getTradeRawData(ot orderType, money, coin string, price, amount float64) ([]byte, error) {
	body := o.tradeBodyMaker(ot, money, coin, price, amount)
	return o.Post(submitOrderURL, body)
}

func (o *OKCoin) tradeBodyMaker(ot orderType, money, coin string, price, amount float64) io.Reader {
	v := url.Values{}
	v.Set("key", o.PublicKey)
	nowTime := fmt.Sprint(time.Now().Unix())
	v.Set("time", nowTime)
	md5 := o.md5(nowTime)
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
func (o *OKCoin) CancelOrder(money, coin string, orderID int) (bool, error) {
	rawData, err := o.getCancelOrderRawData(money, coin, orderID)
	if err != nil {
		return false, err
	}

	return handleCancelOrderRawData(rawData)
}

func (o *OKCoin) getCancelOrderRawData(money, coin string, orderID int) ([]byte, error) {
	body := o.cancelOrderBodyMaker(money, coin, orderID)
	return o.Post(cancelOrderURL, body)
}

func (o *OKCoin) cancelOrderBodyMaker(money, coin string, orderID int) io.Reader {
	v := url.Values{}
	v.Set("key", o.PublicKey)
	nowTime := fmt.Sprint(time.Now().Unix())
	v.Set("time", nowTime)
	md5 := o.md5(nowTime)
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

type order struct {
	ID        int       `json:"id,string"`
	Coin      string    `json:"coinname"`
	OrderType orderType `json:"type,string"`
	Amount    float64   `json:"amount,string"`
	Price     float64   `json:"price,string"`
	Time      string    `json:"time"`
}

//getMyOrders 下单交易
func (o *OKCoin) getMyOrders(money, coin string) ([]order, error) {
	rawData, err := o.getMyOrdersRawData(money, coin)
	if err != nil {
		return nil, err
	}

	return handleMyOrdersRawData(rawData)
}

func (o *OKCoin) getMyOrdersRawData(money, coin string) ([]byte, error) {
	body := o.myOrdersBodyMaker(money, coin)
	return o.Post(getOrderListURL, body)
}

func (o *OKCoin) myOrdersBodyMaker(money, coin string) io.Reader {
	v := url.Values{}
	v.Set("key", o.PublicKey)
	nowTime := fmt.Sprint(time.Now().Unix())
	v.Set("time", nowTime)
	md5 := o.md5(nowTime)
	v.Set("md5", md5)

	v.Set("mk_type", money)
	v.Set("coinname", coin)
	encoded := v.Encode()

	return strings.NewReader(encoded)
}

func handleMyOrdersRawData(rawData []byte) ([]order, error) {
	resp := []order{}

	err := ec.JSONDecode(rawData, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

type myTrade struct {
	ID       int     `json:"id,string"`
	BuyerID  string  `json:"buyer_id"`
	SellerID string  `json:"seller_id"`
	Volume   float64 `json:"volume,string"`
	Price    float64 `json:"price,string"`
	Coin     string  `json:"coinname"`
	Time     string  `json:"time"`
}

//getMyTrades 下单交易
func (o *OKCoin) getMyTrades(money, coin string, page int) ([]myTrade, error) {
	rawData, err := o.getMyTradesRawData(money, coin, page)
	if err != nil {
		return nil, err
	}

	return handleMyTradesRawData(rawData)
}

func (o *OKCoin) getMyTradesRawData(money, coin string, page int) ([]byte, error) {
	body := o.myTradesBodyMaker(money, coin, page)
	return o.Post(getMyTradeListURL, body)
}

func (o *OKCoin) myTradesBodyMaker(money, coin string, page int) io.Reader {
	v := url.Values{}
	v.Set("key", o.PublicKey)
	nowTime := fmt.Sprint(time.Now().Unix())
	v.Set("time", nowTime)
	md5 := o.md5(nowTime)
	v.Set("md5", md5)

	v.Set("mk_type", money)
	v.Set("coinname", coin)
	v.Set("page", strconv.Itoa(page))
	encoded := v.Encode()

	return strings.NewReader(encoded)
}

func handleMyTradesRawData(rawData []byte) ([]myTrade, error) {
	resp := []myTrade{}

	err := ec.JSONDecode(rawData, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
