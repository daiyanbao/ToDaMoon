package exchanges

import (
	"ToDaMoon/database"
)

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

//newTrade 返回了一个*Trade变量。
func newTrade() *Trade {
	return &Trade{}
}

//TradesDB 用来存放exchange的历史交易数据的数据库
type TradesDB struct {
	db database.DBer
}

//OpenTradesDB 连接上一个filename对应的数据库文件
func OpenTradesDB(filename string) (*TradesDB, error) {
	createStatement := `create table raw (
		tid integer primary key,
		date integer NOT NULL,
		price real NOT NULL,
		amount real NOT NULL,
		type text NOT NULL);`

	db, err := database.Connect(filename, createStatement)
	if err != nil {
		return nil, err
	}

	return &TradesDB{
		db: db,
	}, nil
}

//Trades 是*Trade的切片
//因为会有很多关于[]*Trade的操作，所以，设置了这个方法。
type Trades []*Trade
