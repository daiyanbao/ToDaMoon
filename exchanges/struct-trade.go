package exchanges

//Trade 记录一个成交记录的细节
type Trade struct {
	Tid    int64
	Date   int64
	Price  float64
	Amount float64
	Type   string
}

//Attributes 实现了database.Attributer接口
func (t *Trade) Attributes() []interface{} {
	return []interface{}{&t.Tid, &t.Date, &t.Price, &t.Amount, &t.Type}
}

//Trades 是*Trade的切片
//因为会有很多关于[]Trade的操作，所以，设置了这个方法。
type Trades []*Trade
