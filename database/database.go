package database

import (
	"database/sql"
	"log"
	"sync"

	_ "github.com/tursodatabase/go-libsql"
)

var (
	once sync.Once
	db   *sql.DB
)

func OpenDB() *sql.DB {
	once.Do(func() {
		var err error
		db, err = sql.Open("libsql", "file:database/local/db.sqlite3")
		if err != nil {
			log.Fatal(err)
		}
	})

	return db
}
