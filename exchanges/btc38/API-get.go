package btc38

import (
	ec "ToDaMoon/exchanges"
	"fmt"
	"log"
	"net/url"
)

//Ticker 可以返回coin的ticker信息
func (a *API) Ticker(money, coin string) (*ec.Ticker, error) {
	rawData, err := a.tickerRawData(money, coin)
	if err != nil {
		return nil, err
	}

	if a.ShowDetail {
		log.Printf(`rawData btc38.Ticker("%s","%s")=%s`, money, coin, string(rawData))
	}

	resp := tickerResp{}
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
	rawData, err := a.tickerRawData(money, "all")
	if err != nil {
		return nil, err
	}

	if a.ShowDetail {
		log.Printf(`rawData btc38.AllTicker("%s")=%s`, money, string(rawData))
	}

	resp := make(map[string]tickerResp)
	err = ec.JSONDecode(rawData, &resp)
	if err != nil {
		return nil, err
	}

	if a.ShowDetail {
		log.Printf(`After JSONDecode: btc38.AllTicker("%s")=%v`, money, resp)
	}

	result := make(map[string]*ec.Ticker)
	for k, v := range resp {
		result[k] = v.Ticker.normalize()
	}
	return result, nil
}

func (a *API) tickerRawData(money, coin string) ([]byte, error) {
	path := urlMaker(tickerURL, money, coin)
	return a.Get(path)
}

//Depth 是反馈市场深度信息
func (a *API) Depth(money, coin string) (*ec.Depth, error) {
	rawData, err := a.depthRawData(money, coin)
	if err != nil {
		return nil, err
	}

	if a.ShowDetail {
		log.Printf(`rawData btc38.Depth("%s","%s")=%s`, money, coin, string(rawData))
	}

	resp := ec.Depth{}
	err = ec.JSONDecode(rawData, &resp)
	if err != nil {
		return nil, err
	}

	if a.ShowDetail {
		log.Printf(`After JSONDecode: btc38.Depth("%s","%s")=%v`, money, coin, resp)
	}

	return &resp, nil
}

func (a *API) depthRawData(money, coin string) ([]byte, error) {
	path := urlMaker(depthURL, money, coin)
	return a.Get(path)
}

//TransRecords 返回市场的交易记录
//当tid<=0时，返回最新的30条记录
func (a *API) TransRecords(money, coin string, tid int64) (ec.Trades, error) {
	rawData, err := a.transRecordsRawData(money, coin, tid)
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

func (a *API) transRecordsRawData(money, coin string, tid int64) ([]byte, error) {
	path := urlMaker(transRecordsURL, money, coin)

	if tid > 0 {
		postfix := fmt.Sprintf("&tid=%d", tid)
		path += postfix
	}

	return a.Get(path)
}

func urlMaker(URL string, money, coin string) string {
	v := url.Values{}

	v.Set("c", coin)
	v.Set("mk_type", money)

	return ec.Path(URL, v)
}
