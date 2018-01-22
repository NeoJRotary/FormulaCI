package main

import (
	"fmt"
	"net/http"
	"sync"
)

const dataPath = "/formulaci/data/"

var mutex = &sync.Mutex{}

func init() {
	// runCmd("bash", "init.sh")
	sqlite.connect()
	repo.init()
	ci.init()
	fmt.Println("Formula CI Server Init Done")
}

func main() {
	mux := http.NewServeMux()
	startWebServer(mux)
	startFormulaServer(mux)
	fmt.Println("FormulaCI Listen On : 8099")
	http.ListenAndServe(":8099", mux)
}
