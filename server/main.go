package main

import (
	"fmt"
	"log"
	"net/http"

	"./executer"
	"./formula"
	"./litedb"
	"./repo"
	D "github.com/NeoJRotary/describe-go"
)

const dataPath = "/formulaci/data/"

var (
	db       *litedb.DB
	repoM    *repo.Manager
	formulaM *formula.Master
	execEX   *executer.Exec
)

func init() {
	var e error
	// runCmd("bash", "init.sh")
	db, e = litedb.Init(dataPath + "sqlite.db")
	if D.IsErr(e) {
		log.Fatalln(D.ErrWithTitle(e, "litedb init :"))
	}
	repoM, e = repo.Init(dataPath+"repo/", db)
	if D.IsErr(e) {
		log.Fatalln(D.ErrWithTitle(e, "repo init :"))
	}
	formulaM, e = formula.Init(".formulaci.yaml", repoM, execEX)
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
