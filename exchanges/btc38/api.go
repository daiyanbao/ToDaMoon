package btc38

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

// Ticker returns okcoin's latest ticker data
func (o *BTC38) Ticker(symbol string) (*ec.Ticker, error) {
	resp := TickerResponse{}

	v := url.Values{}
	v.Set("c", symbol)
	v.Set("mk_type", "cny")

	ansChan := make(chan answer)
	o.ask <- ask{Type: get,
		Path:       ec.Path(apiURL+tickerURL, v),
		AnswerChan: ansChan}
	ans := <-ansChan

	if ans.err != nil {
		return nil, ans.err
	}

	err := ec.JSONDecode(ans.body, &resp)
	if err != nil {
		return nil, err
	}
	t := resp.Ticker

	return &ec.Ticker{
		Buy:  t.Buy,
		High: t.High,
		Last: t.Last,
		Low:  t.Low,
		Sell: t.Sell,
		Vol:  t.Vol,
	}, nil
}