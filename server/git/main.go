package git

import (
	"io/ioutil"

	"../executer"
	D "github.com/NeoJRotary/describe-go"
)

// Agent ...
type Agent struct {
	execEX   *executer.Exec
	DataPath string
	Email    string
	SSHKey   string
}

// NewAgent ...
func NewAgent(DataPath string, execEX *executer.Exec) *Agent {
	agn := Agent{DataPath: DataPath}
	if execEX == nil {
		agn.execEX = executer.NewExec(false)
	} else {
		agn.execEX = execEX
	}
	return &agn
}

// GetPubKey ...
func (agn *Agent) GetPubKey() string {
	buf, err := ioutil.ReadFile(agn.SSHKey + ".pub")
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
func (agn *Agent) GenerateSSH() (err error) {
	defer D.RecoverErr(nil)

	if agn.Email == "" {
		return D.NewErr("Agent Email is empty")
	}
	if agn.DataPath == "" {
		return D.NewErr("Agent DataPath is empty")
	}

	execEX := agn.execEX
	execEX.Run("rm", agn.DataPath+"id_rsa")
	execEX.Run("rm", agn.DataPath+"id_rsa.pub")
	_, err = execEX.Run(CmdSSHKeygenGenerate(agn.Email, agn.DataPath)...)
	D.CheckErr(err)
	agn.SSHKey = agn.DataPath + "id_rsa"
	return err
}
