package litedb

// TableConfig ...
type TableConfig struct {
	Key   string
	Value string
}

// TableRepo ...
type TableRepo struct {
	Dir    string
	Src    string
	Hub    string
	User   string
	Name   string
	Branch string
}

// TableHistory ...
type TableHistory struct {
	Result int     `json:"result"`
	Repo   string  `json:"repo"`
	Branch string  `json:"branch"`
	Flow   string  `json:"flow"`
	Log    string  `json:"log"`
	Time   int64   `json:"time"`
	Dur    float64 `json:"dur"`
	RowID  int64   `json:"rowid"`

	flow []string
	log  map[string]string
}

// InitTables ...
func (db *DB) InitTables() {
	db.Exec("CREATE TABLE config (key text unique, value text);")
	db.Exec(`CREATE TABLE repo (
		dir text,
		src text,
		hub text,
		user text,
		name text,
		branch text
		);`)
	db.Exec(`CREATE TABLE history (
		result integer,
		repo text,
		branch text,
		flow text,
		log text,
		time integer,
		dur real
		);`)
}
