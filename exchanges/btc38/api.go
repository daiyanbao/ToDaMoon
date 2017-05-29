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
func (b *BTC38) Ticker(coin, money string) (*ec.Ticker, error) {
	rawData, err := b.getTickerRawData(coin, money)
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
func (b *BTC38) allTicker(money string) (map[string]*ec.Ticker, error) {
	rawData, err := b.getTickerRawData("all", money)
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

// Ticker returns okcoin's latest ticker data
func (b *BTC38) getTickerRawData(coin, money string) ([]byte, error) {
	path := tickerURLMaker(coin, money)
	return b.Get(path)
}

func tickerURLMaker(coin, money string) string {
	return urlMaker(tickerURL, coin, money)
}

func urlMaker(URL string, coin, money string) string {
	v := url.Values{}
	v.Set("c", coin)
	v.Set("mk_type", money)

	return ec.Path(URL, v)
}

//Depth 是反馈市场深度信息
func (b *BTC38) Depth(coin, money string) (*ec.Depth, error) {
	rawData, err := b.getDepthRawData(coin, money)
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

// Ticker returns okcoin's latest ticker data
func (b *BTC38) getDepthRawData(coin, money string) ([]byte, error) {
	path := depthURLMaker(coin, money)
	return b.Get(path)
}

func depthURLMaker(coin, money string) string {
	return urlMaker(depthURL, coin, money)
}

//Trades 返回市场的交易记录
//当tid<=0时，返回最新的30条记录
func (b *BTC38) Trades(coin, money string, tid int64) (ec.Trades, error) {
	rawData, err := b.getTradesRawData(coin, money, tid)
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

func (b *BTC38) getTradesRawData(coin, money string, tid int64) ([]byte, error) {
	path := tradesURLMaker(coin, money, tid)
	return b.Get(path)
}

func tradesURLMaker(coin, money string, tid int64) string {
	path := urlMaker(tradesURL, coin, money)
	if tid <= 0 {
		return path
	}
	postfix := fmt.Sprintf("&tid=%d", tid)
	return path + postfix
}

//Balance 返回市场的交易记录
//TODO: 把返回的数据修改成ec.Balance
func (b *BTC38) Balance() (MyBalance, error) {
	rawData, err := b.getMyBalanceRawData()
	if err != nil {
		return nil, err
	}

	resp := MyBalance{}
	err = ec.JSONDecode(rawData, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (b *BTC38) getMyBalanceRawData() ([]byte, error) {
	body := b.myBalanceBodyMaker()
	return b.Post(getMyBalanceURL, body)
}

func (b *BTC38) myBalanceBodyMaker() io.Reader {
	v := url.Values{}
	v.Set("key", b.PublicKey)
	nowTime := fmt.Sprint(time.Now().Unix())
	v.Set("time", nowTime)
	md5 := b.md5(nowTime)
	v.Set("md5", md5)

	encoded := v.Encode()
	return strings.NewReader(encoded)
}

func (b *BTC38) md5(time string) string {
	md := fmt.Sprintf("%s_%d_%s_%s", b.PublicKey, b.ID, b.SecretKey, time)
	md5 := ec.MD5([]byte(md))
	return ec.HexEncodeToString(md5)
}

type orderType int

const (
	//BUY 是使用money换coin的过程
	BUY orderType = 1 //不用iota是因为btc38的api指定了数字1为买入
	//SELL 是使用coin换money的过程
	SELL orderType = 2
)

//Trade 下单交易
//TODO: 把money改成枚举类型，所有的
func (b *BTC38) Trade(ot orderType, coin, money string, price, amount float64) (int, error) {
	rawData, err := b.getTradeRawData(ot, coin, money, price, amount)
	if err != nil {
		return 0, err
	}
	fmt.Println(string(rawData)) //TODO: 删除此处内容
	return handleTradeRawData(rawData)
}

func (b *BTC38) getTradeRawData(ot orderType, coin, money string, price, amount float64) ([]byte, error) {
	body := b.tradeBodyMaker(ot, coin, money, price, amount)
	return b.Post(submitOrderURL, body)
}

func (b *BTC38) tradeBodyMaker(ot orderType, coin, money string, price, amount float64) io.Reader {
	v := url.Values{}
	v.Set("key", b.PublicKey)
	nowTime := fmt.Sprint(time.Now().Unix())
	v.Set("time", nowTime)
	md5 := b.md5(nowTime)
	v.Set("md5", md5)

	v.Set("type", fmt.Sprint(ot))
	v.Set("coinname", coin)
	v.Set("mk_type", money)
	v.Set("price", strconv.FormatFloat(price, 'f', -1, 64))
	v.Set("amount", strconv.FormatFloat(amount, 'f', -1, 64))
	encoded := v.Encode()

	fmt.Println("order body:", encoded) //TODO: 删除此处内容
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
//TODO: 把money改成枚举类型，所有的
func (b *BTC38) CancelOrder(coin, money string, orderID int) (bool, error) {
	rawData, err := b.getCancelOrderRawData(coin, money, orderID)
	if err != nil {
		return false, err
	}

	return handleCancelOrderRawData(rawData)
}

func (b *BTC38) getCancelOrderRawData(coin, money string, orderID int) ([]byte, error) {
	body := b.cancelOrderBodyMaker(coin, money, orderID)
	return b.Post(cancelOrderURL, body)
}

func (b *BTC38) cancelOrderBodyMaker(coin, money string, orderID int) io.Reader {
	v := url.Values{}
	v.Set("key", b.PublicKey)
	nowTime := fmt.Sprint(time.Now().Unix())
	v.Set("time", nowTime)
	md5 := b.md5(nowTime)
	v.Set("md5", md5)

	v.Set("mk_type", money)
	v.Set("order_id", strconv.Itoa(orderID))
	v.Set("coinname", coin)
	encoded := v.Encode()

	fmt.Println("order body:", encoded) //TODO: 删除此处内容
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
	ID        int     `json:"order_id,string"`
	OrderType string  `json:"order_type,string"`
	Coin      string  `json:"order_coinname,string"`
	Amount    float64 `json:"order_amount,string"`
	Price     float64 `json:"order_price,string"`
}

//getMyOrders 下单交易
//TODO: 把money改成枚举类型，所有的
func (b *BTC38) getMyOrders(coin, money string) ([]order, error) {
	rawData, err := b.getMyOrdersRawData(coin, money)
	if err != nil {
		return nil, err
	}

	return handleMyOrdersRawData(rawData)
}

func (b *BTC38) getMyOrdersRawData(coin, money string) ([]byte, error) {
	body := b.myOrdersBodyMaker(coin, money)
	return b.Post(getOrderListURL, body)
}

func (b *BTC38) myOrdersBodyMaker(coin, money string) io.Reader {
	v := url.Values{}
	v.Set("key", b.PublicKey)
	nowTime := fmt.Sprint(time.Now().Unix())
	v.Set("time", nowTime)
	md5 := b.md5(nowTime)
	v.Set("md5", md5)

	v.Set("mk_type", money)
	v.Set("coinname", coin)
	encoded := v.Encode()

	fmt.Println("orders body:", encoded) //TODO: 删除此处内容
	return strings.NewReader(encoded)
}

func handleMyOrdersRawData(rawData []byte) ([]order, error) {
	resp := []order{}

	//TODO: 删除此处内容
	fmt.Println(string(rawData))

	err := ec.JSONDecode(rawData, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
