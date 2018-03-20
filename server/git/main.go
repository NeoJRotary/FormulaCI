package git

import (
	"io/ioutil"
	"strings"

	"../executer"
	D "github.com/NeoJRotary/describe-go"
)

// Agent ...
type Agent struct {
	execEX   *executer.Exec
	DataPath string
	Email    string
}

// NewAgent ...
func NewAgent(DataPath string) *Agent {
	return &Agent{DataPath: DataPath}
}

// GetPubKey ...
func (agn *Agent) GetPubKey() string {
	buf, err := ioutil.ReadFile(agn.DataPath + "id_rsa.pub")
	if D.IsErr(err) {
		return ""
	}
	return string(buf)
}

// SetEmail ...
func (agn *Agent) SetEmail(email string) error {
	_, err := agn.execEX.Run(CmdConfigGlobalUserEmail(email)...)
	if D.IsErr(err) {
		return err
	}
	agn.Email = email
	return nil
}

// GenerateSSH ...
func (agn *Agent) GenerateSSH() error {
	if agn.Email == "" {
		return D.NewErr("Agent Email is empty")
	}
	if agn.DataPath == "" {
		return D.NewErr("Agent DataPath is empty")
	}

	execEX := agn.execEX
	execEX.Run("rm", agn.DataPath+"id_rsa")
	execEX.Run("rm", agn.DataPath+"id_rsa.pub")
	_, err := cmdEX.run(CmdSSHKeygenGenerate(agn.Email, agn.DataPath)...)
	if D.IsErr(err) {
		resFunc(500, nil, err.Error())
		return
	}
	pubkey, err := cmdEX.run("cat", "/formulaci/data/id_rsa.pub")
	if D.IsErr(err) {
		resFunc(500, nil, err.Error())
		return
	}
	resFunc(200, pubkey, "")
}

func (agn *Agent) pullRepo(name string, branch string, src string) (err error) {
	dir := repoPath + name + "/" + branch
	exist := fileExist(dir)

	defer func() {
		if !exist && D.IsErr(err) {
			cmdEX.run("rm", "-rf", dir)
		}
	}()

	// _, err = cmdEX.run(dir, "ssh-add", "/formulaci/data/id_rsa")

	if !exist {
		if src == "" {
			return D.NewErr("empty src")
		}
		cmdEX.run("mkdir", "-p", dir)
		cmdEX.runDir(dir, "git", "init")
		_, err = cmdEX.runDir(dir, "git", "remote", "add", "origin", src)
		if D.IsErr(err) {
			return err
		}
		if branch != "master" {
			_, err = cmdEX.runDir(dir, "git", "checkout", "-b", branch)
			if D.IsErr(err) {
				return err
			}
		}
	} else {
		_, err = cmdEX.runDir(dir, "git", "reset", "--hard", "HEAD")
		if D.IsErr(err) {
			return err
		}
		_, err = cmdEX.runDir(dir, "git", "clean", "-df")
		if D.IsErr(err) {
			return err
		}
	}

	_, err = cmdEX.runDir(dir, "git", "pull", "origin", branch)
	if D.IsErr(err) {
		return err
	}

	return nil
}

func (agn *Agent) headSHA(name string, branch string) (string, error) {
	dir := repoPath + name + "/" + branch
	o, err := cmdEX.runDir(dir, "git", "rev-parse", "HEAD")
	if D.IsErr(err) {
		return "", err
	}

	return strings.Replace(o, "\n", "", -1), nil
}

func (agn *Agent) diffWithHEAD(name string, branch string, sha string) ([]string, error) {
	dir := repoPath + name + "/" + branch
	o, err := cmdEX.runDir(dir, "git", "diff", "--name-only", "HEAD", sha)
	if D.IsErr(err) {
		return nil, err
	}
	return strings.Split(o, "\n"), nil
}
