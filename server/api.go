package main

// func (rpm *Manager) getList(data interface{}, resFunc wsResFunc) {
// 	list := []repository{}
// 	for _, b := range rpm.list {
// 		for _, r := range b {
// 			list = append(list, *r)
// 		}
// 	}

// 	b, err := json.Marshal(list)
// 	if isErr(err) {
// 		resFunc(500, "", err.Error())
// 	}
// 	resFunc(200, string(b), "")
// }

// func (rpm *Manager) add(data interface{}, resFunc wsResFunc) {
// 	d := data.(map[string]interface{})
// 	name := d["name"].(string)
// 	src := d["src"].(string)
// 	branch := d["branch"].(string)
// 	err := sqlite.queryrow("SELECT name FROM repo WHERE name=? AND branch=?;", name, branch).Scan()
// 	if !isErr(err) {
// 		resFunc(400, "", err.Error())
// 		return
// 	} else if err != sql.ErrNoRows {
// 		resFunc(400, "", "same name:branch already exist")
// 		return
// 	}

// 	dir := repoPath + name + "/" + branch

// 	err = git.pullRepo(name, branch, src)
// 	if isErr(err) {
// 		resFunc(500, "", err.Error())
// 		return
// 	}

// 	err = ci.install(name, branch)
// 	if isErr(err) {
// 		resFunc(400, "", err.Error())
// 		cmdEX.run("rm", "-rf", dir)
// 		return
// 	}

// 	_, err = sqlite.exec("INSERT INTO repo (name, src, branch) VALUES (?, ?, ?);", name, src, branch)
// 	if isErr(err) {
// 		resFunc(500, "", err.Error())
// 		return
// 	}

// 	atI := strings.Index(src, "@")
// 	colonI := strings.Index(src[atI:], ":")
// 	hub := src[atI+1 : atI+colonI]
// 	mutex.Lock()
// 	if _, ok := rpm.list[name]; !ok {
// 		rpm.list[name] = map[string]*repository{}
// 	}
// 	repo := repository{
// 		Name:   name,
// 		Src:    src,
// 		Branch: branch,
// 	}
// 	rpm.list[name][branch] = &repo
// 	rpm.hub[hub] = append(rpm.hub[hub], &repo)
// 	mutex.Unlock()
// 	resFunc(200, "", "")
// }

// func (rpm *Manager) remove(data interface{}, resFunc wsResFunc) {
// 	d := data.(map[string]interface{})
// 	name := d["name"].(string)
// 	branch := d["branch"].(string)

// 	_, err := sqlite.exec("DELETE FROM repo WHERE name=? AND branch=?;", name, branch)
// 	if isErr(err) {
// 		resFunc(500, "", err.Error())
// 		return
// 	}

// 	mutex.Lock()
// 	delete(repo.list[name], branch)
// 	delete(ci.list[name], branch)
// 	mutex.Unlock()

// 	cmdEX.run("rm", "-rf", repoPath+name+"/"+branch)
// 	resFunc(200, "", "")
// }

// func (rpm *Manager) trigger(data interface{}, resFunc wsResFunc) {
// 	d := data.(map[string]interface{})
// 	name := d["name"].(string)
// 	branch := d["branch"].(string)
// 	ci.trigger(name, branch)
// }

// func (fma *Master) getHistory(data interface{}, resFunc wsResFunc) {
// 	rows, err := sqlite.query("SELECT result, repo, branch, flow, log, time, dur, rowid FROM history ORDER BY rowid DESC LIMIT 20;")
// 	if isErr(err) {
// 		resFunc(500, "", err.Error())
// 	}
// 	defer rows.Close()

// 	list := []tableHistory{}
// 	for name, bPipe := range fci.pipelines {
// 		for branch := range bPipe {
// 			list = append(list, tableHistory{
// 				Result: -1,
// 				Repo:   name,
// 				Branch: branch,
// 			})
// 		}
// 	}
// 	for rows.Next() {
// 		var v tableHistory
// 		rows.Scan(&v.Result, &v.Repo, &v.Branch, &v.Flow, &v.Log, &v.Time, &v.Dur, &v.RowID)
// 		list = append(list, v)
// 	}
// 	resFunc(200, jsonString(list), "")
// }
