// Package database to manage database
package database

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"

	"github.com/aQuaYi/ToDaMoon/util"

	//导入sqlite3的驱动
	_ "github.com/mattn/go-sqlite3"
)

//Attributer 返回了struct的属性的指针位置组成的切片
//NOTICE: 切片中元素的顺序，要与对应查询语句中的元素顺序相同
type Attributer interface {
	Attributes() []interface{}
}

//DBer 是定制数据库的接口
type DBer interface {
	Name() string
	Insert(data []Attributer, statement string) error
	GetRows(statement string, newItem func() Attributer) ([]interface{}, error)
	GetValues(statement string, dest ...interface{}) error
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
		msg := fmt.Sprintf("%s无法启动一个insert事务:%s", db.Name(), err)
		return errors.New(msg)
	}
	defer transaction.Commit()

	//为insert事务进行准备工作
	stmt, err := transaction.Prepare(insertStatement)
	if err != nil {
		msg := fmt.Sprintf("%s的insert事务的准备以下insert语句时失败\n%s\n失败原因: %s", db.Name(), insertStatement, err)
		return errors.New(msg)
	}
	defer stmt.Close()

	//按行插入
	for _, d := range data {
		_, err := stmt.Exec(d.Attributes()...)
		if err != nil {
			attrs := fmt.Sprint(d)
			msg := fmt.Sprintf("%s在插入%s出错: %s", db.Name(), attrs, err)
			//NOTICE: 经过再三的思考，我决定在插入出错后，不要直接关闭程序。由程序的调用方来决定，如何处理错误。
			return errors.New(msg)
		}
	}

	return nil
}

//GetRows 描述了从数据库中查询的过程
func (db *DB) GetRows(queryStatement string, newItem func() Attributer) ([]interface{}, error) {
	rows, err := db.Query(queryStatement)
	if err != nil {
		msg := fmt.Sprintf("对%s使用以下语句查询\n%s\n出现错误:%s", db.Name(), queryStatement, err)
		return nil, errors.New(msg)
	}
	defer rows.Close()

	result := []interface{}{}
	for rows.Next() {
		item := newItem()
		err := rows.Scan(item.Attributes()...)
		if err != nil {
			msg := fmt.Sprintf("对%s查询%s出来的rows进行Scan时，出错:%s", db.Name(), queryStatement, err)
			//NOTICE: 经过再三的思考，我决定在查询出错后，不要直接关闭程序。由程序的调用方来决定，如何处理错误。
			return nil, errors.New(msg)
		}
		result = append(result, item)
	}

	err = rows.Err()
	if err != nil {
		msg := fmt.Sprintf("对%s查询%s出来的rows，Scan完毕后，出错:%s", db.Name(), queryStatement, err)
		return nil, errors.New(msg)
	}

	return result, nil
}

//GetValues 获取查询语句的多个值
//在dest存放变量的指针
//NOTICE: dest中变量的指针的顺序，需要与statement中查询的一样
func (db *DB) GetValues(statement string, dest ...interface{}) error {
	stmt, err := db.Prepare(statement)
	if err != nil {
		msg := fmt.Sprintf("对%s使用以下语句查询\n%s\n出现错误:%s", db.Name(), statement, err)
		return errors.New(msg)
	}
	defer stmt.Close()

	err = stmt.QueryRow().Scan(dest...) //NOTICE: Scan的参数必须打上...
	if err != nil {
		msg := fmt.Sprintf("database.GetValues：对%s查询%s出来的值，Scan完毕后，出错:%s", db.Name(), statement, err)
		return errors.New(msg)
	}

	return nil
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

//AttributerMaker 把别的类型的切片转换成[]Attributer
func AttributerMaker(is []interface{}) []Attributer {
	as := make([]Attributer, len(is))
	for i, v := range is {
		as[i] = v.(Attributer)
	}
	return as
}
