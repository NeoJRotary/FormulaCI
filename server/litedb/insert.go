package litedb

// InsertRepo ...
func (db *DB) InsertRepo(dir, src, hub, user, name, branch string) (row TableRepo, err error) {
	query := "INSERT INTO repo (dir, src, hub, user, name, branch) VALUES (?, ?, ?, ?, ?, ?);"
	args := []interface{}{dir, src, hub, user, name, branch}
	_, err = db.Query(query, args...)
	row = TableRepo{
		Dir:    dir,
		Src:    src,
		Hub:    hub,
		User:   user,
		Name:   name,
		Branch: branch,
	}
	return row, err
}
