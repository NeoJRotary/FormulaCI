package litedb

import (
	D "github.com/NeoJRotary/describe-go"
)

// GetAllRepo ...
func (db *DB) GetAllRepo() ([]TableRepo, error) {
	rows, err := db.Query("SELECT * FROM repo;")
	if D.IsErr(err) {
		return nil, err
	}
	defer rows.Close()
	list := []TableRepo{}
	for rows.Next() {
		var r TableRepo
		rows.Scan(&r.Name, &r.Src, &r.Branch)
		list = append(list, r)
	}
	return list, nil
}
