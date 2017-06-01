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

//Creater 返回用于创建数据库的table的语句
type Creater interface {
	CreateStatement() string
}

//Inserter 向数据库插入数据的接口
type Inserter interface {
	//返回，向数据库中插入数据的插入语句
	Statement() string

	//输出插入数据组成的切片，数据的顺序要求与插入语句中的顺序一致
	Attributer
}

//Querier 从数据库查询内容的接口
type Querier interface {
	Statement() string
	NewItem() Attributer
	Attributer
}

//Attributer 返回了struct的属性的指针位置组成的切片
//NOTICE: 切片中元素的顺序，要与对应查询语句中的元素顺序相同
type Attributer interface {
	Attributes() []interface{}
}

//Datar 反应了数据库中数据的特性
type Datar interface {
	Creater
	NewItem() Attributer
}

//DBer 是定制数据库的接口
type DBer interface {
	Name() string
	Insert([]Attributer, string) error
	QueryBy(string, func() Attributer) ([]interface{}, error)
}

//DB 定制的sql数据库
type DB struct {
	name string
	*sql.DB
}

//Name 返回数据库的名称，也是数据库的存放地址
func (db *DB) Name() string {
	return db.name
}

//Connect 返回一个数据库对象
func Connect(filename string, createStatement string) (DBer, error) {
	db, err := open(filename, createStatement)
	if err != nil {
		return nil, err
	}

	return &DB{
		name: filename,
		DB:   db,
	}, nil
}

var mutex sync.Mutex

//Open 链接上了数据库
func open(filename string, createStatement string) (*sql.DB, error) {
	{ //为了不重复创建数据库，加个锁
		mutex.Lock()
		defer mutex.Unlock()

		//如果不存在数据库文件不存在，就创建一个新的
		if !util.Exist(filename) {
			if err := createDB(filename, createStatement); err != nil {
				return nil, err
			}
		}
	}

	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		msg := fmt.Sprintf("无法打开sqlite3文件%s，出错原因:%s", filename, err)
		return nil, errors.New(msg)
	}

	return db, nil
}

//Insert 描述了向DB内插入数据的过程
func (db *DB) Insert(data []Attributer, insertStatement string) error {
	//启动insert事务
	transaction, err := db.Begin()
	if err != nil {
		msg := fmt.Sprintf("%s无法启动一个insert事务:%s", db.name, err)
		return errors.New(msg)
	}
	defer transaction.Commit()

	//为insert事务进行准备工作
	stmt, err := transaction.Prepare(insertStatement)
	if err != nil {
		msg := fmt.Sprintf("%s的insert事务的准备以下insert语句时失败\n%s\n失败原因: %s", db.name, insertStatement, err)
		return errors.New(msg)
	}
	defer stmt.Close()

	//按行插入
	for _, d := range data {
		_, err := stmt.Exec(d.Attributes()...)
		if err != nil {
			attrs := fmt.Sprint(d)
			msg := fmt.Sprintf("%s在插入\n%s\n出错: %s", db.name, attrs, err)
			//TODO: 直接关闭程序是否太严格了，考虑一下换成报错
			log.Fatalln(msg)
		}
	}

	return nil
}

//QueryBy 描述了从数据库中查询的过程
func (db *DB) QueryBy(queryStatement string, newItem func() Attributer) ([]interface{}, error) {
	rows, err := db.Query(queryStatement)
	if err != nil {
		msg := fmt.Sprintf("对%s使用以下语句查询\n%s\n出现错误:%s", db.name, queryStatement, err)
		return nil, errors.New(msg)
	}
	defer rows.Close()

	result := []interface{}{}
	for rows.Next() {
		item := newItem()
		err := rows.Scan(item.Attributes()...)
		if err != nil {
			msg := fmt.Sprintf("对%s查询%s出来的rows进行Scan时，出错:%s", db.name, queryStatement, err)
			//TODO: 直接关闭程序是否太严格了，考虑一下换成报错
			log.Fatalln(msg)
		}
		result = append(result, item)
	}

	err = rows.Err()
	if err != nil {
		msg := fmt.Sprintf("对%s查询%s出来的rows，Scan完毕后，出错:%s", db.name, queryStatement, err)
		return nil, errors.New(msg)
	}

	return result, nil
}

func createDB(filename string, createStatement string) error {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		msg := fmt.Sprintf("无法创建%s数据库，出错原因:%s", filename, err)
		return errors.New(msg)
	}
	defer db.Close()

	_, err = db.Exec(createStatement) //数据库执行创建语句
	if err != nil {
		msg := fmt.Sprintf("在%s数据库中执行以下create语句失败，\n%s\n失败原因:%s", filename, createStatement, err)
		return errors.New(msg)
	}

	return nil
}
