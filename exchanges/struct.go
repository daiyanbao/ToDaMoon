package exchanges

import "fmt"

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
	Asks [][2]float64
	Bids [][2]float64
}

func (d *Depth) String() string {
	str := fmt.Sprintln("Asks:")
	for _, v := range d.Asks {
		str += "\t" + fmt.Sprintln(v)
	}
	str += fmt.Sprintln("Bids:")
	for _, v := range d.Bids {
		str += "\t" + fmt.Sprintln(v)
	}
	return str
}

//Order 是交易所的订单信息
//TODO: 写完Order
type Order struct {
}
