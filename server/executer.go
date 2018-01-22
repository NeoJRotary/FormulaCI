package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type executer struct {
	log    bool
	logger *log.Logger
}

type execPipeline struct {
	ex       *executer
	stopChan chan bool
	stop     bool
	cmds     []*executerCMD
}

type executerCMD struct {
	cmd    *exec.Cmd
	str    string
	stdout *bytes.Buffer
	stderr *bytes.Buffer
	done   chan error
}

type executerResult struct {
	cmd    string
	output string
}

var cmdEX = executer{
	log:    true,
	logger: log.New(os.Stdout, "\n[EXEC] ", log.LstdFlags),
}

func (excmd *executerCMD) run() (string, error) {
	err := excmd.cmd.Run()
	if isErr(err) {
		// errMsg := err.Error() + ": " + stderr.String()
		return "", errors.New(excmd.stderr.String())
	}
	return excmd.stdout.String(), nil
}

func (excmd *executerCMD) start() {
	excmd.done = make(chan error, 1)
	err := excmd.cmd.Start()
	if isErr(err) {
		excmd.done <- err
		return
	}
}

func (excmd *executerCMD) wait() {
	excmd.done <- excmd.cmd.Wait()
}

func (excmd *executerCMD) output() string {
	return excmd.str + "\n" + excmd.stdout.String()
}

func (excmd *executerCMD) cancel() {
	excmd.cmd.Process.Kill()
}

func (*executer) newCMD(dir string, args ...string) *executerCMD {
	cmd := exec.Command(args[0], args[1:]...)
	if dir == "" {
		dir = "./"
	}
	cmd.Dir = dir
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	return &executerCMD{
		cmd:    cmd,
		str:    "(" + dir + ") > " + strings.Join(args, " "),
		stdout: &stdout,
		stderr: &stderr,
	}
}

func (ex *executer) run(args ...string) (out string, err error) {
	cmd := ex.newCMD("", args...)
	out, err = cmd.run()
	if ex.log {
		if isErr(err) {
			ex.logger.Println(cmd.str, "\n", out, err.Error())
		} else {
			ex.logger.Println(cmd.str, "\n", out)
		}
	}
	return out, err
}

func (ex *executer) runDir(dir string, args ...string) (out string, err error) {
	cmd := ex.newCMD(dir, args...)
	out, err = cmd.run()
	if ex.log {
		if isErr(err) {
			ex.logger.Println(cmd.str, "\n", out, err.Error())
		} else {
			ex.logger.Println(cmd.str, "\n", out)
		}
	}
	return out, err
}

func (ex *executer) runSets(v ...interface{}) (result []executerResult, err error) {
	for i := 0; i < len(v); i += 2 {
		dir := v[i].(string)
		args := v[i+1].([]string)
		cmd := ex.newCMD(dir, args...)
		out, err := cmd.run()
		if ex.log {
			if isErr(err) {
				ex.logger.Println(cmd.str, "\n", out, err.Error())
			} else {
				ex.logger.Println(cmd.str, "\n", out)
			}
		}
		// outputs = append(outputs, o)
		result = append(result, executerResult{
			cmd:    cmd.str,
			output: out,
		})
		if isErr(err) {
			return result, err
		}
	}
	return result, nil
}

func (ex *executer) newPipeline() *execPipeline {
	pipe := execPipeline{
		stopChan: make(chan bool, 1),
		stop:     false,
		ex:       ex,
		// cmds:     []*executerCMD{},
	}
	return &pipe
}

func (pipe *execPipeline) start(v ...interface{}) {
	pipe.cmds = []*executerCMD{}
	for i := 0; i < len(v); i += 2 {
		dir := v[i].(string)
		args := v[i+1].([]string)
		cmd := pipe.ex.newCMD(dir, args...)
		pipe.cmds = append(pipe.cmds, cmd)
	}
}

func (pipe *execPipeline) cancel() {
	if pipe.stop {
		return
	}
	pipe.stopChan <- true
	pipe.stop = true
}

func (pipe *execPipeline) wait() (result []executerResult, cancel bool, err error) {
	if pipe.stop {
		return nil, true, nil
	}
	for _, c := range pipe.cmds {
		if pipe.ex.log {
			fmt.Println(c.str)
		}
		c.start()
		go c.wait()
		select {
		case <-pipe.stopChan:
			c.cancel()
			return nil, true, nil
		case err = <-c.done:
			result = append(result, executerResult{
				cmd:    c.str,
				output: c.stdout.String(),
			})
			if pipe.ex.log {
				if isErr(err) {
					pipe.ex.logger.Println(c.stdout.String(), err.Error())
				} else {
					pipe.ex.logger.Println(c.stdout.String())
				}
			}
			if isErr(err) {
				return result, false, err
			}
		}
	}
	return result, false, nil
}
