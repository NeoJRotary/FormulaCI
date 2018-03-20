package executer

import (
	"fmt"

	D "github.com/NeoJRotary/describe-go"
)

// Pipeline ...
type Pipeline struct {
	ex       *Exec
	stopChan chan bool
	Stop     bool
	cmds     []*Cmd
	group    *PipelineGroup
}

// NewPipeline ...
func (ex *Exec) NewPipeline() *Pipeline {
	pipe := Pipeline{
		stopChan: make(chan bool, 1),
		Stop:     false,
		ex:       ex,
		// cmds:     []*Cmd{},
	}
	return &pipe
}

// Start ...
func (pipe *Pipeline) Start(v ...interface{}) {
	pipe.cmds = []*Cmd{}
	for i := 0; i < len(v); i += 2 {
		dir := v[i].(string)
		args := v[i+1].([]string)
		cmd := pipe.ex.NewCMD(dir, args...)
		pipe.cmds = append(pipe.cmds, cmd)
	}
}

// Cancel ...
func (pipe *Pipeline) Cancel() {
	if pipe.Stop {
		return
	}
	pipe.stopChan <- true
	pipe.stop()
}

func (pipe *Pipeline) stop() {
	pipe.Stop = true
	if pipe.group != nil {
		pipe.group.childStop()
	}
}

// Wait ...
func (pipe *Pipeline) Wait() (result []Result, cancel bool, err error) {
	if pipe.Stop {
		return nil, true, nil
	}
	for _, c := range pipe.cmds {
		if pipe.ex.doStdLog {
			fmt.Println(c.str)
		}
		c.Start()
		go c.Wait()
		select {
		case <-pipe.stopChan:
			c.Cancel()
			return nil, true, nil
		case err = <-c.done:
			result = append(result, Result{
				cmd:    c.str,
				output: c.stdout.String(),
			})
			if pipe.ex.doStdLog {
				if D.IsErr(err) {
					pipe.ex.stdLogger.Println(c.stdout.String(), err.Error())
				} else {
					pipe.ex.stdLogger.Println(c.stdout.String())
				}
			}
			if D.IsErr(err) {
				return result, false, err
			}
		}
	}
	pipe.stop()
	return result, false, nil
}
