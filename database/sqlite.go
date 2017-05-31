// Package database to manage database
package database

import (
	"ToDaMoon/util"
	"database/sql"
	"errors"
	"fmt"
	"sync"
)

type creater interface {
	CreateStatement() string
}

type insertier interface {
	InsertStatement() string
	Attributes() []interface{}
}

type querier interface {
	QueryStatement() string
	NewItem() attributer
	Attributes() []interface{}
}

type attributer interface {
	Attributes() []interface{}
}
type dataer interface {
	creater
	InsertStatement() string
	QueryStatement() string
	NewItem() attributer
	Attributes() []interface{}
}

//DBer 是定制数据库的接口
type DBer interface {
}

//DB 定制的sql数据库
type DB struct {
	name string
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
func (d *DB) Insert(ds []dataer) error {
	//插入长度为0的数据，直接返回
	if len(ds) == 0 {
		return nil
	}

	//如果数据库本身数据的类型与待插入数据的类型不一致，无法插入
	if !util.IsTypeEqual(d.dataer, ds[0]) {
		msg := fmt.Sprintf("数据库%s的数据的原始类型为%T，待插入数据的原始数据类型为%T，两者不符，无法插入。", d.name, d.dataer, ds[0])
		return errors.New(msg)
	}
	return d.insert(ds)
}

func (d *DB) insert(ds []dataer) error {

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
