package exchanges

import (
	"fmt"
	"sort"
)

//Ticker 是ticker的数据结构。
type Ticker struct {
	Last float64
	Buy  float64
	Sell float64
	High float64
	Low  float64
	Vol  float64
}

func (t *Ticker) String() string {
	str := fmt.Sprintf("Last:%f\n", t.Last)
	str += fmt.Sprintf("Buy :%f\n", t.Buy)
	str += fmt.Sprintf("Sell:%f\n", t.Sell)
	str += fmt.Sprintf("High:%f\n", t.High)
	str += fmt.Sprintf("Low :%f\n", t.Low)
	str += fmt.Sprintf("Vol :%f\n", t.Vol)
	return str
}

//Depth 记录深度信息
type Depth struct {
	Asks Quotations
	Bids Quotations
}

func (d *Depth) String() string {
	str := "Asks"
	str += d.Asks.String()
	str += "Bids"
	str += d.Bids.String()
	return str
}

//Quotation 是报价单的意思。
type Quotation struct {
	Price, Amount float64
}

//Quotations 是卖出价的报价单表
type Quotations []Quotation

func (qs Quotations) String() string {
	str := fmt.Sprint("\tPrice\t\tAmount\n")
	for _, q := range qs {
		str += fmt.Sprintf("\t%f\t%f\n", q.Price, q.Amount)
	}
	return str
}

//IsAskSorted 判断一个Quotations是否是按照升序排列的。
func (qs Quotations) IsAskSorted() bool {
	for i := 1; i < qs.Len(); i++ {
		if qs[i-1].Price > qs[i].Price {
			return false
		}
	}
	return true
}

//IsBidSorted 判断一个Quotations是否是按照降序排列的。
func (qs Quotations) IsBidSorted() bool {
	for i := 1; i < qs.Len(); i++ {
		if qs[i-1].Price < qs[i].Price {
			return false
		}
	}
	return true
}

//Len returns length of ts
func (qs Quotations) Len() int {
	return len(qs)
}

//Less 决定了是升序还是降序
func (qs Quotations) Less(i, j int) bool {
	return qs[i].Price < qs[j].Price
}

//Swap 是交换方式
func (qs Quotations) Swap(i, j int) {
	qs[i], qs[j] = qs[j], qs[i]
}

//SortAsks 以Asks的升序进行原地排列
func (qs Quotations) SortAsks() {
	sort.Sort(qs)
}

//SortBids 以Bids的降序进行原地排序
func (qs Quotations) SortBids() {
	sort.Sort(sort.Reverse(qs))
}

//Order 是交易所的订单信息
type Order struct {
	ID     int64
	Date   int64
	Money  string
	Price  float64
	Coin   string
	Amount float64
	Type   string
}

func (o *Order) String() string {
	str := fmt.Sprintf("ID    :%d\n", o.ID)
	str += fmt.Sprintf("Date  :%d\n", o.Date)
	str += fmt.Sprintf("Money :%s\n", o.Money)
	str += fmt.Sprintf("Price :%f\n", o.Price)
	str += fmt.Sprintf("Coin  :%s\n", o.Coin)
	str += fmt.Sprintf("Amount:%f\n", o.Amount)
	str += fmt.Sprintf("Type  :%s\n", o.Type)
	return str
}

//OrderType 指定了交易的类型
type OrderType string

const (
	//BUY 是使用money换coin的过程
	BUY OrderType = "buy"
	//SELL 是使用coin换money的过程
	SELL OrderType = "sell"
)
