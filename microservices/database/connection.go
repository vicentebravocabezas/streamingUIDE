package database

import (
	"database/sql"
	"log"
	"sync"

	_ "github.com/tursodatabase/go-libsql"
)

var (
	dbFile = "file:./local/db.sqlite3"
	once   sync.Once
	db     *sql.DB
)

// esta función solo debe usarse para testing. Al llamarse, se cerrará cualquier conexión a una DB.
// Los siguientes llamados a DB() utilizarán la nueva ruta de la base de datos
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
