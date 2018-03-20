package git

// CmdConfigGlobalUserEmail git config --global user.email {email}
func CmdConfigGlobalUserEmail(email string) []string {
	return []string{"git", "config", "--global", "user.email", email}
}

// CmdSSHKeygenGenerate ssh-keygen -t rsa -P "" -C {email} -b 4096 -f {path}id_rsa
func CmdSSHKeygenGenerate(email, dataPath string) []string {
	return []string{"ssh-keygen", "-t", "rsa", "-P", "", "-C", email, "-b", "4096", "-f", dataPath + "id_rsa"}
}
