// Package database to manage database
package database

import (
	"ToDaMoon/util"
	"database/sql"
	"sync"
)

type creater interface {
	CreateStatement() string
}

type insertQuerier interface {
	InsertStatement() string
	QueryStatement() string
	Attributes() []interface{}
}

type dataer interface {
	creater
	insertQuerier
}

//DBer 是定制数据库的接口
type DBer interface {
}

//DB 定制的sql数据库
type DB struct {
	*sql.DB
	dataer
}

var mutex sync.Mutex

//Open 链接上了数据库
func Open(filename string, d dataer) (DBer, error) {
	{ //为了不重复创建数据库，加个锁
		mutex.Lock()

		if !util.Exist(filename) {
			if err := createDB(filename, d); err != nil {
				return nil, err
			}
		}

		mutex.Unlock()
	}

	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, util.Err("无法open"+filename, err)
	}

	result := &DB{
		DB:     db,
		dataer: d,
	}

	return result, nil
}

//Insert 向DB内插入数据
func (d *DB) Insert() {

}

func createDB(filename string, c creater) error {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return util.Err("无法创建数据库"+filename, err)
	}
	defer db.Close()

	//在db中创建一个表，并清空表中所有的行。
	stmt := c.CreateStatement()
	_, err = db.Exec(stmt) //数据库执行创建语句
	if err != nil {
		return util.Err("执行创建以下创建语句失败:"+stmt, err)
	}
	return nil
}
