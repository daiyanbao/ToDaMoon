package btc38

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"

	"github.com/aQuaYi/GoKit"

	ec "github.com/aQuaYi/exchanges"
)

// Ticker 可以返回coin的ticker信息
func (a *API) Ticker(money, coin string) (*ec.Ticker, error) {
	rawData, err := a.tickerRawData(money, coin)
	if err != nil {
		return nil, err
	}

	if a.IsLog {
		log.Printf(`rawData %s.Ticker("%s","%s")=%s`, a.Name(), money, coin, string(rawData))
	}

	resp := tickerResp{}
	err = json.Unmarshal(rawData, &resp)
	if err != nil {
		return nil, err
	}

	if a.IsLog {
		log.Printf(`After json.Unmarshal: %s.Ticker("%s","%s")=%v`, a.Name(), money, coin, resp)
	}

	return resp.Ticker.normalize(), nil
}

// AllTicker 返回money市场中全部coin的ticker
// btc38.com 独有的API
func (a *API) AllTicker(money string) (map[string]*ec.Ticker, error) {
	rawData, err := a.tickerRawData(money, "all")
	if err != nil {
		return nil, GoKit.Err(err, "无法获取%s.AllTicker(%s)的rawData", a.Name(), money)
	}

	if a.IsLog {
		log.Printf(`rawData %s.AllTicker("%s")=%s`, a.Name(), money, string(rawData))
	}

	resp, err := handleAllTickerRawData(rawData)
	if err != nil {
		return nil, GoKit.Err(err, "无法unmarshal %s.AllTicker(%s)的rawData到映射", a.Name(), money)
	}

	if a.IsLog {
		log.Printf(`After json.Unmarshal: %s.AllTicker("%s")=`, a.Name(), money)
		for k, v := range resp {
			log.Println(k, v)
		}
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

func handleAllTickerRawData(rawData []byte) (map[string]tickerResp, error) {
	rawMap := make(map[string]*json.RawMessage)
	err := json.Unmarshal(rawData, &rawMap)
	if err != nil {
		return nil, GoKit.Err(err, "获取rawMap时出现错误")
	}

	/*
		分两次unmarshal是因为，btc38网站返回的json中，有的币的内容可能是{"ticker":""}，这样会导致Unmarshal失败
	*/
	resp := make(map[string]tickerResp)
	for coin, data := range rawMap {
		t := tickerResp{}
		if err := json.Unmarshal(*data, &t); err != nil {
			if err.Error() == "json: cannot unmarshal string into Go struct field tickerResp.ticker of type btc38.ticker" {
				continue
			}
			return nil, GoKit.Err(err, "获取resp时出现错误")
		}
		resp[coin] = t
	}

	return resp, nil
}

// CheckAllTicker 测试btc38.AllTicker()
func CheckAllTicker(a *API, money string) (result string) {
	method := fmt.Sprintf(`%s.AllTicker("%s")`, a.Name(), money)
	fmt.Printf("\n==测试%s==\n", method)

	at, err := a.AllTicker(money)

	if err != nil {
		result = fmt.Sprintf("%s Error: %s\n", method, err)
		return
	}

	fmt.Printf("%s=\n%s\n", method, at)

	return
}

//Depth 是反馈市场深度信息
func (a *API) Depth(money, coin string) (*ec.Depth, error) {
	rawData, err := a.depthRawData(money, coin)
	if err != nil {
		return nil, err
	}

	if a.IsLog {
		log.Printf(`rawData %s.Depth("%s","%s")=%s`, a.Name(), money, coin, string(rawData))
		log.Println()
	}

	resp := depth{}
	err = json.Unmarshal(rawData, &resp)
	if err != nil {
		return nil, err
	}

	if a.IsLog {
		log.Printf(`After json.Unmarshal: %s.Depth("%s","%s")=%v`, a.Name(), money, coin, resp)
		log.Println()
	}

	res := depthNormalize(resp)

	return res, nil
}

func (a *API) depthRawData(money, coin string) ([]byte, error) {
	path := urlMaker(depthURL, money, coin)
	return a.Get(path)
}

func depthNormalize(d depth) *ec.Depth {
	a := quotations(d.Asks)
	if !a.IsAskSorted() {
		a.SortAsks()
	}

	b := quotations(d.Bids)
	if !b.IsBidSorted() {
		b.SortBids()
	}

	return &ec.Depth{
		Asks: a,
		Bids: b,
	}
}

// TransRecords 返回市场的交易记录
// 当tid<=0时，返回最新的30条记录
func (a *API) TransRecords(money, coin string, tid int64) (ec.Trades, error) {
	rawData, err := a.transRecordsRawData(money, coin, tid)
	if err != nil {
		return nil, err
	}

	if a.IsLog {
		log.Printf(`rawData %s.TransRecords("%s", "%s", %d)=%s`, a.Name(), money, coin, tid, string(rawData))
		log.Println()
	}

	resp := ec.Trades{}
	err = json.Unmarshal(rawData, &resp)
	if err != nil {
		return nil, err
	}

	if a.IsLog {
		log.Printf(`After json.Unmarshal: %s.TransRecords("%s", "%s", %d)=%v`, a.Name(), money, coin, tid, resp[:5])
		log.Println()
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
