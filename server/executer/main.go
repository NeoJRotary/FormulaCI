package executer

import (
	"log"
	"os"

	D "github.com/NeoJRotary/describe-go"
)

// Exec ...
type Exec struct {
	doStdLog  bool
	stdLogger *log.Logger
}

// Result ...
type Result struct {
	cmd    string
	output string
}

// RunFunc ...
type RunFunc func(...string) (string, error)

// NewExec ...
func NewExec(doStdLog bool) *Exec {
	return &Exec{
		doStdLog:  doStdLog,
		stdLogger: log.New(os.Stdout, "\n[EXEC] ", log.LstdFlags),
	}
}

// Run ...
func (ex *Exec) Run(args ...string) (out string, err error) {
	cmd := ex.NewCMD("", args...)
	out, err = cmd.Run()
	if ex.doStdLog {
		if D.IsErr(err) {
			ex.stdLogger.Println(cmd.str, "\n", out, err.Error())
		} else {
			ex.stdLogger.Println(cmd.str, "\n", out)
		}
	}
	return out, err
}

// RunDir ...
func (ex *Exec) RunDir(dir string) RunFunc {
	return func(args ...string) (out string, err error) {
		cmd := ex.NewCMD(dir, args...)
		out, err = cmd.Run()
		if ex.doStdLog {
			if D.IsErr(err) {
				ex.stdLogger.Println(cmd.str, "\n", out, err.Error())
			} else {
				ex.stdLogger.Println(cmd.str, "\n", out)
			}
		}
		return out, err
	}
}

// RunSets ...
func (ex *Exec) RunSets(v ...interface{}) (result []Result, err error) {
	for i := 0; i < len(v); i += 2 {
		dir := v[i].(string)
		args := v[i+1].([]string)
		cmd := ex.NewCMD(dir, args...)
		out, err := cmd.Run()
		if ex.doStdLog {
			if D.IsErr(err) {
				ex.stdLogger.Println(cmd.str, "\n", out, err.Error())
			} else {
				ex.stdLogger.Println(cmd.str, "\n", out)
			}
		}
		// outputs = append(outputs, o)
		result = append(result, Result{
			cmd:    cmd.str,
			output: out,
		})
		if D.IsErr(err) {
			return result, err
		}
	}
	return result, nil
}
