package btc38

import (
	ec "ToDaMoon/exchanges"
	"fmt"
)

// TickerResponse is TickerResponse
type TickerResponse struct {
	Date   string
	Ticker Ticker
}

// Ticker is Ticker
type Ticker struct {
	Buy  float64 `json:"buy,string"`
	High float64 `json:"high,string"`
	Last float64 `json:"last,string"`
	Low  float64 `json:"low,string"`
	Sell float64 `json:"sell,string"`
	Vol  float64 `json:"vol,string"`
}

// UserInfo cantains user's details in okcoin.cn
type UserInfo struct {
	Info struct {
		Funds struct {
			Asset struct {
				Net   float64 `json:"net,string"`
				Total float64 `json:"total,string"`
			} `json:"asset"`
			Borrow struct {
				BTC float64 `json:"btc,string"`
				LTC float64 `json:"ltc,string"`
				CNY float64 `json:"cny,string"`
			} `json:"borrow"`
			Free struct {
				BTC float64 `json:"btc,string"`
				LTC float64 `json:"ltc,string"`
				CNY float64 `json:"cny,string"`
			} `json:"free"`
			Freezed struct {
				BTC float64 `json:"btc,string"`
				LTC float64 `json:"ltc,string"`
				CNY float64 `json:"cny,string"`
			} `json:"freezed"`
			UnionFund struct {
				BTC float64 `json:"btc,string"`
				LTC float64 `json:"ltc,string"`
			} `json:"union_fund"`
		} `json:"funds"`
	} `json:"info"`
	Result    bool  `json:"result"`
	ErrorCode int64 `json:"error_code"`
}

func (ui *UserInfo) String() string {
	str := fmt.Sprintf("Result: %t\n", ui.Result)
	str += fmt.Sprint("Info:\n\tFunds:\n")
	str += fmt.Sprint("\t\tAsset:\n")
	str += fmt.Sprintf("\t\t\tNet:%f\n", ui.Info.Funds.Asset.Net)
	str += fmt.Sprintf("\t\t\tTotal:%f\n", ui.Info.Funds.Asset.Total)
	str += fmt.Sprint("\t\tBorrow:\n")
	str += fmt.Sprintf("\t\t\tBTC:%f\n", ui.Info.Funds.Borrow.BTC)
	str += fmt.Sprintf("\t\t\tLTC:%f\n", ui.Info.Funds.Borrow.LTC)
	str += fmt.Sprintf("\t\t\tCNY:%f\n", ui.Info.Funds.Borrow.CNY)
	str += fmt.Sprint("\t\tFree:\n")
	str += fmt.Sprintf("\t\t\tBTC:%f\n", ui.Info.Funds.Free.BTC)
	str += fmt.Sprintf("\t\t\tLTC:%f\n", ui.Info.Funds.Free.LTC)
	str += fmt.Sprintf("\t\t\tCNY:%f\n", ui.Info.Funds.Free.CNY)
	str += fmt.Sprint("\t\tFreezed:\n")
	str += fmt.Sprintf("\t\t\tBTC:%f\n", ui.Info.Funds.Freezed.BTC)
	str += fmt.Sprintf("\t\t\tLTC:%f\n", ui.Info.Funds.Freezed.LTC)
	str += fmt.Sprintf("\t\t\tCNY:%f\n", ui.Info.Funds.Freezed.CNY)
	str += fmt.Sprint("\t\tUnionFund:\n")
	str += fmt.Sprintf("\t\t\tBTC:%f\n", ui.Info.Funds.UnionFund.BTC)
	str += fmt.Sprintf("\t\t\tLTC:%f\n", ui.Info.Funds.UnionFund.LTC)
	return str
}

// Trade :okcoin's Trade struct
type Trade struct {
	Amount  float64 `json:"amount,string"`
	Date    int64   `json:"date"`
	Price   float64 `json:"price,string"`
	TradeID int64   `json:"tid"`
	Type    string  `json:"type"`
}

// Normalize change Trade date to standard formate
func (t Trade) Normalize() ec.Trade {
	return ec.Trade{
		Amount: t.Amount,
		Date:   t.Date,
		Price:  t.Price,
		Tid:    t.TradeID,
		Type:   t.Type,
	}
}
