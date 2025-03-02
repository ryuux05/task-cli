package storage

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

const dbFile = "db/task.db"
func NewSqlite() (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		return nil, err
	}
	
	return db, nil

}