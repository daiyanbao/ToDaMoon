package exchanges

import (
	"ToDaMoon/database"
	"ToDaMoon/util"
	"errors"
	"fmt"
)

//Trade 记录一个成交记录的细节
type Trade struct {
	Tid    int64
	Date   int64
	Price  float64
	Amount float64
	Type   string
}

func (t Trade) String() string {
	str := "*****************\n"
	str += fmt.Sprintf("Tid   :%d\n", t.Tid)
	str += fmt.Sprintf("Date  :%d (%s)\n", t.Date, util.DateOf(t.Date))
	str += fmt.Sprintf("Price :%f\n", t.Price)
	str += fmt.Sprintf("Amount:%f\n", t.Amount)
	str += fmt.Sprintf("Type  :%s\n", t.Type)
	return str
}

//Attributes 实现了database.Attributer接口
func (t *Trade) Attributes() []interface{} {
	return []interface{}{&t.Tid, &t.Date, &t.Price, &t.Amount, &t.Type}
}

//newTrade 返回了一个*Trade变量。
func newTrade() database.Attributer {
	return &Trade{}
}



//Trades 是*Trade的切片
//因为会有很多关于[]*Trade的操作，所以，设置了这个方法。
type Trades []*Trade
