package litedb

import (
	"database/sql"

	D "github.com/NeoJRotary/describe-go"
	// go-sqlite3
	_ "github.com/mattn/go-sqlite3"
)

// DB ...
type DB struct {
	sqlDB    *sql.DB
	dataPath string
}

// Query ...
func (db *DB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return db.sqlDB.Query(query, args...)
}

// QueryRow ...
func (db *DB) QueryRow(query string, args ...interface{}) *sql.Row {
	return db.sqlDB.QueryRow(query, args...)
}

// Exec ...
func (db *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return db.sqlDB.Exec(query, args...)
}

// Init ...
func Init(dataPath string) (*DB, error) {
	sqldb, err := sql.Open("sqlite3", dataPath)
	if D.IsErr(err) {
		return nil, err
	}
	// defer db.Close()
	db := DB{
		sqlDB:    sqldb,
		dataPath: dataPath,
	}

	db.InitTables()

	return &db, nil
}
