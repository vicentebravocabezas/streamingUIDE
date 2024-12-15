package database

import (
	"database/sql"
	"log"
	"sync"

	_ "github.com/tursodatabase/go-libsql"
)

var (
	once   sync.Once
	db     *sql.DB
	dbFile = "file:./database/local/db.sqlite3"
)

// esta función solo debe usarse para testing, al llamarse, se cerrara cualquier conexión a una DB
// y los siguientes llamados a DB() utilizaran la nueva ruta de la base de datos
func SetDBFile(dbPath string) {
	if db != nil {
		db.Close()
	}
	dbFile = dbPath
	once = sync.Once{}
}

func DB() *sql.DB {
	once.Do(func() {
		var err error
		db, err = sql.Open("libsql", dbFile)
		if err != nil {
			log.Fatal(err)
		}

		if err := db.Ping(); err != nil {
			log.Fatal(err)
		}
	})

	return db
}
