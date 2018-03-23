package repo

import (
	"path"
	"sync"

	"../executer"
	"../git"
	"../litedb"
	D "github.com/NeoJRotary/describe-go"
)

// Repository ...
type Repository struct {
	litedb.TableRepo
	Git *git.Repository
}

// Manager ...
type Manager struct {
	RepoPath string
	db       *litedb.DB
	execEX   *executer.Exec
	git      *git.Agent
	Map      map[string]*Repository
	List     []*Repository
}

var mutex = &sync.Mutex{}

// Init ...
func Init(repoPath string, db *litedb.DB, execEX *executer.Exec, gitA *git.Agent) (man *Manager, err error) {
	defer D.RecoverErr(nil)

	list, err := db.GetAllRepo()
	D.CheckErr(err)

	man = &Manager{
		RepoPath: repoPath,
		git:      gitA,
		db:       db,
		execEX:   execEX,
		Map:      map[string]*Repository{},
		List:     []*Repository{},
	}

	for _, l := range list {
		repo, err := gitA.NewRepository(l.Dir, "origin", l.Src, l.Branch)
		D.CheckErr(err)

		man.AddRepo(l, repo)
	}

	return man, nil
}

// GetListKey ...
func GetListKey(hub, user, name, branch string) string {
	return D.StringSlice([]string{hub, user, name, branch}).Join("/").Get()
}

// AddRepo ...
func (man *Manager) AddRepo(row litedb.TableRepo, repo *git.Repository) {
	r := &Repository{row, repo}
	key := GetListKey(row.Hub, row.User, row.Name, row.Branch)
	mutex.Lock()
	man.List = append(man.List, r)
	man.Map[key] = r
	mutex.Unlock()
}

// HasRepo ...
func (man *Manager) HasRepo(key string) bool {
	_, ok := man.Map[key]
	return ok
}

// NewRepo ...
func (man *Manager) NewRepo(src, branch string) (err error) {
	dir := path.Join(man.RepoPath, UUID())
	defer D.RecoverErr(func(err error) {
		man.execEX.Run("rm", "-rf", dir)
	})

	_, err = man.execEX.Run("mkdir", "-p", dir)
	D.CheckErr(err)

	hub, user, name, err := ResolveSrc(src)
	D.CheckErr(err)

	repo, err := man.git.NewRepository(dir, "origin", src, branch)
	D.CheckErr(err)

	row, err := man.db.InsertRepo(dir, src, hub, user, name, branch)
	D.CheckErr(err)

	man.AddRepo(row, repo)

	return nil
}
