// Package btc38 wrap btc38.com api in singleton pattern
// 此模块是单例模式，请使用btc38.Instance()来生成实例
package btc38

import (
	"net/url"
	"strconv"

	ec "ToDaMoon/exchanges"
)

const (
	apiURL    = "http://api.btc38.com/v1/"
	tickerURL = "ticker.php"
	balance   = "getMyBalance.php"
	hsitory   = "trades.php"
	/*
		TICKER            = "ticker.do"
		DEPTH             = "depth.do"
		TRADES            = "trades.do"
		KLINE             = "kline.do"

		TRADE             = "trade.do"
		HISTORY           = "trade_history.do"
		TRADE_BATCH       = "batch_trade.do"
		CANCEL_ORDER      = "cancel_order.do"
		ORDER_INFO        = "order_info.do"
		ORDERS_INFO       = "orders_info.do"
		ORDER_HISTORY     = "order_history.do"
		WITHDRAW          = "withdraw.do"
		CANCEL_WITHDRAW   = "cancel_withdraw.do"
		WITHDRAW_INFO     = "withdraw_info.do"
		ORDER_FEE         = "order_fee.do"
		LEND_DEPTH        = "lend_depth.do"
		BORROWS_INFO      = "borrows_info.do"
		BORROW_MONEY      = "borrow_money.do"
		BORROW_CANCEL     = "cancel_borrow.do"
		BORROW_ORDER_INFO = "borrow_order_info.do"
		REPAYMENT         = "repayment.do"
		UNREPAYMENTS_INFO = "unrepayments_info.do"
		ACCOUNT_RECORDS   = "account_records.do"
	*/
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

//AllCoins returns okcoin's latest ticker data
func (o *BTC38) AllCoins() ([]string, error) {
	resp := make(map[string]interface{})

	v := url.Values{}
	v.Set("c", "all")
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
	result := make([]string, len(resp))
	i := 0
	for k := range resp {
		result[i] = k
		i++
	}

	return result, nil
}

/*
以下是具体的post方法。
*/

//Balance 返回账户信息
func (o *BTC38) Balance() (MyBalance, error) {
	result := new(MyBalance)
	err := o.post(balance, url.Values{}, result)

	if err != nil {
		return nil, err
	}

	// if result.ErrorCode > 0 {
	// 	s := fmt.Sprintln("获取用户信息出错:", o.restErrors[result.ErrorCode])
	// 	return nil, errors.New(s)
	// }

	return *result, nil
}

// UserInfo give user's information
func (o *BTC38) UserInfo() (*Balance, error) {
	result := new(Balance)
	err := o.post(balance, url.Values{}, result)

	if err != nil {
		return nil, err
	}

	// if result.ErrorCode > 0 {
	// 	s := fmt.Sprintln("获取用户信息出错:", o.restErrors[result.ErrorCode])
	// 	return nil, errors.New(s)
	// }

	return result, nil
}

// TradeHistory returns trade history , no personal.
func (o *BTC38) TradeHistory(symbol string, TradeID int64) (ec.Trades, error) {
	result := []Trade{}
	v := url.Values{}
	v.Set("symbol", symbol+"_cny")
	v.Set("since", strconv.FormatInt(TradeID, 10))

	err := o.post(hsitory, v, &result)

	if err != nil {
		return nil, err
	}

	mt := ec.Trades{}
	for _, v := range result {
		mt = append(mt, v.Normalize())
	}

	return mt, nil
}
