package main

import (
	"encoding/json"
	"fmt"
	"github.com/zacharyworks/sync/packages/dir"
	"io"
	"log"
	"net/http"
	"os"
)

//const toPath = "/Users/zacharybriggs/sync/to"

func main() {
	// start HTTP server
	h := dirHandler{toPath: os.Args[1]}
	// provide handler for (CREATE / UPDATE) / DELETE
	http.HandleFunc("/update", h.handleRequest)
	log.Fatalln(http.ListenAndServe(":8090", nil))
}

type dirHandler struct {
	toPath string
}

func (d dirHandler) handleRequest(w http.ResponseWriter, r *http.Request) {
	// read body
	defer r.Body.Close()
	file, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var fileChange dir.FileChange
	if err := json.Unmarshal(file, &fileChange); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// very primitive create/delete via os package
	switch fileChange.Modification {
	case dir.Create, dir.Update:
		if err := os.WriteFile(fmt.Sprintf("%s/%s", d.toPath, fileChange.F.Name), fileChange.F.Content, 0600); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	case dir.Delete:
		if err := os.Remove(fmt.Sprintf("%s/%s", d.toPath, fileChange.F.Name)); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
