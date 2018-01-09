package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

var (
	db *sql.DB
)

type dbHandler struct{}

// Hoge do this.
func Hoge() {
	fmt.Println("hoge")
}

func (h dbHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "DB Page\n\n")

	// connects to db server and create dummytable if not exists
	err := initDB()
	if err != nil {
		log.Fatal(err)
	}

	// check url has /db/add/*** format and if so, insert record to the dummytable
	if strings.HasPrefix(r.URL.Path, "/db/add/") {
		addnew := strings.TrimPrefix(r.URL.Path, "/db/add/")
		if addnew != "" {
			insertOne(addnew)
			fmt.Fprintf(w, "New item inserted. \n\n")
		}
	}

	fmt.Fprintf(w, "Access /db/add/{some text} to insert new item. \n\n")

	fmt.Fprintf(w, "Items from dummy table...\n\n")
	// list records from the dummytable
	items := listitems()
	for c, obkey := range items {
		fmt.Fprintln(w, c+1, ". ", obkey)
	}
}

func initDB() error {

	connInfo := "user=postgres dbname=postgres password=mypass host=db port=5432 sslmode=disable"

	var err error
	db, err = sql.Open("postgres", connInfo)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 5; i++ {
		time.Sleep(time.Duration(i) * time.Second)

		if err = db.Ping(); err == nil {
			break
		}
		log.Println(err)
	}
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(
		`create table if not exists dummytable (
			id SERIAL,
			dummytext TEXT NOT NULL,
			CONSTRAINT dummytable_pkey PRIMARY KEY (id)
		)`)

	return err
}

func listitems() []string {

	rows, err := db.Query("select dummytext from dummytable")
	if err != nil {
		log.Fatal(err)
	}

	var items []string
	for rows.Next() {
		var dummytext string
		err := rows.Scan(&dummytext)
		if err != nil {
			log.Fatal(err)
		}
		items = append(items, dummytext)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return items
}

func insertOne(text string) {
	r, err := db.Exec("insert into dummytable(dummytext) values('" + text + "')")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(r)
}
