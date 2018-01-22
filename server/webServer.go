package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"path/filepath"
)

type webServer struct {
	list map[string]webStaticFile
	mux  *http.ServeMux
}

type webStaticFile struct {
	mime string
	b    []byte
}

const webFilesPath = "/formulaci/web"

var webSvr = webServer{
	list: map[string]webStaticFile{},
}

func startWebServer(mux *http.ServeMux) {
	webSvr.mux = mux
	webSvr.readFiles(webFilesPath)
	mux.HandleFunc("/", webSvr.httpHandler)
}

func (wsvr *webServer) readFiles(base string) {
	files, err := ioutil.ReadDir(base)
	if isErr(err) {
		log.Fatal(err)
	}

	for _, file := range files {
		fp := filepath.Join(base, file.Name())
		if file.IsDir() {
			wsvr.readFiles(fp)
			return
		}
		sf := webStaticFile{}
		sf.mime = mime.TypeByExtension(filepath.Ext(file.Name()))
		b, err := ioutil.ReadFile(fp)
		if isErr(err) {
			log.Fatal(err)
		}
		var buf bytes.Buffer
		zipper := gzip.NewWriter(&buf)
		_, err = zipper.Write(b)
		if isErr(err) {
			log.Fatal(err)
		}
		zipper.Close()
		sf.b = buf.Bytes()
		key, err := filepath.Rel(base, fp)
		if isErr(err) {
			log.Fatal(err)
		}
		wsvr.list[key] = sf
		fmt.Println(key)
		// wsvr.mux.HandleFunc(key, wsvr.httpHandler)
	}
}

func (wsvr *webServer) httpHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	path := r.URL.Path

	if path[0] == '/' {
		path = path[1:]
	}
	if _, ok := wsvr.list[path]; !ok {
		// w.WriteHeader(http.StatusNotFound)
		// return
		path = "index.html"
	}

	sf := wsvr.list[path]
	w.Header().Set("Content-Type", sf.mime)
	w.Header().Set("Content-Encoding", "gzip")
	w.Write(sf.b)
}

// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 	if pusher, ok := w.(http.Pusher); ok {
// 			// Push is supported.
// 			if err := pusher.Push("/app.js", nil); err != nil {
// 					log.Printf("Failed to push: %v", err)
// 			}
// 	}
// 	// ...
// })
