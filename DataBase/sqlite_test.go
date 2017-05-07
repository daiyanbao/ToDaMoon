// sqlite_test.go
package database

import (
	"os"
	"testing"
)

const (
	dir    string = "."
	market string = "test"
	coin   string = "test"
	money  string = "test"
)

func Test_dbFilePath(t *testing.T) {
	if dbFilePath(dir, market, coin, money) != "./test/test-test.db" {
		t.Error("无法创建正确的数据库文件路径。")
	}
}

func Test_creat(t *testing.T) {
	db := dbFilePath(dir, market, coin, money)
	os.Remove(db)

	err := creat(db)
	if err != nil {
		t.Fatalf("创建%s出错。", db)
	}
}

func Test_Singleton(t *testing.T) {
	db1, _ := New(dir, market, coin, money)
	db2, _ := New(dir, market, coin, money)
	if db1 != db2 || db1 == nil {
		t.Error("sqlite.go不是单例模式。")
	}
}

func Test_MaxTid(t *testing.T) {
	testDB, _ := New(dir, market, coin, money)
	maxTid, _ := testDB.MaxTid()
	if maxTid != 1 {
		t.Error("无法正确地从空白数据库中读取MaxTid。")
	}

}
