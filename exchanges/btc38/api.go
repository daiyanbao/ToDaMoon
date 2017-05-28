package btc38

import (
	ec "ToDaMoon/exchanges"
	"net/url"
)

const (
	baseURL           = "http://api.btc38.com/v1/"
	tickerURL         = baseURL + "ticker.php"
	depthURL          = baseURL + "depth.php"
	tradesURL         = baseURL + "trades.php"
	getBalanceURL     = baseURL + "getMyBalance.php"
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
