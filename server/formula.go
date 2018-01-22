package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"sync"
	"time"

	"github.com/gobwas/glob"
	yaml "gopkg.in/yaml.v2"
)

type formulaCI struct {
	// list map[string]
	list      map[string]map[string][]*formula
	pipelines map[string]map[string][]*execPipeline
}

var ci = formulaCI{
	list:      map[string]map[string][]*formula{},
	pipelines: map[string]map[string][]*execPipeline{},
}

type formula struct {
	Repo    string
	Branch  string
	Mode    string
	Name    string
	Setup   []string
	Trigger []formulaTrigger
	Flow    []string
	Steps   map[string]formulaStep
	Deploy  formulaDeploy
	Webhook formulaWebhook
}

type formulaTrigger struct {
	Tag     string
	Changes []string
}

type formulaStep struct {
	Env     map[string]string
	Trigger []formulaTrigger
	Cmd     []string
}

type formulaDeploy struct {
	Target     string
	Kubernetes struct {
		Type          string
		Namespace     string
		Name          string
		ContainerName string `yaml:"containerName"`
		ImageHub      string `yaml:"imageHub"`
		Image         string
		// Image         struct {
		// 	Hub     string
		// 	HubName string
		// 	Name    string
		// 	Tag     string
		// }
	}
}

type formulaWebhook struct {
	Slack string
}

const formulaYAML = ".formulaci.yaml"

func (fci *formulaCI) init() {
	var wg sync.WaitGroup
	for repoName, branchMap := range repo.list {
		for branch := range branchMap {
			n := repoName
			b := branch
			wg.Add(1)
			go func() {
				fci.install(n, b, nil)
				wg.Done()
			}()
		}
	}
	wg.Wait()
}

func (fci *formulaCI) getHistory(data interface{}, resFunc wsResFunc) {
	rows, err := sqlite.query("SELECT result, repo, branch, flow, log, time, dur FROM history ORDER BY rowid DESC;")
	if isErr(err) {
		resFunc(500, "", err.Error())
	}
	defer rows.Close()

	list := []tableHistory{}
	for name, bPipe := range fci.pipelines {
		for branch := range bPipe {
			list = append(list, tableHistory{
				Result: -1,
				Repo:   name,
				Branch: branch,
			})
		}
	}
	for rows.Next() {
		var v tableHistory
		rows.Scan(&v.Result, &v.Repo, &v.Branch, &v.Flow, &v.Log, &v.Time, &v.Dur)
		list = append(list, v)
	}
	resFunc(200, jsonString(list), "")
}

func (fci *formulaCI) install(name string, branch string, pipe *execPipeline) error {
	fmt.Println("Formula installing ", name, branch)
	dir := repoPath + name + "/" + branch + "/"
	// _, err := runCmdPath(dir, "git", "checkout", branch)
	// if isErr(err) {
	// 	return err
	// }
	// _, err = runCmdPath(dir, "git", "pull", "origin", branch)
	// if isErr(err) {
	// 	return err
	// }
	b, err := ioutil.ReadFile(dir + formulaYAML)
	if isErr(err) {
		return err
	}
	bb := bytes.Split(b, []byte("---\n"))
	flist := []*formula{}
	for _, s := range bb {
		var f formula
		err := yaml.Unmarshal(s, &f)
		if isErr(err) {
			return err
		}
		f.Repo = name
		f.Branch = branch

		cmds := []interface{}{}
		for _, setup := range f.Setup {
			cmds = append(cmds, dir, strings.Fields(setup))
		}

		// if pipe != nil {
		// 	pipe.start(cmds)
		// 	_, cancel, err := pipe.wait()
		// 	if cancel {
		// 		return errors.New("pipeline canceled")
		// 	} else if isErr(err) {
		// 		return err
		// 	}
		// } else {
		// 	_, err = cmdEX.runSets(cmds...)
		// 	if isErr(err) {
		// 		return err
		// 	}
		// }

		_, err = cmdEX.runSets(cmds...)
		if isErr(err) {
			return err
		}

		flist = append(flist, &f)
	}

	mutex.Lock()
	if _, ok := fci.list[name]; !ok {
		fci.list[name] = map[string][]*formula{}
	}
	fci.list[name][branch] = flist
	mutex.Unlock()

	return nil
}

func (fci *formulaCI) trigger(repoName string, hookBranch string) {
	if _, ok := repo.list[repoName]; !ok {
		return
	}
	if _, ok := repo.list[repoName][hookBranch]; !ok {
		return
	}

	mutex.Lock()
	if _, ok := fci.pipelines[repoName]; !ok {
		fci.pipelines[repoName] = map[string][]*execPipeline{}
	}
	if pipesBuf, ok := fci.pipelines[repoName][hookBranch]; ok {
		for _, p := range pipesBuf {
			p.cancel()
		}
	}
	pipe := cmdEX.newPipeline()
	fci.pipelines[repoName][hookBranch] = []*execPipeline{pipe}
	mutex.Unlock()

	var (
		res    []executerResult
		cancel bool
		err    error
	)

	defer func() {
		if cancel {
			fmt.Println("trigger canceled : ", repoName, hookBranch)
		} else if isErr(err) {
			fmt.Println("trigger failed : ", repoName, hookBranch)
		} else {
			fmt.Println("trigger done : ", repoName, hookBranch)
		}
	}()

	now := time.Now()

	dir := repoPath + repoName + "/" + hookBranch
	pipe.start(
		dir, []string{"git", "rev-parse", "HEAD"},
		dir, []string{"git", "reset", "--hard", "HEAD"},
		dir, []string{"git", "pull", "origin", hookBranch},
	)
	res, cancel, err = pipe.wait()
	if cancel || isErr(err) {
		return
	}
	prevHead := strings.Replace(res[0].output, "\n", "", -1)

	pipe.start(dir, []string{"git", "diff", "--name-only", "HEAD", prevHead})
	res, cancel, err = pipe.wait()
	if cancel || isErr(err) {
		return
	}

	changes := strings.Split(res[0].output, "\n")

	for _, s := range changes {
		if pipe.stop {
			return
		}

		if strings.Index(s, formulaYAML) != -1 {
			fci.install(repoName, hookBranch, pipe)
			break
		}
	}

	formulas := fci.list[repoName][hookBranch]
	for _, f := range formulas {
		if pipe.stop {
			return
		}
		if fci.validTrigger(&f.Trigger, changes) {
			go ci.run(f, changes, now)
		}
	}

	pipe.cancel()
}

func (fci *formulaCI) validTrigger(triggers *[]formulaTrigger, changes []string) bool {
	if len(*triggers) == 0 {
		return true
	}
	for _, t := range *triggers {
		for _, tc := range t.Changes {
			for _, c := range changes {
				if c == "" {
					continue
				}
				if glob.MustCompile(tc).Match(c) {
					return true
				}
			}
		}
	}
	return false
}

// func (fci *formulaCI) run(repoName string, branch string, i int) {
func (fci *formulaCI) run(f *formula, changes []string, now time.Time) {
	mutex.Lock()
	pipe := cmdEX.newPipeline()
	fci.pipelines[f.Repo][f.Branch] = append(fci.pipelines[f.Repo][f.Branch], pipe)
	mutex.Unlock()

	fmt.Println("New pipeline : ", f.Repo, f.Branch)

	var (
		res    []executerResult
		cancel bool
		err    error
	)

	defer func() {
		if cancel {
			fmt.Println("pipeline canceled : ", f.Repo, f.Branch)
		} else if isErr(err) {
			fmt.Println("pipeline failed : ", f.Repo, f.Branch)
		} else {
			fmt.Println("pipeline done : ", f.Repo, f.Branch)
		}
	}()

	hs := tableHistory{
		Result: -1,
		Repo:   f.Repo,
		Branch: f.Branch,

		flow: []string{},
		log:  map[string]string{},
	}

	dir := repoPath + f.Repo + "/" + f.Branch

	if f.Webhook.Slack != "" {
		go webhooksM.sendToSlack(f.Webhook.Slack, hs)
	}

	for _, fw := range f.Flow {
		if hs.Result != -1 {
			break
		}

		output := []string{}
		if _, ok := f.Steps[fw]; !ok {
			return
		}
		hs.flow = append(hs.flow, fw)
		step := f.Steps[fw]

		if !fci.validTrigger(&step.Trigger, changes) {
			hs.log[fw] = "skip"
			continue
		}

		cmds := []interface{}{}
		for _, cmd := range step.Cmd {
			cmds = append(cmds, dir, strings.Fields(cmd))
		}

		pipe.start(cmds...)
		res, cancel, err = pipe.wait()
		for _, r := range res {
			output = append(output, r.cmd, r.output)
		}

		if cancel {
			hs.Result = -2
		} else if isErr(err) {
			output = append(output, err.Error())
			hs.Result = 0
		}
		hs.log[fw] = strings.Join(output, "\n")
	}

	if hs.Result == -1 && f.Deploy.Target == "kubernetes" {
		k8s := f.Deploy.Kubernetes
		output := []string{}
		pipe.start(
			dir, []string{"docker", "build", "-t", k8s.Image, "."},
			dir, []string{"docker", "tag", k8s.Image, k8s.Image},
			dir, []string{"gcloud", "docker", "--", "push", k8s.Image},
			"", []string{"./kube-update-img.sh", k8s.Namespace, k8s.Type, k8s.Name, k8s.ContainerName, k8s.Image},
		)
		res, cancel, err = pipe.wait()
		for _, r := range res {
			output = append(output, r.cmd, r.output)
		}
		if cancel {
			hs.Result = -2
		} else if isErr(err) {
			output = append(output, err.Error())
			hs.Result = 0
		}
		hs.flow = append(hs.flow, "deploy")
		hs.log["deploy"] = strings.Join(output, "\n")
	}

	if hs.Result == -1 {
		hs.Result = 1
	}

	hs.Flow = jsonString(hs.flow)
	hs.Log = jsonString(hs.log)
	hs.Time = now.Unix()
	hs.Dur = time.Since(now).Seconds()

	mutex.Lock()
	pipe.cancel()
	allStop := true
	for _, p := range fci.pipelines[f.Repo][f.Branch] {
		if !p.stop {
			allStop = false
		}
	}
	if allStop {
		delete(fci.pipelines[f.Repo], f.Branch)
	}
	mutex.Unlock()

	sqlite.exec(
		"INSERT INTO history (result, repo, branch, flow, log, time, dur) VALUES (?, ?, ?, ?, ?, ?, ?)",
		hs.Result, hs.Repo, hs.Branch, hs.Flow, hs.Log, hs.Time, hs.Dur,
	)

	if f.Webhook.Slack != "" {
		go webhooksM.sendToSlack(f.Webhook.Slack, hs)
	}
}
