package formula

import "../executer"

type processor struct {
	master       *Master
	pipeGroupMap map[string]*executer.PipelineGroup
}

func (fma *Master) initProcessor() {
	fma.processor = &processor{
		master:       fma,
		pipeGroupMap: map[string]*executer.PipelineGroup{},
	}
}

func (pcr *processor) newPipelineGroup(repoName, branch string) *executer.PipelineGroup {
	key := repoName + ":" + branch
	if current, ok := pcr.pipeGroupMap[key]; ok {
		current.Cancel()
	}
	pcs := pcr.master.execEX.NewPipelineGroup()
	mutex.Lock()
	pcr.pipeGroupMap[key] = pcs
	mutex.Unlock()
	return pcs
}

func (pcr *processor) lifeCheck(repoName, branch string) {
	key := repoName + ":" + branch
	if current, ok := pcr.pipeGroupMap[key]; ok {
		if current.Stop {
			delete(pcr.pipeGroupMap, key)
		}
	}
}

// func (prc *process) newTask() *processTask {
// 	pipe := prc.pipelinegroup.NewPipeline()
// }

// func (prc *process) cancel(repoName, branch string) {
// 	child.PipelineGroup.Cancel()
// 	mutex.Lock()
// 	pipe.Cancel()
// 	allStop := true
// 	for _, p := range fma.Pipelines[repoName][hookBranch] {
// 		if !p.Stop {
// 			allStop = false
// 		}
// 	}
// 	if allStop {
// 		delete(fma.Pipelines[repoName], hookBranch)
// 	}
// 	mutex.Unlock()
// }
