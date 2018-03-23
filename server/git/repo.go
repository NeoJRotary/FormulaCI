package git

import (
	"strings"

	"../executer"
	D "github.com/NeoJRotary/describe-go"
)

// Repository ...
type Repository struct {
	run           executer.RunFunc
	SSHKey        string
	Dir           string
	TaragetRemote string
	Remotes       []Remote
	CurrentBranch string
	Branches      []string
}

// Remote ...
type Remote struct {
	Name string
	Src  string
}

// NewRepository ...
func (agn *Agent) NewRepository(dir, remote, src, branch string) (repo *Repository, err error) {
	defer D.RecoverErr(nil)

	repo = &Repository{
		run:    agn.execEX.RunDir(dir),
		SSHKey: agn.SSHKey,
		Dir:    dir,
	}

	err = repo.GetRemoteList()
	D.CheckErr(err)

	err = repo.GetBranchList()
	D.CheckErr(err)

	repo.Init(remote, src, branch)

	return repo, nil
}

// GetRemoteList ...
func (repo *Repository) GetRemoteList() (err error) {
	defer D.RecoverErr(nil)

	out, err := repo.run(CmdRemoteList()...)
	D.CheckErr(err)

	for _, name := range strings.Split(out, "\n") {
		if name == "" {
			continue
		}
		out, err = repo.run(CmdRemoteGetURL(name)...)
		D.CheckErr(err)
		repo.Remotes = append(repo.Remotes, Remote{
			Name: name,
			Src:  out,
		})
	}
	return err
}

// HasRemote ...
func (repo *Repository) HasRemote(name, src string) bool {
	for _, r := range repo.Remotes {
		if r.Name == name && r.Src == src {
			return true
		}
	}
	return false
}

// HasRemoteName ...
func (repo *Repository) HasRemoteName(name string) bool {
	for _, r := range repo.Remotes {
		if r.Name == name {
			return true
		}
	}
	return false
}

// GetBranchList ...
func (repo *Repository) GetBranchList() (err error) {
	defer D.RecoverErr(nil)

	out, err := repo.run(CmdBranchList()...)
	D.CheckErr(err)

	ls := D.String(out).Split("\n")
	current := ls.FindHasPrefix("*")
	if D.Found(current) {
		repo.CurrentBranch = current.Get()
		current.Trim("*").TrimSpace().SetInto(ls)
	}
	repo.Branches = ls.Get()

	return nil
}

// HasBranch ...
func (repo *Repository) HasBranch(branch string) bool {
	return D.StringSlice(repo.Branches).Has(branch)
}

// Init ...
func (repo *Repository) Init(remote, src, branch string) (err error) {
	defer D.RecoverErr(nil)

	if !repo.HasRemote(remote, src) {
		if repo.HasRemoteName(remote) {
			_, err = repo.run(CmdRemoteAdd(remote, src)...)
		} else {
			_, err = repo.run(CmdRemoteSetURL(remote, src)...)
		}
		D.CheckErr(err)
	}
	repo.TaragetRemote = remote

	if repo.CurrentBranch != branch {
		if repo.HasBranch(branch) {
			_, err = repo.run(CmdCheckoutBranch(branch)...)
		} else {
			_, err = repo.run(CmdCheckoutNewBranch(branch)...)
		}
	}

	_, err = repo.run(CmdSSHAdd(repo.SSHKey)...)
	D.CheckErr(err)
	err = repo.ResetPull()
	D.CheckErr(err)

	return nil
}

// ResetPull ...
func (repo *Repository) ResetPull() (err error) {
	_, err = repo.run(CmdPullRemoteBranch(repo.TaragetRemote, repo.CurrentBranch)...)
	if D.IsErr(err) {
		// try reset / clean up and pull again
		_, err = repo.run(CmdResetHardLatest()...)
		_, err = repo.run(CmdCleanUntracked()...)
		_, err = repo.run(CmdPullRemoteBranch(repo.TaragetRemote, repo.CurrentBranch)...)
	}
	return err
}

// HeadSHA ...
func (repo *Repository) HeadSHA(name string, branch string) (string, error) {
	o, err := repo.run(CmdRevParseHead()...)
	if D.IsErr(err) {
		return "", err
	}
	return D.String(o).Trim("\n").Get(), nil
}

// CompareWithHead ...
func (repo *Repository) CompareWithHead(sha string) ([]string, error) {
	o, err := repo.run(CmdDiffHeadNameOnly(sha)...)
	if D.IsErr(err) {
		return nil, err
	}
	return D.String(o).Split("\n").Get(), nil
}

// // PullNewBranch ...
// func (agn *Agent) PullNewBranch(dir, branch, src string) error {
// 	RunDir := agn.execEX.RunDir
// 	_, err := RunDir(dir, CmdCheckoutNewBranch(name, src)...)
// 	_, err := RunDir(dir, CmdPullRemoteBranch(name, src)...)
// 	return err
// }

// // ResetPull ...
// func (agn *Agent) ResetPull(dir, remote, branch string) (err error) {
// 	dir := repoPath + name + "/" + branch
// 	exist := fileExist(dir)

// 	defer func() {
// 		if !exist && D.IsErr(err) {
// 			cmdEX.run("rm", "-rf", dir)
// 		}
// 	}()

// 	// _, err = cmdEX.run(dir, "ssh-add", "/formulaci/data/id_rsa")

// 	if !exist {
// 		if src == "" {
// 			return D.NewErr("empty src")
// 		}
// 		cmdEX.run("mkdir", "-p", dir)
// 		cmdEX.runDir(dir, "git", "init")
// 		_, err = cmdEX.runDir(dir, "git", "remote", "add", "origin", src)
// 		if D.IsErr(err) {
// 			return err
// 		}
// 		if branch != "master" {
// 			_, err = cmdEX.runDir(dir, "git", "checkout", "-b", branch)
// 			if D.IsErr(err) {
// 				return err
// 			}
// 		}
// 	} else {
// 		_, err = cmdEX.runDir(dir, "git", "reset", "--hard", "HEAD")
// 		if D.IsErr(err) {
// 			return err
// 		}
// 		_, err = cmdEX.runDir(dir, "git", "clean", "-df")
// 		if D.IsErr(err) {
// 			return err
// 		}
// 	}

// 	_, err = cmdEX.runDir(dir, "git", "pull", "origin", branch)
// 	if D.IsErr(err) {
// 		return err
// 	}

// 	return nil
// }
