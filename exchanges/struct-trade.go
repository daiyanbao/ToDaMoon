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

//Attributes 实现了database.Attributer接口
func (t *Trade) Attributes() []interface{} {
	return []interface{}{&t.Tid, &t.Date, &t.Price, &t.Amount, &t.Type}
}

//newTrade 返回了一个*Trade变量。
func newTrade() database.Attributer {
	return &Trade{}
}

//TradesDB 用来存放exchange的历史交易数据的数据库
type TradesDB struct {
	db database.DBer
}

//Name 返回数据库的名称
func (t *TradesDB) Name() string {
	return t.db.Name()
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

func (t *TradesDB) count() (int64, error) {
	stmt := "select count(tid) from raw"
	var result int64

	if err := t.db.GetValues(stmt, &result); err != nil {
		return 0, util.Err("*TradesDB.count(): ", err)
	}
	return result, nil
}

//MaxTid 返回数据库中最大的tid值
func (t *TradesDB) MaxTid() (int64, error) {
	c, err := t.count()
	if err != nil {
		return 0, util.Err("*TradesDB.MaxTid(): ", err)
	} else if c == 0 {
		return 1, nil
	}

	stmt := "select max(tid) from raw"
	var result int64
	if err := t.db.GetValues(stmt, &result); err != nil {
		return 0, util.Err("*TradesDB.MaxTid(): ", err)
	}
	return result, nil
}

//MaxDate 返回数据库中最大的date值
func (t *TradesDB) MaxDate() (int64, error) {
	maxTid, err := t.MaxTid()
	if err != nil {
		return 0, util.Err("*TradesDB.MaxDate(): ", err)
	} else if maxTid == 1 {
		msg := fmt.Sprintf("%s数据库中，还没有数据。无法读取MaxDate()", t.Name())
		return 0, errors.New(msg)
	}

	stmt := fmt.Sprintf("select date from raw where tid = %d", maxTid)
	var result int64
	if err := t.db.GetValues(stmt, &result); err != nil {
		return 0, util.Err("*TradesDB.MaxDate(): ", err)
	}

	return result, nil
}

//Trades 获取数据库中的交易记录。返回结果会包含beginTID，**不包含**endTID的记录。
//endID为-1说明，需要查询从startID到末尾的全部数据。
func (t *TradesDB) Trades(beginTID, endTID int64) (Trades, error) {
	if endTID == -1 {
		endTID = 9223372036854775807 //int64的最大值
	}

	stmt := fmt.Sprintf("select * from raw where %d <= tid and tid < %d", beginTID, endTID)
	rows, err := t.db.GetRows(stmt, newTrade)
	if err != nil {
		return nil, util.Err("*TradesDB.Trades(): ", err)
	}

	ts, err := convertToTrades(rows)
	if err != nil {
		return nil, util.Err("*Trades.Trades(): ", err)
	}

	return ts, nil
}

func convertToTrades(rows []interface{}) (Trades, error) {
	ts := make(Trades, len(rows))
	ok := false
	for i, r := range rows {
		ts[i], ok = r.(*Trade)
		if !ok {
			return nil, errors.New("无法转换%v到Trade类型。")
		}
	}

	return ts, nil
}

//Trades 是*Trade的切片
//因为会有很多关于[]*Trade的操作，所以，设置了这个方法。
type Trades []*Trade
