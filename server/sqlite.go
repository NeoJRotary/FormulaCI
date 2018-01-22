package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type sqliteDB struct {
	db *sql.DB
}

type tableHistory struct {
	Result int     `json:"result"`
	Repo   string  `json:"repo"`
	Branch string  `json:"branch"`
	Flow   string  `json:"flow"`
	Log    string  `json:"log"`
	Time   int64   `json:"time"`
	Dur    float64 `json:"dur"`

	flow []string
	log  map[string]string
}

var sqlite = sqliteDB{}

func (sq *sqliteDB) query(query string, args ...interface{}) (*sql.Rows, error) {
	return sq.db.Query(query, args...)
}

func (sq *sqliteDB) queryrow(query string, args ...interface{}) *sql.Row {
	return sq.db.QueryRow(query, args...)
}

func (sq *sqliteDB) exec(query string, args ...interface{}) (sql.Result, error) {
	return sq.db.Exec(query, args...)
}

func (*sqliteDB) connect() {
	db, err := sql.Open("sqlite3", dataPath+"sqlite.db")
	if err != nil {
		log.Fatalln(err)
	}
	// defer db.Close()
	sqlite.db = db

	sqlite.initTables()

	rows, _ := db.Query("SELECT * FROM config;")
	// rows.Columns()
	defer rows.Close()
	for rows.Next() {
		var key string
		var val string
		rows.Scan(&key, &val)
		switch key {
		case "git/email":
			git.email = val
			cmdEX.run("git", "config", "--global", "user.email", val)
		case "git/webhookToken":
			git.webhookToken = val
		case "gcloud/project":
			gcloud.project = val
			cmdEX.run("gcloud", "config", "set", "project", val)
		case "gcloud/gkeZone":
			gcloud.gkeZone = val
		case "gcloud/gkeName":
			gcloud.gkeName = val
			// case "gitlabEmail":
			// 	git.gitlabEmail = val
			// case "githubEmail":
			// 	git.githubEmail = val
		}
	}

	if gcloud.gkeName != "" {
		cmdEX.run("gcloud", "container", "clusters", "get-credentials", gcloud.gkeName, "--zone", gcloud.gkeZone, "--project", gcloud.project)
	}
}

func (sq *sqliteDB) initTables() {
	sq.exec("CREATE TABLE config (key text unique, value text);")
	sq.exec("CREATE TABLE repo (name text, src text, branch text);")
	sq.exec(`CREATE TABLE history (
		result integer,
		repo text,
		branch text,
		flow text,
		log text,
		time integer,
		dur real
		);`)
}
