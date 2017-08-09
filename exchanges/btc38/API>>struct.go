package btc38

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/aQuaYi/GoKit"
	"github.com/aQuaYi/ToDaMoon/exchanges"
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

func (t ticker) normalize() *exchanges.Ticker {
	return &exchanges.Ticker{
		High: t.High,
		Low:  t.Low,
		Last: t.Last,
		Vol:  t.Vol,
		Buy:  t.Buy,
		Sell: t.Sell,
	}
}

type depth struct {
	Asks [][2]float64
	Bids [][2]float64
}

func quotations(depthData [][2]float64) exchanges.Quotations {
	res := make(exchanges.Quotations, len(depthData))
	for i, d := range depthData {
		res[i] = exchanges.Quotation{
			Price:  d[0],
			Amount: d[1],
		}
	}

	return res
}

//MyBalance 是btc38的账户信息
type myBalance map[string]string

func (m myBalance) normalize(coins []string) (*exchanges.Account, error) {
	mba, err := convertMyBalanceAmount(m)
	if err != nil {
		msg := fmt.Sprintf("转换时失败：%s", err)
		return nil, errors.New(msg)
	}

	a := exchanges.NewAccount()
	a.TotalCNY = mba["cny_balance"]

	coins = append(coins, "cny")

	for _, coin := range coins {
		tKey, fKey := getMyBalanceKeys(coin)
		total := mba[tKey]
		freezed := mba[fKey]
		available := total - freezed
		a.Coins[coin] = exchanges.CoinStatus{
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
func (t Trade) normalize() *exchanges.Trade {
	return &exchanges.Trade{
		Amount: t.Amount,
		Date:   t.Date,
		Price:  t.Price,
		Tid:    t.TradeID,
		Type:   t.Type,
	}
}

type order struct {
	ID        int64   `json:"id,string"`
	Coin      string  `json:"coinname"`
	OrderType int     `json:"type,string"`
	Amount    float64 `json:"amount,string"`
	Price     float64 `json:"price,string"`
	Time      string  `json:"time"`
}

var oTypeMap = map[int]string{
	1: "buy",
	2: "sell",
}

func (o order) normalize(money string) (*exchanges.Order, error) {
	t, err := time.Parse("2006-01-02 15:04:05", o.Time)
	if err != nil {
		msg := fmt.Sprintf(`无法把"%s"转换成time.time:%s`, o.Time, err) + "\n"
		return nil, errors.New(msg)
	}

	date := t.Unix()
	oType := oTypeMap[o.OrderType]

	return &exchanges.Order{
		ID:     o.ID,
		Date:   date,
		Money:  money,
		Price:  o.Price,
		Coin:   o.Coin,
		Amount: o.Amount,
		Type:   oType,
	}, nil
}

//MyTrade 是btc38的交易记录的格式
type MyTrade struct {
	ID       int64   `json:"id,string"`
	BuyerID  int     `json:"buyer_id,string"`
	SellerID int     `json:"seller_id,string"`
	Volume   float64 `json:"volume,string"`
	Price    float64 `json:"price,string"`
	Coin     string  `json:"coinname"`
	Time     string  `json:"time"`
}

func (mt MyTrade) String() string {
	str := fmt.Sprintf("ID      :%d\n", mt.ID)
	str += fmt.Sprintf("BuyerID :%d\n", mt.BuyerID)
	str += fmt.Sprintf("SellerID:%d\n", mt.SellerID)
	str += fmt.Sprintf("Volume  :%f\n", mt.Volume)
	str += fmt.Sprintf("Price   :%f\n", mt.Price)
	str += fmt.Sprintf("Coin    :%s\n", mt.Coin)
	str += fmt.Sprintf("Time    :%s\n", mt.Time)
	return str
}

func (mt MyTrade) normalize(myUserID int) (*exchanges.Trade, error) {
	d, err := GoKit.ParseLocalTime(mt.Time)
	if err != nil {
		msg := fmt.Sprintf("无法把%s转换timestamp: %s", mt.Time, err)
		return nil, errors.New(msg)
	}

	t := ""
	switch myUserID {
	case mt.BuyerID:
		t = "buy"
	case mt.SellerID:
		t = "sell"
	default:
		msg := fmt.Sprintf("交易记录不包含本人ID=%d", myUserID)
		return nil, errors.New(msg)
	}

	return &exchanges.Trade{
		Tid:    mt.ID,
		Date:   d.Unix(),
		Price:  mt.Price,
		Amount: mt.Volume,
		Type:   t,
	}, nil
}
