package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"sync"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	os.Remove("./foo.db")

	ncsdbfilename := "./noCreateStatement.db"
	os.Remove(ncsdbfilename)

	dbt, err := sql.Open("sqlite3", ncsdbfilename)
	if err != nil {
		msg := fmt.Sprintf("无法创建%s数据库，出错原因:%s", ncsdbfilename, err)
		log.Println(msg)
	}
	defer dbt.Close()

	sqlStmt := `
	create table foo (id integer not null primary key, name text);
	`
	//NOTICE: 	delete from foo;
	_, err = dbt.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	db2, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		log.Fatal("db2", err)
	}
	//	defer db2.Close()

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go insert(db, 0, 0, wg)
	go insert(db2, time.Second*2, 100, wg)

	go query(db2)

	wg.Wait()
	go queryMaxID(db2)

	time.Sleep(time.Millisecond * 500)
	wg.Add(1)
	insert(db2, time.Second*2, 10000, wg)
}

func insert(db *sql.DB, waitTime time.Duration, index int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("insert之前，先休息%d秒\n", waitTime)
	time.Sleep(waitTime)

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	defer tx.Commit()
	fmt.Println(index, "After begin")

	stmt, err := tx.Prepare("insert into foo(id, name) values(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for i := index; i < index+10; i++ {
		p := newPeople(i, fmt.Sprintf("こんにちわ世界%05d", i))
		_, err = stmt.Exec(p.attributes()...)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(index, i)
		if i == 10000+5 {
			return
		}
		time.Sleep(time.Millisecond * 500)
	}
	fmt.Println(index, "commit")
}

func query(db *sql.DB) {
	fmt.Println("entered query")
	time.Sleep(time.Millisecond * 250)

	i := 0
	for {
		time.Sleep(time.Millisecond * 500)

		go func() {
			fmt.Println("\t\t\t", i)
			rows, err := db.Query("select id, name from foo")
			if err != nil {
				log.Fatal(err)
			}
			defer rows.Close()
			for rows.Next() {
				p := &people{}
				err = rows.Scan(p.attributes()...)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println("\t\t\t", p.id, p.name)
			}
			err = rows.Err()
			if err != nil {
				log.Fatal(err)
			}
		}()
		i++
	}
}
func queryMaxID(db *sql.DB) {
	fmt.Println("entered query")
	time.Sleep(time.Millisecond * 250)

	i := 0
	for {
		time.Sleep(time.Millisecond * 500)

		go func() {
			fmt.Println("\t\t\t", i)
			rows, err := db.Query("select max(id) from foo")
			if err != nil {
				log.Fatal(err)
			}
			defer rows.Close()
			for rows.Next() {
				result := new(int)
				err = rows.Scan(result)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println("\t\tmax id\t", *result)
			}
			err = rows.Err()
			if err != nil {
				log.Fatal(err)
			}
		}()
		i++
	}
}

type people struct {
	id   int
	name string
}

func (p *people) attributes() []interface{} {
	return []interface{}{&p.id, &p.name}
}

func newPeople(id int, name string) *people {
	return &people{
		id:   id,
		name: name,
	}
}
