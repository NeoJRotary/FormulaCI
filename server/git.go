package main

import (
	"errors"
	"io/ioutil"
	"strings"
)

type gitAPI struct {
	email        string
	webhookToken string
}

var git = gitAPI{}

func (git *gitAPI) getInfo(data interface{}, resFunc wsResFunc) {
	var pubkey = ""
	buf, err := ioutil.ReadFile("/formulaci/data/id_rsa.pub")
	if !isErr(err) {
		pubkey = string(buf)
	}
	obj := map[string]interface{}{
		"email":        git.email,
		"webhookToken": git.webhookToken,
		"pubkey":       pubkey,
	}
	resFunc(200, obj, "")
}

func (git *gitAPI) setEmail(data interface{}, resFunc wsResFunc) {
	email := data.(string)
	_, err := cmdEX.run("git", "config", "--global", "user.email", email)
	if isErr(err) {
		resFunc(500, nil, err.Error())
	}
	updateConfig("git/email", email)
	git.email = email
	resFunc(200, nil, "")
}

func (git *gitAPI) setWebhookToken(data interface{}, resFunc wsResFunc) {
	token := data.(string)
	updateConfig("git/webhookToken", token)
	git.webhookToken = token
	resFunc(200, nil, "")
}

func (git *gitAPI) generateSSH(data interface{}, resFunc wsResFunc) {
	cmdEX.run("rm", "/formulaci/data/id_rsa")
	cmdEX.run("rm", "/formulaci/data/id_rsa.pub")
	_, err := cmdEX.run("ssh-keygen", "-t", "rsa", "-P", "", "-C", git.email, "-b", "4096", "-f", "/formulaci/data/id_rsa")
	if isErr(err) {
		resFunc(500, nil, err.Error())
		return
	}
	pubkey, err := cmdEX.run("cat", "/formulaci/data/id_rsa.pub")
	if isErr(err) {
		resFunc(500, nil, err.Error())
		return
	}
	resFunc(200, pubkey, "")
}

func (git *gitAPI) pullRepo(name string, branch string, src string) (err error) {
	dir := repoPath + name + "/" + branch
	exist := fileExist(dir)

	defer func() {
		if !exist && isErr(err) {
			cmdEX.run("rm", "-rf", dir)
		}
	}()

	// _, err = cmdEX.run(dir, "ssh-add", "/formulaci/data/id_rsa")

	if !exist {
		if src == "" {
			return errors.New("empty src")
		}
		cmdEX.run("mkdir", "-p", dir)
		cmdEX.runDir(dir, "git", "init")
		_, err = cmdEX.runDir(dir, "git", "remote", "add", "origin", src)
		if isErr(err) {
			return err
		}
		if branch != "master" {
			_, err = cmdEX.runDir(dir, "git", "checkout", "-b", branch)
			if isErr(err) {
				return err
			}
		}
	} else {
		_, err = cmdEX.runDir(dir, "git", "reset", "--hard", "HEAD")
		if isErr(err) {
			return err
		}
		_, err = cmdEX.runDir(dir, "git", "clean", "-df")
		if isErr(err) {
			return err
		}
	}

	_, err = cmdEX.runDir(dir, "git", "pull", "origin", branch)
	if isErr(err) {
		return err
	}

	return nil
}

func (git *gitAPI) headSHA(name string, branch string) (string, error) {
	dir := repoPath + name + "/" + branch
	o, err := cmdEX.runDir(dir, "git", "rev-parse", "HEAD")
	if isErr(err) {
		return "", err
	}

	return strings.Replace(o, "\n", "", -1), nil
}

func (git *gitAPI) diffWithHEAD(name string, branch string, sha string) ([]string, error) {
	dir := repoPath + name + "/" + branch
	o, err := cmdEX.runDir(dir, "git", "diff", "--name-only", "HEAD", sha)
	if isErr(err) {
		return nil, err
	}
	return strings.Split(o, "\n"), nil
}
