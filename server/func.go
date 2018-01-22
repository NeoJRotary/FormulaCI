package main

import (
	"encoding/json"
	"os"
)

// var execLog = log.New(os.Stdout, "\n[EXEC] ", log.LstdFlags)

func isErr(err error) bool {
	return err != nil
}

func fileExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func updateConfig(key string, val string) {
	_, err := sqlite.db.Exec(`INSERT INTO config VALUES (?, ?);`, key, val)
	if isErr(err) {
		sqlite.db.Exec(`UPDATE config SET "value"=? WHERE "key"=?;`, val, key)
	}
	// sqlite.db.Exec(`INSERT INTO config VALUES ('gitEmail', '');`)
}

func jsonString(v interface{}) string {
	b, err := json.Marshal(v)
	if isErr(err) {
		return ""
	}
	return string(b)
}
