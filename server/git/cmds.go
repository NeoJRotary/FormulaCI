package git

// CmdConfigGlobalUserEmail git config --global user.email {email}
func CmdConfigGlobalUserEmail(email string) []string {
	return []string{"git", "config", "--global", "user.email", email}
}

// CmdSSHKeygenGenerate ssh-keygen -t rsa -P "" -C {email} -b 4096 -f {path}id_rsa
func CmdSSHKeygenGenerate(email, dataPath string) []string {
	return []string{"ssh-keygen", "-t", "rsa", "-P", "", "-C", email, "-b", "4096", "-f", dataPath + "id_rsa"}
}

// CmdSSHAdd ssh-add {path}
func CmdSSHAdd(path string) []string {
	return []string{"ssh-add", path}
}

// CmdRemoteList git remote
func CmdRemoteList() []string {
	return []string{"git", "remote"}
}

// CmdRemoteAdd git remote add {name} {src}
func CmdRemoteAdd(name, src string) []string {
	return []string{"git", "remote", "add", name, src}
}

// CmdRemoteRemove git remote remove {name}
func CmdRemoteRemove(name string) []string {
	return []string{"git", "remote", "remove", name}
}

// CmdRemoteGetURL git remote get-url {name}
func CmdRemoteGetURL(name string) []string {
	return []string{"git", "remote", "get-url", name}
}

// CmdRemoteSetURL git remote set-url {name} {src}
func CmdRemoteSetURL(name, src string) []string {
	return []string{"git", "remote", "set-url", name, src}
}

// CmdCheckoutBranch git checkout {branch}
func CmdCheckoutBranch(branch string) []string {
	return []string{"git", "checkout", branch}
}

// CmdCheckoutNewBranch git checkout -b {branch}
func CmdCheckoutNewBranch(branch string) []string {
	return []string{"git", "checkout", "-b", branch}
}

// CmdBranchList git branch
func CmdBranchList() []string {
	return []string{"git", "branch"}
}

// CmdResetHardLatest git reset --hard HEAD
func CmdResetHardLatest() []string {
	return []string{"git", "reset", "--hard", "HEAD"}
}

// CmdCleanUntracked git clean -df
func CmdCleanUntracked() []string {
	return []string{"git", "clean", "-df"}
}

// CmdPullRemoteBranch git pull {remote} {branch}
func CmdPullRemoteBranch(remote, branch string) []string {
	return []string{"git", "pull", remote, branch}
}

// CmdRevParseHead git rev-parse HEAD
func CmdRevParseHead() []string {
	return []string{"git", "rev-parse", "HEAD"}
}

// CmdDiffHeadNameOnly git diff --name-only HEAD {sha}
func CmdDiffHeadNameOnly(sha string) []string {
	return []string{"git", "diff", "--name-only", "HEAD", sha}
}
