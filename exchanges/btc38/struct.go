package btc38

import (
	ec "ToDaMoon/exchanges"
	"errors"
	"fmt"
	"strconv"
)

type tickerResp struct {
	Ticker ticker `json:"ticker"`
}

type ticker struct {
	High float64 `json:"high,float64"`
	Low  float64 `json:"low,float64"`
	Last float64 `json:"last,float64"`
	Vol  float64 `json:"vol,float64"`
	Buy  float64 `json:"buy,float64"`
	Sell float64 `json:"sell,float64"`
}

func (t ticker) normalize() *ec.Ticker {
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

func (m myBalance) normalize(coins []string) (*ec.Account, error) {
	mba, err := convertMyBalanceAmount(m)
	if err != nil {
		msg := fmt.Sprintf("转换时失败：%s", err)
		return nil, errors.New(msg)
	}

	a := ec.NewAccount()
	a.TotalCNY = mba["cny_balance"]

	coins = append(coins, "cny")

	for _, coin := range coins {
		tKey, fKey := getMyBalanceKeys(coin)
		total := mba[tKey]
		freezed := mba[fKey]
		available := total - freezed
		a.Coins[coin] = ec.CoinStatus{
			Total:     total,
			Freezed:   freezed,
			Available: available,
		}
	}

	return a, nil
}

func convertMyBalanceAmount(m myBalance) (map[string]float64, error) {
	result := make(map[string]float64)

	for k, v := range m {
		if v == "0.000000" { //因为0项很多，所以用这个来加速。
			result[k] = 0
			continue
		}

		amount, err := strconv.ParseFloat(v, 64)
		if err != nil {
			msg := fmt.Sprintf("无法把%s的%s转换成float64: %s", k, v, err)
			return nil, errors.New(msg)
		}
		result[k] = amount
	}

	return result, nil
}

func getMyBalanceKeys(coin string) (totalKey, freezedKey string) {
	totalKey = fmt.Sprintf("%s_balance", coin)
	freezedKey = fmt.Sprintf("%s_balance_lock", coin)
	return totalKey, freezedKey
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
