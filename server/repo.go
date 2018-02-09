package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"strings"
)

type repository struct {
	Name   string `json:"name"`
	Src    string `json:"src"`
	Branch string `json:"branch"`
}

type repositories struct {
	list map[string]map[string]*repository
	hub  map[string][]*repository
}

var repo = repositories{
	list: map[string]map[string]*repository{},
	hub:  map[string][]*repository{},
}
var repoPath = dataPath + "repo/"

func (rp *repositories) init() {
	rows, err := sqlite.query("SELECT * FROM repo;")
	if isErr(err) {
		log.Fatalln(err)
	}
	defer rows.Close()
	for rows.Next() {
		var r repository
		rows.Scan(&r.Name, &r.Src, &r.Branch)
		atI := strings.Index(r.Src, "@")
		colonI := strings.Index(r.Src[atI:], ":")
		hub := r.Src[atI+1 : atI+colonI]
		mutex.Lock()
		if _, ok := rp.list[r.Name]; !ok {
			rp.list[r.Name] = map[string]*repository{}
		}
		rp.list[r.Name][r.Branch] = &r
		rp.hub[hub] = append(rp.hub[hub], &r)
		mutex.Unlock()
	}
}

func (rp *repositories) getList(data interface{}, resFunc wsResFunc) {
	list := []repository{}
	for _, b := range rp.list {
		for _, r := range b {
			list = append(list, *r)
		}
	}

	b, err := json.Marshal(list)
	if isErr(err) {
		resFunc(500, "", err.Error())
	}
	resFunc(200, string(b), "")
}

func (rp *repositories) add(data interface{}, resFunc wsResFunc) {
	d := data.(map[string]interface{})
	name := d["name"].(string)
	src := d["src"].(string)
	branch := d["branch"].(string)
	err := sqlite.queryrow("SELECT name FROM repo WHERE name=? AND branch=?;", name, branch).Scan()
	if !isErr(err) {
		resFunc(400, "", err.Error())
		return
	} else if err != sql.ErrNoRows {
		resFunc(400, "", "same name:branch already exist")
		return
	}

	dir := repoPath + name + "/" + branch

	err = git.pullRepo(name, branch, src)
	if isErr(err) {
		resFunc(500, "", err.Error())
		return
	}

	err = ci.install(name, branch)
	if isErr(err) {
		resFunc(400, "", err.Error())
		cmdEX.run("rm", "-rf", dir)
		return
	}

	_, err = sqlite.exec("INSERT INTO repo (name, src, branch) VALUES (?, ?, ?);", name, src, branch)
	if isErr(err) {
		resFunc(500, "", err.Error())
		return
	}

	atI := strings.Index(src, "@")
	colonI := strings.Index(src[atI:], ":")
	hub := src[atI+1 : atI+colonI]
	mutex.Lock()
	if _, ok := rp.list[name]; !ok {
		rp.list[name] = map[string]*repository{}
	}
	repo := repository{
		Name:   name,
		Src:    src,
		Branch: branch,
	}
	rp.list[name][branch] = &repo
	rp.hub[hub] = append(rp.hub[hub], &repo)
	mutex.Unlock()
	resFunc(200, "", "")
}

func (rp *repositories) remove(data interface{}, resFunc wsResFunc) {
	d := data.(map[string]interface{})
	name := d["name"].(string)
	branch := d["branch"].(string)

	_, err := sqlite.exec("DELETE FROM repo WHERE name=? AND branch=?;", name, branch)
	if isErr(err) {
		resFunc(500, "", err.Error())
		return
	}

	mutex.Lock()
	delete(repo.list[name], branch)
	delete(ci.list[name], branch)
	mutex.Unlock()

	cmdEX.run("rm", "-rf", repoPath+name+"/"+branch)
	resFunc(200, "", "")
}

func (rp *repositories) trigger(data interface{}, resFunc wsResFunc) {
	d := data.(map[string]interface{})
	name := d["name"].(string)
	branch := d["branch"].(string)
	ci.trigger(name, branch)
}
