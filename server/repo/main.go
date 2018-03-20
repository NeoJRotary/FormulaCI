package repo

import (
	"log"
	"strings"
	"sync"

	"../litedb"
	D "github.com/NeoJRotary/describe-go"
)

// Repository ...
type Repository struct {
	Name   string `json:"name"`
	Src    string `json:"src"`
	Branch string `json:"branch"`
}

// Manager ...
type Manager struct {
	RepoPath string
	db       *litedb.DB
	List     map[string]map[string]*Repository
	Hub      map[string][]*Repository
}

var mutex = &sync.Mutex{}

// Init ...
func Init(repoPath string, db *litedb.DB) (*Manager, error) {
	list, err := db.GetAllRepo()
	if D.IsErr(err) {
		log.Fatalln(err)
	}

	rpm := Manager{
		RepoPath: repoPath,
		db:       db,
		List:     map[string]map[string]*Repository{},
		Hub:      map[string][]*Repository{},
	}

	for _, l := range list {
		r := Repository{
			Name:   l.Name,
			Src:    l.Src,
			Branch: l.Branch,
		}
		atI := strings.Index(r.Src, "@")
		colonI := strings.Index(r.Src[atI:], ":")
		hub := r.Src[atI+1 : atI+colonI]
		mutex.Lock()
		if _, ok := rpm.List[r.Name]; !ok {
			rpm.List[r.Name] = map[string]*Repository{}
		}
		rpm.List[r.Name][r.Branch] = &r
		rpm.Hub[hub] = append(rpm.Hub[hub], &r)
		mutex.Unlock()
	}

	return &rpm, nil
}

// HasRepo ...
func (man *Manager) HasRepo(name string, branch string) bool {
	if _, ok := man.List[name]; !ok {
		return false
	}
	if _, ok := man.List[name][branch]; !ok {
		return false
	}
	return true
}
