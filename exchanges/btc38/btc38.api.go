// Package btc38 wrap btc38.com api in singleton pattern 
// 此模块是单例模式，请使用btc38.Instance()来生成实例
package btc38
import (
	"errors"
	"fmt"
	"net/url"
	"strconv"

	ec "ToDaMoon/exchanges"
)

const (
	apiURL    = "https://www.okcoin.cn/api/v1/"
	tickerURL = "ticker.do"
	userinfo  = "userinfo.do"
	hsitory   = "trade_history.do"
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
func (o *OKCoin) Ticker(symbol string) (*ec.Ticker, error) {
	resp := TickerResponse{}

	v := url.Values{}
	v.Set("symbol", symbol+"_cny")

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

/*
以下是具体的post方法。
*/

// UserInfo give user's information
func (o *OKCoin) UserInfo() (*UserInfo, error) {
	result := new(UserInfo)
	err := o.post(userinfo, url.Values{}, result)

	if err != nil {
		return nil, err
	}

	if result.ErrorCode > 0 {
		s := fmt.Sprintln("获取用户信息出错:", o.restErrors[result.ErrorCode])
		return nil, errors.New(s)
	}

	return result, nil
}

// TradeHistory returns trade history , no personal.
func (o *OKCoin) TradeHistory(symbol string, TradeID int64) (ec.Trades, error) {
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
