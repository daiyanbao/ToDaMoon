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

	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
	create table foo (id integer not null primary key, name text);
	delete from foo;
	`
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
		_, err = stmt.Exec(i, fmt.Sprintf("こんにちわ世界%05d", i))
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
				var id int
				var name string
				err = rows.Scan(&id, &name)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println("\t\t\t", id, name)
			}
			err = rows.Err()
			if err != nil {
				log.Fatal(err)
			}
		}()
		i++
	}
}
