package executer

import (
	"bytes"
	"os/exec"
	"strings"

	D "github.com/NeoJRotary/describe-go"
)

// Cmd ...
type Cmd struct {
	cmd    *exec.Cmd
	str    string
	stdout *bytes.Buffer
	stderr *bytes.Buffer
	done   chan error
}

// NewCMD ...
func (*Exec) NewCMD(dir string, args ...string) *Cmd {
	cmd := exec.Command(args[0], args[1:]...)
	if dir == "" {
		dir = "./"
	}
	cmd.Dir = dir
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	return &Cmd{
		cmd:    cmd,
		str:    "(" + dir + ") > " + strings.Join(args, " "),
		stdout: &stdout,
		stderr: &stderr,
	}
}

// Run ...
func (c *Cmd) Run() (string, error) {
	err := c.cmd.Run()
	if D.IsErr(err) {
		// errMsg := err.Error() + ": " + stderr.String()
		return "", D.NewErr(c.stderr.String())
	}
	return c.stdout.String(), nil
}

// Start ...
func (c *Cmd) Start() {
	c.done = make(chan error, 1)
	err := c.cmd.Start()
	if D.IsErr(err) {
		c.done <- err
		return
	}
}

// Wait ...
func (c *Cmd) Wait() {
	err := c.cmd.Wait()
	if D.IsErr(err) {
		c.done <- D.NewErr(c.stderr.String())
	} else {
		c.done <- nil
	}
}

// Output ...
func (c *Cmd) Output() string {
	return c.str + "\n" + c.stdout.String()
}

// Cancel ...
func (c *Cmd) Cancel() {
	c.cmd.Process.Kill()
}
