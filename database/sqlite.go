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
}

//DB 定制的sql数据库
type DB struct {
	name string
	*sql.DB
	Datar
}

//Name 返回数据库的名称，其实也就是数据库的存放地址
func (db *DB) Name() string {
	return db.name
}

func New(filename string, d Datar) (DBer, error) {
	db, err := open(filename, d)
	if err != nil {
		return nil, err
	}

	return &DB{
		name:  filename,
		DB:    db,
		Datar: d,
	}, nil
}

var mutex sync.Mutex

//Open 链接上了数据库
func open(filename string, d Datar) (*sql.DB, error) {
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
		msg := fmt.Sprintf("无法打开sqlite3文件%s，出错原因:%s", filename, err)
		return nil, errors.New(msg)
	}

	return db, nil
}

//Insert 向DB内插入数据
func (db *DB) Insert(ds []Datar) error {
	//插入长度为0的数据，直接返回
	//检查长度是因为，后面会用到ds[0]，如果len(ds)==0，ds[0]会引起panic
	if len(ds) == 0 {
		return nil
	}

	//如果数据库本身数据的类型与待插入数据的类型不一致，无法插入
	item := db.NewItem()
	if !util.IsTypeEqual(item, ds[0]) {
		msg := fmt.Sprintf("数据库%s的数据的原始类型为%T，待插入数据的原始数据类型为%T，两者不符，无法插入。", db.name, item, ds[0])
		return errors.New(msg)
	}

	return insert(db, ds)
}

func insert(db *DB, ds []Datar) error {
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

func (db *DB) QueryBy(statement string) ([]interface{}, error) {
	rows, err := db.Query(statement)
	if err != nil {
		msg := fmt.Sprintf("对%s使用%s语句查询，出现错误:%s", db.name, statement, err)
		return nil, errors.New(msg)
	}
	defer rows.Close()

	result := []interface{}{}
	for rows.Next() {
		item := db.NewItem()
		err := rows.Scan(item.Attributes())
		if err != nil {
			msg := fmt.Sprintf("对%s查询%s出来的rows进行Scan时，出错:%s", db.name, statement, err)
			//TODO: 直接关闭程序是否太严格了，考虑一下换成报错
			log.Fatalln(msg)
		}
		result = append(result, item)
	}

	err = rows.Err()
	if err != nil {
		msg := fmt.Sprintf("对%s查询%s出来的rows，Scan完毕后，出错:%s", db.name, statement, err)
		return nil, errors.New(msg)
	}

	return result, nil
}

func createDB(filename string, c Creater) error {
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
