// Package database to manage database
package database

import (
	"ToDaMoon/util"
	"database/sql"
	"errors"
	"fmt"
	"log"
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
	return insert(d, ds)
}

func insert(db *DB, ds []dataer) error {
	//启动insert事务
	transaction, err := db.Begin()
	if err != nil {
		msg := fmt.Sprintf("%s无法启动一个insert事务:%s", db.name, err)
		return errors.New(msg)
	}
	defer transaction.Commit()

	//为insert事务进行准备工作
	statement := db.InsertStatement()
	stmt, err := transaction.Prepare(statement)
	if err != nil {
		msg := fmt.Sprintf("%s的insert事务的准备工作失败: %s", db.name, err)
		return errors.New(msg)
	}
	defer stmt.Close()

	//按行插入
	for _, d := range ds {
		_, err := stmt.Exec(d.Attributes())
		if err != nil {
			attrs := fmt.Sprint(d.Attributes())
			msg := fmt.Sprintf("%s在插入[%s]时出错: %s", db.name, attrs, err)
			//TODO: 直接关闭程序是否太严格了，考虑一下换成报错
			log.Fatalln(msg)
		}
	}

	return nil
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
