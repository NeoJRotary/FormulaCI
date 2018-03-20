package formula

import (
	"log"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"../executer"
	"../repo"
	D "github.com/NeoJRotary/describe-go"
	"github.com/gobwas/glob"
)

// Master ...
type Master struct {
	execEX    *executer.Exec
	repoM     *repo.Manager
	YAMLPath  string
	List      map[string]map[string][]*Conf
	processor *processor
}

var mutex = &sync.Mutex{}

// Init ...
func Init(yamlPath string, repoM *repo.Manager, execEX *executer.Exec) (*Master, error) {
	fma := Master{
		repoM:    repoM,
		YAMLPath: yamlPath,
		List:     map[string]map[string][]*Conf{},
	}
	fma.initProcessor()

	var wg sync.WaitGroup
	for repoName, branchMap := range repoM.List {
		for branch := range branchMap {
			n := repoName
			b := branch
			wg.Add(1)
			go func() {
				fma.Install(n, b)
				wg.Done()
			}()
		}
	}
	wg.Wait()

	return &fma, nil
}

// Install ...
func (fma *Master) Install(repoName string, branch string) error {
	log.Println("Formula installing ", repoName, branch)
	repoM := fma.repoM
	execEX := fma.execEX

	dir := repoM.RepoPath + repoName + "/" + branch + "/"
	confs, err := GetConfFromYAML(dir + fma.YAMLPath)
	if D.IsErr(err) {
		return err
	}

	for _, f := range confs {
		if f.Branch != branch {
			continue
		}

		f.Repo = repoName

		cmds := []interface{}{}
		for _, setup := range f.Setup {
			cmds = append(cmds, dir, strings.Fields(setup))
		}

		_, err = execEX.RunSets(cmds...)
		if D.IsErr(err) {
			return err
		}
	}

	mutex.Lock()
	if _, ok := fma.List[repoName]; !ok {
		fma.List[repoName] = map[string][]*Conf{}
	}
	fma.List[repoName][branch] = confs
	mutex.Unlock()

	return nil
}

// Trigger ...
func (fma *Master) Trigger(repoName string, hookBranch string) {
	execEX := fma.execEX
	repoM := fma.repoM

	if !repoM.HasRepo(repoName, hookBranch) {
		return
	}

	pipeGroup := fma.processor.newPipelineGroup(repoName, hookBranch)
	pipe := pipeGroup.NewPipeline()

	var (
		result []executer.Result
		cancel bool
		err    error
	)

	defer func() {
		if cancel {
			log.Println("trigger canceled : ", repoName, hookBranch)
		} else if D.IsErr(err) {
			log.Println("trigger failed : ", repoName, hookBranch)
		} else {
			log.Println("trigger done : ", repoName, hookBranch)
		}

		fma.processor.lifeCheck(repoName, hookBranch)
	}()

	now := time.Now()

	dir := repoM.RepoPath + repoName + "/" + hookBranch
	pipe.Start(
		dir, []string{"git", "rev-parse", "HEAD"},
		dir, []string{"git", "reset", "--hard", "HEAD"},
		dir, []string{"git", "clean", "-df"},
		dir, []string{"git", "pull", "origin", hookBranch},
	)
	res, cancel, err = pipe.wait()
	if cancel || D.IsErr(err) {
		return
	}
	prevHead := strings.Replace(res[0].output, "\n", "", -1)

	pipe.start(dir, []string{"git", "diff", "--name-only", "HEAD", prevHead})
	res, cancel, err = pipe.wait()
	if cancel || D.IsErr(err) {
		return
	}

	changes := strings.Split(res[0].output, "\n")
	for _, s := range changes {
		if pipe.stop {
			return
		}

		if strings.Index(s, formulaYAML) != -1 {
			fma.install(repoName, hookBranch)
			break
		}
	}

	for _, f := range fma.List[repoName][hookBranch] {
		if pipe.stop {
			return
		}
		ls := fma.validTrigger(&f.Trigger, changes)

		if len(ls) != 0 || ls == nil {
			go ci.run(f, changes, ls, now)
		}
	}
}

func (fma *Master) validTrigger(triggers *[]formulaTrigger, changes []string) []*formulaTrigger {
	if len(*triggers) == 0 {
		return nil
	}

	ls := []*formulaTrigger{}
	for _, t := range *triggers {
		for _, tc := range t.Changes {
			if len(changes) == 0 {
				if glob.MustCompile(tc).Match("") {
					ls = append(ls, &t)
				}
			} else {
				for _, c := range changes {
					// if c == "" {
					// 	continue
					// }
					if glob.MustCompile(tc).Match(c) {
						ls = append(ls, &t)
					}
				}
			}
		}
	}
	return ls
}

func (fma *Master) replaceTriggerVar(triggers []*formulaTrigger, cmd string) string {
	for _, t := range triggers {
		if t.Name != "" {
			cmd = strings.Replace(cmd, "${"+t.Name+"}", t.Value, -1)
		}
	}

	return cmd
}

// func (fma *Master) run(repoName string, branch string, i int) {
func (fma *Master) run(f *formula, changes []string, globalTriggers []*formulaTrigger, now time.Time) {
	mutex.Lock()
	pipe := cmdEX.newPipeline()
	fma.Pipelines[f.Repo][f.Branch] = append(fma.Pipelines[f.Repo][f.Branch], pipe)
	mutex.Unlock()

	log.Println("New pipeline : ", f.Repo, f.Branch)

	var (
		res    []executerResult
		cancel bool
		err    error
	)

	defer func() {
		if cancel {
			log.Println("pipeline canceled : ", f.Repo, f.Branch)
		} else if D.IsErr(err) {
			log.Println("pipeline failed : ", f.Repo, f.Branch)
		} else {
			log.Println("pipeline done : ", f.Repo, f.Branch)
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

	for _, hook := range f.Webhooks {
		go webhooksM.send(hook, hs, globalTriggers)
	}

	for _, fw := range f.Flow {
		if hs.Result != -1 {
			break
		}

		output := []string{}
		if _, ok := f.Steps[fw]; !ok {
			hs.log[fw] = "has flow but step is not defined"
			continue
		}
		hs.flow = append(hs.flow, fw)
		step := f.Steps[fw]

		ls := fma.validTrigger(&step.Trigger, changes)
		if len(ls) == 0 && ls != nil {
			hs.log[fw] = "skip"
			continue
		}

		triggers := append([]*formulaTrigger{}, globalTriggers...)
		triggers = append(triggers, ls...)

		cmds := []interface{}{}
		for _, cmd := range step.Cmd {
			cmd = fma.replaceTriggerVar(triggers, cmd)
			cmds = append(cmds, dir, strings.Fields(cmd))
		}

		pipe.start(cmds...)
		res, cancel, err = pipe.wait()
		for _, r := range res {
			output = append(output, r.cmd, r.output)
		}

		if cancel {
			hs.Result = -2
		} else if D.IsErr(err) {
			output = append(output, err.Error())
			hs.Result = 0
		}
		hs.log[fw] = strings.Join(output, "\n")
	}

	if hs.Result == -1 && f.Deploy.Target == "kubernetes" {
		k8s := f.Deploy.Kubernetes
		output := []string{}
		name := fma.replaceTriggerVar(globalTriggers, k8s.Name)
		containerName := fma.replaceTriggerVar(globalTriggers, k8s.ContainerName)
		image := fma.replaceTriggerVar(globalTriggers, k8s.Image)
		dockerPath := filepath.Join(dir, fma.replaceTriggerVar(globalTriggers, f.Deploy.Path))
		pipe.start(
			dockerPath, []string{"docker", "build", "-t", image, "."},
			dir, []string{"docker", "tag", image, image},
			dir, []string{"gcloud", "docker", "--", "push", image},
			"", []string{"./kube-update-img.sh", k8s.Namespace, k8s.Type, name, containerName, image},
		)
		res, cancel, err = pipe.wait()
		for _, r := range res {
			output = append(output, r.cmd, r.output)
		}
		if cancel {
			hs.Result = -2
		} else if D.IsErr(err) {
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
	for _, p := range fma.Pipelines[f.Repo][f.Branch] {
		if !p.stop {
			allStop = false
		}
	}
	if allStop {
		delete(fma.Pipelines[f.Repo], f.Branch)
	}
	mutex.Unlock()

	sqlite.exec(
		"INSERT INTO history (result, repo, branch, flow, log, time, dur) VALUES (?, ?, ?, ?, ?, ?, ?)",
		hs.Result, hs.Repo, hs.Branch, hs.Flow, hs.Log, hs.Time, hs.Dur,
	)

	for _, hook := range f.Webhooks {
		go webhooksM.send(hook, hs, globalTriggers)
	}
}
