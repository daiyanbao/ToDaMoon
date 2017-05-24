// Package database to manage database
package database

import (
	ec "ToDaMoon/exchanges"
	"ToDaMoon/util"
	"aQuaGo/common"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	// go-sqlite3 is sqlite3 interface
	_ "github.com/mattn/go-sqlite3"
)

// DBM manage some one database
type DBM chan DBCmd

// DBCmd contains all command needed for DBM(datamanager)
type DBCmd struct {
	Action     DBAction
	Stmt       string
	StartID    int64 //included
	EndID      int64 //not included
	InData     ec.Trades
	TradesChan chan<- ec.Trades
	ResultChan chan<- int64
	ErrChan    chan<- error
}

// DBAction is enumeration type for DB Action
type DBAction int

// enumeration data
const (
	Insert DBAction = iota
	QueryTrades
	QueryOne
)

var dBMap = make(map[string]DBM)
var mutex sync.Mutex

// New makes new DBM
func Connect(dir, market, coin, money string) (DBM, error) {
	var dbm DBM
	db := dbFilePath(dir, market, coin, money)
	if !common.Exist(db) {
		err := creat(db)
		if err != nil {
			return nil, common.WrapErr("sqlite.go:New():", err)
		}
	}

	if m, exist := dBMap[db]; exist {
		dbm = m
	} else {
		mutex.Lock()
		defer mutex.Unlock()
		dbm = make(DBM, 10)
		dBMap[db] = dbm
		go dbm.start(db)
	}

	return dbm, nil
}

func (dbm DBM) start(dbPath string) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalln("DBM start():", err)
	}
	defer db.Close()

	for cmd := range dbm {
		switch cmd.Action {
		case Insert:
			err = insert(db, cmd.InData)
			cmd.ErrChan <- err
		case QueryTrades:
			ts, err := queryTrades(db, cmd.StartID, cmd.EndID)
			cmd.TradesChan <- ts
			cmd.ErrChan <- err
		case QueryOne:
			r, err := queryOne(db, cmd.Stmt)
			cmd.ResultChan <- r
			cmd.ErrChan <- err
		default:
			log.Println("DBM.strat(): switch: Wrong Action.")
		}
	}
}

//Insert data
func (dbm DBM) Insert(ts ec.Trades) error {
	if unique, ids := ts.IsUnique(); !unique {
		log.Println("重复id是", ids)
		fmt.Println(ts)
		return errors.New("插入的数据，不是唯一的")
	}
	ec := make(chan error)
	dbm <- DBCmd{Action: Insert,
		InData:  ts,
		ErrChan: ec,
	}
	return <-ec
}

// MaxTid returns the max tid of database
func (dbm DBM) MaxTid() (int64, error) {
	if length, err := dbm.Len(); length == 0 && err == nil {
		return 1, nil
	}
	reply := make(chan int64)
	ec := make(chan error)
	dbm <- DBCmd{Action: QueryOne,
		Stmt:       "select max(tid) from raw",
		ResultChan: reply,
		ErrChan:    ec,
	}
	return <-reply, <-ec
}

// MaxDate returns the max date of this database
func (dbm DBM) MaxDate() (int64, error) {
	tid, _ := dbm.MaxTid()
	if tid == 1 {
		return 0, errors.New("this database is empty")
	}
	ec := make(chan error)
	reply := make(chan int64)

	dbm <- DBCmd{Action: QueryOne,
		Stmt:       fmt.Sprintf("select date from raw where tid = %d", tid),
		ResultChan: reply,
		ErrChan:    ec,
	}
	return <-reply, <-ec
}

// MinTid returns the min tid of this database
func (dbm DBM) MinTid() (int64, error) {
	if length, err := dbm.Len(); length == 0 && err == nil {
		return 1, nil
	}
	reply := make(chan int64)
	ec := make(chan error)
	dbm <- DBCmd{Action: QueryOne,
		Stmt:       "select min(tid) from raw",
		ResultChan: reply,
		ErrChan:    ec,
	}
	return <-reply, <-ec
}

// MinDate returns the min date of this database
func (dbm DBM) MinDate() (int64, error) {
	tid, _ := dbm.MinTid()
	ec := make(chan error)
	reply := make(chan int64)

	dbm <- DBCmd{Action: QueryOne,
		Stmt:       fmt.Sprintf("select date from raw where tid = %d", tid),
		ResultChan: reply,
		ErrChan:    ec,
	}
	return <-reply, <-ec
}

// MaxTidBefore returns the max tid before some date
func (dbm DBM) MaxTidBefore(date int64) (tid int64, err error) {
	if length, _ := dbm.Len(); length < 3 {
		err = errors.New("This database is almost empty")
		return 0, err
	}
	sTid, _ := dbm.MinTid()
	eTid, _ := dbm.MaxTid()

	for {
		fmt.Println("sTid, eTid\n", sTid, eTid)
		if eTid-sTid < 10 {
			fmt.Println("return sTid, nil")
			return sTid, nil
		}
		tid = (eTid + sTid) / 2
		newDate, tid := dbm.MaxDateBefore(tid)
		if newDate < date && date-newDate < 3600 {
			fmt.Println("tid-sTid=", tid-sTid)
			fmt.Println("return tid, nil")
			return tid, nil
		}
		if newDate < date {
			sTid = tid
		} else {
			eTid = tid
		}
		fmt.Println("===============", time.Now())
		time.Sleep(time.Second)
	}
}

// MaxDateBefore is Max Date Before a tid
func (dbm DBM) MaxDateBefore(tid int64) (int64, int64) {
	date, err := dbm.DateOf(tid)
	for err != nil {
		tid--
		date, err = dbm.DateOf(tid)
		fmt.Println("MaxDateBefore: tid, date \n ", tid, date)
	}
	return date, tid
}

// DateOf return the date of a tid
func (dbm DBM) DateOf(tid int64) (date int64, err error) {
	reply := make(chan int64)
	ec := make(chan error)

	format := "select date from raw where tid = %d"
	stmt := fmt.Sprintf(format, tid)

	dbm <- DBCmd{Action: QueryOne,
		Stmt:       stmt,
		ResultChan: reply,
		ErrChan:    ec,
	}
	return <-reply, <-ec
}

// Len returns the count of this database
func (dbm DBM) Len() (int64, error) {
	reply := make(chan int64)
	ec := make(chan error)
	dbm <- DBCmd{Action: QueryOne,
		Stmt:       "select count(*) from raw",
		ResultChan: reply,
		ErrChan:    ec,
	}
	return <-reply, <-ec
}

// Trades 获取数据库中的交易记录。返回结果会包含startID，**不包含**endID的记录。
func (dbm DBM) Trades(startID, endID int64) (ec.Trades, error) {
	//endID为0说明，需要查询从startID到末尾的全部数据。
	if endID == 0 {
		endID = 9223372036854775807 //int64的最大值
	}
	reply := make(chan ec.Trades)
	ec := make(chan error)
	dbm <- DBCmd{Action: QueryTrades,
		StartID:    startID,
		EndID:      endID,
		TradesChan: reply,
		ErrChan:    ec,
	}
	return <-reply, <-ec
}

// TradesChan 返回的结果会包含startID，**不包含**endID的记录。
func (dbm DBM) TradesChan(startTid, endTid int64) <-chan ec.Trades {
	out := make(chan ec.Trades, 16)
	//var err error

	if endTid == 0 {
		endTid = 1<<63 - 1
	}
	var finalDate, finalTid int64

	go func() {
		sTid := startTid
		done := false
		for !done {
			eTid := sTid + 2000000 // 2M records every time
			if endTid <= eTid {
				eTid = endTid
				done = true
			}

			ts, err := dbm.Trades(sTid, eTid)
			if err != nil {
				log.Println("TradesChan:", err)
				continue
			}
			sTid = eTid

			if ts.Len() == 0 {
				fmt.Println("*****in the TradesChan*****")
				fmt.Println("\tfinalDate=", finalDate, util.DateOf(finalDate))
				fmt.Println("\tfinalTid=", finalTid)
				fmt.Println("***************************")
				break
			}
			finalDate = ts[ts.Len()-1].Date
			finalTid = ts[ts.Len()-1].Tid
			out <- ts
		}
		close(out)
	}()

	return out
}

//以下是基础方法

//往数据库内插入数据。
func insert(db *sql.DB, ts ec.Trades) error {
	tx, err := db.Begin()
	if err != nil {
		return errors.New("insert(): db.Begin: " + err.Error())
	}

	format := "insert into raw(tid, date, price, amount, type) values(?,?,?,?,?)"
	stmt, err := tx.Prepare(format)
	if err != nil {
		return errors.New("insert(): tx.Prepare: " + err.Error())
	}
	defer stmt.Close()
	if unique, _ := ts.IsUnique(); !unique {
		log.Fatalln("插入前的检查，还是有重复数据。")
	}
	for _, t := range ts {
		_, err := stmt.Exec(t.Tid, t.Date, t.Price, t.Amount, t.Type)
		if err != nil {
			if unique, ids := ts.IsUnique(); !unique {
				log.Println("插入数据库过程中，数据库已经报警数据不唯一", ids)
			}
			log.Println("看看")
			ts.PrintIDDiff()
			log.Fatalln("sqlite.go: insert(): stme.Exec(): Tid: ", t.Tid, err)
		}
	}
	tx.Commit()
	return nil
}

//查询批量数据
func queryTrades(db *sql.DB, StartID, EndID int64) (ec.Trades, error) {
	format := "select * from raw where %d <= tid and tid < %d"
	stmt := fmt.Sprintf(format, StartID, EndID)

	rows, err := db.Query(stmt)
	if err != nil {
		return nil, errors.New("queryTrades: " + err.Error())
	}
	defer rows.Close()
	ts := ec.Trades{}
	for rows.Next() {
		var tid, date int64
		var price, amount float64
		var tType string
		err = rows.Scan(&tid, &date, &price, &amount, &tType)
		if err != nil {
			return nil, errors.New("queryTrades: rows.Scan: " + err.Error())
		}
		ts = append(ts, ec.Trade{Tid: tid,
			Date:   date,
			Price:  price,
			Amount: amount,
			Type:   tType})
	}
	err = rows.Err()
	if err != nil {
		return nil, errors.New("queryTrades: rows.Err: " + err.Error())
	}
	return ts, nil
}

//查询单个数据。
func queryOne(db *sql.DB, s string) (int64, error) {
	stmt, err := db.Prepare(s)
	if err != nil {
		return 0, errors.New("queryOne(): db.Prepare: " + err.Error())
	}
	defer stmt.Close()

	var result int64
	err = stmt.QueryRow().Scan(&result)
	if err != nil {
		return 1, err
	}
	return result, nil
}

//创建SQLite数据库文件
func creat(dbPath string) error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return common.WrapErr("sqlite.go:creat:"+dbPath, err)
	}
	defer db.Close()

	//在db中创建一个表，并清空表中所有的行。
	sqlStmt := `create table raw (	tid integer primary key, 
									date integer NOT NULL,
									price real NOT NULL,
									amount real NOT NULL,
									type text NOT NULL);
				delete from raw;`
	_, err = db.Exec(sqlStmt) //数据库执行以上语句。
	if err != nil {
		return common.WrapErr("sqlite.go:执行创建语句:", err)
	}
	return nil
}

//生成数据库文件路径
func dbFilePath(dir, market, coin, money string) string {
	return fmt.Sprintf("%s/%s/%s-%s.db", dir, market, coin, money)
}
