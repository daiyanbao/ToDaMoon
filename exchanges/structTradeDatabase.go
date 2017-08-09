package exchanges

// import (
// 	"github.com/aQuaYi/ToDaMoon/database"
// 	"github.com/aQuaYi/ToDaMoon/util"
// 	"errors"
// 	"fmt"
// )

// //TradesDB 用来存放exchange的历史交易数据的数据库
// type TradesDB struct {
// 	db database.DBer
// }

// //Name 返回数据库的名称
// func (t *TradesDB) Name() string {
// 	return t.db.Name()
// }

// //OpenTradesDB 连接上一个filename对应的数据库文件
// func OpenTradesDB(filename string) (*TradesDB, error) {
// 	createStatement := `create table raw (
// 		tid integer primary key,
// 		date integer NOT NULL,
// 		price real NOT NULL,
// 		amount real NOT NULL,
// 		type text NOT NULL);`

// 	db, err := database.Connect(filename, createStatement)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &TradesDB{
// 		db: db,
// 	}, nil
// }

// //Len 返回了数据库的长度
// func (t *TradesDB) Len() (int64, error) {
// 	stmt := "select count(tid) from raw"
// 	var result int64

// 	if err := t.db.GetValues(stmt, &result); err != nil {
// 		return 0, GoKit.Err("*TradesDB.Len(): ", err)
// 	}
// 	return result, nil
// }

// //MaxTid 返回数据库中最大的tid值
// func (t *TradesDB) MaxTid() (int64, error) {
// 	c, err := t.Len()
// 	if err != nil {
// 		return 0, GoKit.Err("*TradesDB.MaxTid(): ", err)
// 	} else if c == 0 {
// 		//返回1的原因是，有的交易所会忽略掉tid=0的参数。
// 		return 1, nil
// 	}

// 	stmt := "select max(tid) from raw"
// 	var result int64
// 	if err := t.db.GetValues(stmt, &result); err != nil {
// 		return 0, GoKit.Err("*TradesDB.MaxTid(): ", err)
// 	}
// 	return result, nil
// }

// //MinTid 返回数据库中最大的tid值
// func (t *TradesDB) MinTid() (int64, error) {
// 	c, err := t.Len()
// 	if err != nil {
// 		return 0, GoKit.Err("*TradesDB.MinTid(): ", err)
// 	} else if c == 0 {
// 		//返回1的原因是，为了和MaxTid()一致
// 		return 1, nil
// 	}

// 	stmt := "select min(tid) from raw"
// 	var result int64
// 	if err := t.db.GetValues(stmt, &result); err != nil {
// 		return 0, GoKit.Err("*TradesDB.MinTid(): ", err)
// 	}
// 	return result, nil
// }

// //MaxTidNotGreaterThan 返回数据库中比参数小的tid中的最大值
// //tid通常不是连续的，甚至还有可能被交易所删除了记录，而缺上一大段区间。
// func (t *TradesDB) MaxTidNotGreaterThan(number int64) (int64, error) {
// 	c, err := t.Len()
// 	if err != nil {
// 		return 0, GoKit.Err("*TradesDB.MaxTidNotGreaterThan(): ", err)
// 	} else if c == 0 {
// 		return 1, nil
// 	}

// 	stmt := fmt.Sprintf("select max(tid) from raw where tid <= %d", number)
// 	var result int64
// 	if err := t.db.GetValues(stmt, &result); err != nil {
// 		return 0, GoKit.Err("*TradesDB.MaxTidNotGreaterThan(): ", err)
// 	}

// 	return result, nil
// }

// //MaxDate 返回数据库中最大的date值
// func (t *TradesDB) MaxDate() (int64, error) {
// 	maxTid, err := t.MaxTid()
// 	if err != nil {
// 		return 0, GoKit.Err("*TradesDB.MaxDate(): ", err)
// 	} else if maxTid == 1 {
// 		msg := fmt.Sprintf("%s数据库中，还没有数据。无法读取MaxDate()", t.Name())
// 		return 0, errors.New(msg)
// 	}

// 	stmt := fmt.Sprintf("select date from raw where tid = %d", maxTid)
// 	var result int64
// 	if err := t.db.GetValues(stmt, &result); err != nil {
// 		return 0, GoKit.Err("*TradesDB.MaxDate(): ", err)
// 	}

// 	return result, nil
// }

// //MinDate 返回数据库中最小的date值
// func (t *TradesDB) MinDate() (int64, error) {
// 	minTid, err := t.MinTid()
// 	if err != nil {
// 		return 0, GoKit.Err("*TradesDB.MaxDate(): ", err)
// 	} else if minTid == 1 {
// 		msg := fmt.Sprintf("%s数据库中，还没有数据。无法读取MinDate()", t.Name())
// 		return 0, errors.New(msg)
// 	}

// 	stmt := fmt.Sprintf("select date from raw where tid = %d", minTid)
// 	var result int64
// 	if err := t.db.GetValues(stmt, &result); err != nil {
// 		return 0, GoKit.Err("*TradesDB.MaxDate(): ", err)
// 	}

// 	return result, nil
// }

// //DateOf 返回number所对应的date
// func (t *TradesDB) DateOf(number int64) (int64, error) {
// 	tid, err := t.MaxTidNotGreaterThan(number)
// 	if err != nil {
// 		return 0, GoKit.Err("*TradesDB.DateOf(): ", err)
// 	}

// 	stmt := fmt.Sprintf("select date from raw where tid = %d", tid)

// 	var result int64
// 	if err := t.db.GetValues(stmt, &result); err != nil {
// 		return 0, GoKit.Err("*TradesDB.DateOf(): ", err)
// 	}

// 	return result, nil
// }

// //MaxTidBeforeDate 返回了在指定日期前的最大tid值
// //比如okcoin的btc数据库就非常大，读取全部内容，会非常耗时，只好读取一部分内容。
// func (t *TradesDB) MaxTidBeforeDate(date int64) (int64, error) {
// 	//首先处理数据库中的记录太短的情况
// 	firstDate, err := t.MinDate()
// 	if err != nil {
// 		return 0, GoKit.Err("*TradesDB.MaxTidBeforeDate(): ", err)
// 	}
// 	if date < firstDate {
// 		msg := fmt.Sprintf("*TradesDB.MaxTidBeforeDate(): %s中最老的数据，在%s之后才产生。", t.Name(), GoKit.DateOf(date))
// 		return 0, errors.New(msg)
// 	}

// 	//然后处理，数据库中的数据，还没有更新到最新记录的情况
// 	lastDate, err := t.MaxDate()
// 	if err != nil {
// 		return 0, GoKit.Err("*TradesDB.MaxTidBeforeDate(): ", err)
// 	}
// 	if lastDate < date {
// 		msg := fmt.Sprintf("*TradesDB.MaxTidBeforeDate(): %s中最新的数据，在%s之前。", t.Name(), GoKit.DateOf(date))
// 		return 0, errors.New(msg)
// 	}

// 	//最后处理Date落在数据库中的情况
// 	firstTID, err := t.MinTid()
// 	if err != nil {
// 		return 0, GoKit.Err("*TradesDB.MaxTidBeforeDate(): ", err)
// 	}

// 	lastTID, err := t.MaxTid()
// 	if err != nil {
// 		return 0, GoKit.Err("*TradesDB.MaxTidBeforeDate(): ", err)
// 	}

// 	for firstTID+4096 < lastTID {
// 		number := (firstTID + lastTID) / 2
// 		//因为前面已经排除了两种情况，所以，这里肯定不会出错。
// 		middleTid, _ := t.MaxTidNotGreaterThan(number)
// 		middleDate, _ := t.DateOf(middleTid)

// 		if firstTID == middleTid { //说明date落在了大缺口上
// 			return middleTid, nil
// 		}

// 		if middleDate < date {
// 			firstTID = middleTid
// 		} else {
// 			lastTID = middleTid
// 		}
// 	}

// 	return firstTID, nil
// }

// //Trades 获取数据库中的交易记录。返回结果会包含beginTID，**不包含**endTID的记录。
// //endID为-1说明，需要查询从startID到末尾的全部数据。
// func (t *TradesDB) Trades(beginTID, endTID int64) (Trades, error) {
// 	if endTID == -1 {
// 		endTID = 1<<63 - 1 //int64的最大值 9223372036854775807
// 	}

// 	stmt := fmt.Sprintf("select * from raw where %d <= tid and tid < %d", beginTID, endTID)
// 	rows, err := t.db.GetRows(stmt, newTrade)
// 	if err != nil {
// 		return nil, GoKit.Err("*TradesDB.Trades(): ", err)
// 	}

// 	ts, err := convertToTrades(rows)
// 	if err != nil {
// 		return nil, GoKit.Err("*Trades.Trades(): ", err)
// 	}

// 	return ts, nil
// }

// func convertToTrades(rows []interface{}) (Trades, error) {
// 	ts := make(Trades, len(rows))
// 	ok := false
// 	for i, r := range rows {
// 		ts[i], ok = r.(*Trade)
// 		if !ok {
// 			return nil, errors.New("无法转换%v到Trade类型。")
// 		}
// 	}

// 	return ts, nil
// }

// //Insert 往TradesDB数据库中插入数据。
// func (t *TradesDB) Insert(ts Trades) error {
// 	stmt := "insert into raw(tid, date, price, amount, type) values(?,?,?,?,?)"

// 	//go语言不会自动转换切片，要手动转换
// 	as := make([]database.Attributer, len(ts))
// 	for i, v := range ts {
// 		as[i] = database.Attributer(v)
// 	}

// 	if err := t.db.Insert(as, stmt); err != nil {
// 		return GoKit.Err("*TradesDB.Insert(): ", err)
// 	}

// 	return nil
// }
