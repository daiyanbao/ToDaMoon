package btc38

import (
	ec "ToDaMoon/exchanges"
	"fmt"
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
	rawData, err := b.ticker(coin, money)
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
	rawData, err := b.ticker("all", money)
	if err != nil {
		return nil, err
	}

	resp := make(map[string]TickerResponse)
	fmt.Println(string(rawData))
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
func (b *BTC38) ticker(coin, money string) ([]byte, error) {
	v := url.Values{}
	v.Set("c", coin)
	v.Set("mk_type", money)

	ansChan := make(chan ec.Answer)
	b.Ask <- ec.Ask{Type: ec.Get,
		Path:       ec.Path(tickerURL, v),
		AnswerChan: ansChan,
	}

	ans := <-ansChan
	if ans.Err != nil {
		return nil, ans.Err
	}

	return ans.Body, nil
}
