package litedb

import (
	D "github.com/NeoJRotary/describe-go"
)

// UpdateConfig ...
func (db *DB) UpdateConfig(key string, val string) {
	_, err := db.Exec(`INSERT INTO config VALUES (?, ?);`, key, val)
	if D.IsErr(err) {
		db.Exec(`UPDATE config SET "value"=? WHERE "key"=?;`, val, key)
	}
}
