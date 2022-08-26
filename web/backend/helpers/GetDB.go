package helpers

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func GetDB() *sql.DB {
	db, _ := sql.Open("sqlite3", "./data/database.db")
	return db
}
