package btc38

import ec "ToDaMoon/exchanges"

// TickerResponse is TickerResponse
type TickerResponse struct {
	Ticker Ticker `json:"ticker"`
}

// Ticker is Ticker
type Ticker struct {
	High float64 `json:"high,float64"`
	Low  float64 `json:"low,float64"`
	Last float64 `json:"last,float64"`
	Vol  float64 `json:"vol,float64"`
	Buy  float64 `json:"buy,float64"`
	Sell float64 `json:"sell,float64"`
}

func (t Ticker) normalize() *ec.Ticker {
	return &ec.Ticker{
		High: t.High,
		Low:  t.Low,
		Last: t.Last,
		Vol:  t.Vol,
		Buy:  t.Buy,
		Sell: t.Sell,
	}
}

//MyBalance 是btc38的账户信息
type myBalance map[string]string

func (m myBalance) normalize() *ec.Account {
	//TODO: 完成myBalance的转换程序
	return nil
}

// Trade :okcoin's Trade struct
type Trade struct {
	Amount  float64 `json:"amount,float64"`
	Date    int64   `json:"date,int64"`
	Price   float64 `json:"price,float64"`
	TradeID int64   `json:"tid,int64"`
	Type    string  `json:"type"`
}

// Normalize change Trade date to standard formate
func (t Trade) normalize() *ec.Trade {
	return &ec.Trade{
		Amount: t.Amount,
		Date:   t.Date,
		Price:  t.Price,
		Tid:    t.TradeID,
		Type:   t.Type,
	}
}
