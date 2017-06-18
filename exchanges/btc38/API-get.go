package btc38

import (
	ec "ToDaMoon/exchanges"
	"fmt"
	"log"
	"net/url"
)

//Ticker 可以返回coin的ticker信息
func (a *API) Ticker(money, coin string) (*ec.Ticker, error) {
	rawData, err := a.getTickerRawData(money, coin)
	if err != nil {
		return nil, err
	}

	if a.ShowDetail {
		log.Printf(`btc38.Ticker("%s","%s")=%s`, money, coin, string(rawData))
	}

	resp := TickerResponse{}
	err = ec.JSONDecode(rawData, &resp)
	if err != nil {
		return nil, err
	}

	if a.ShowDetail {
		log.Printf(`After JSONDecode: btc38.Ticker("%s","%s")=%v`, money, coin, resp)
	}

	return resp.Ticker.normalize(), nil
}

//AllTicker 返回money市场中全部coin的ticker
func (a *API) AllTicker(money string) (map[string]*ec.Ticker, error) {
	rawData, err := a.getTickerRawData("all", money)
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

func (a *API) getTickerRawData(money, coin string) ([]byte, error) {
	path := tickerURLMaker(money, coin)
	return a.Get(path)
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
func (a *API) Depth(money, coin string) (*ec.Depth, error) {
	rawData, err := a.getDepthRawData(money, coin)
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

func (a *API) getDepthRawData(money, coin string) ([]byte, error) {
	path := depthURLMaker(money, coin)
	return a.Get(path)
}

func depthURLMaker(money, coin string) string {
	return urlMaker(depthURL, money, coin)
}

//TransRecords 返回市场的交易记录
//当tid<=0时，返回最新的30条记录
//TODO: 修改相关子函数的名称
func (a *API) TransRecords(money, coin string, tid int64) (ec.Trades, error) {
	rawData, err := a.getTradesRawData(money, coin, tid)
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

func (a *API) getTradesRawData(money, coin string, tid int64) ([]byte, error) {
	path := tradesURLMaker(money, coin, tid)
	return a.Get(path)
}

func tradesURLMaker(money, coin string, tid int64) string {
	path := urlMaker(tradesURL, money, coin)
	if tid <= 0 {
		return path
	}
	postfix := fmt.Sprintf("&tid=%d", tid)
	return path + postfix
}
