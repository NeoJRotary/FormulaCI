package main

import (
	"fmt"
	"log"
	"net/http"
	"path"

	"./executer"
	"./formula"
	"./git"
	"./litedb"
	"./repo"
	D "github.com/NeoJRotary/describe-go"
)

const dataPath = "/formulaci/data/"

var (
	db       *litedb.DB
	gitA     *git.Agent
	repoM    *repo.Manager
	formulaM *formula.Master
	execEX   *executer.Exec
)

func init() {
	var err error

	db, err = litedb.Init(path.Join(dataPath, "sqlite.db"))
	if D.IsErr(e) {
		log.Fatalln(D.ErrWithTitle("litedb init :", e))
	}

	gitA, err = git.NewAgent(dataPath)

	repoM, err = repo.Init(path.Join(dataPath, "repo"), db, gitA)
	if D.IsErr(e) {
		log.Fatalln(D.ErrWithTitle("repo init :", e))
	}

	formulaM, err = formula.Init(".formulaci.yaml", repoM, execEX)
	// ci.init()
	fmt.Println("Formula CI Server Init Done")
}

func main() {
	mux := http.NewServeMux()
	startWebServer(mux)
	startFormulaServer(mux)
	fmt.Println("FormulaCI Listen On : 8099")
	http.ListenAndServe(":8099", mux)
}
