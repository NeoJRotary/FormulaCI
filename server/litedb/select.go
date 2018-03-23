package litedb

import (
	D "github.com/NeoJRotary/describe-go"
)

// GetAllRepo ...
func (db *DB) GetAllRepo() ([]TableRepo, error) {
	rows, err := db.Query("SELECT dir, src, hub, user, name, branch FROM repo;")
	if D.IsErr(err) {
		return nil, err
	}
	defer rows.Close()
	list := []TableRepo{}
	for rows.Next() {
		var r TableRepo
		rows.Scan(&r.Dir, &r.Src, &r.Hub, &r.User, &r.Name, &r.Branch)
		list = append(list, r)
	}
	return list, nil
}

// GetAllConfig ...
func (db *DB) GetAllConfig() ([]TableConfig, error) {
	rows, err := db.Query("SELECT key, value FROM config;")
	if D.IsErr(err) {
		return nil, err
	}
	defer rows.Close()
	list := []TableConfig{}
	for rows.Next() {
		var r TableConfig
		rows.Scan(&r.Key, &r.Value)
		list = append(list, r)
	}
	return list, nil
}

// GetConfig ...
func (db *DB) GetConfig(key string) string {
	row := db.QueryRow("SELECT value FROM config WHERE key=?;", key)

	val := ""
	err := row.Scan(&val)
	if D.IsErr(err) {
		return ""
	}
	return val
}
