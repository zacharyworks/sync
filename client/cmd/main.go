package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/zacharyworks/sync/packages/dir"
	"github.com/zacharyworks/sync/packages/persistance"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

//const fromPath = "/Users/zacharybriggs/sync/from"

func main() {
	fromPath := os.Args[1]
	store := persistance.New()
	for {
		files, err := readDirFiles(fromPath)
		if err != nil {
			log.Println(err)
		}
		changes := store.Update(files)
		for _, c := range changes {
			fmt.Printf("%+v\n", c)
			if err := postFileChange(c); err != nil {
				log.Println(err)
			}
		}
		time.Sleep(5 * time.Second)
	}
}

// https://stackoverflow.com/questions/52406709/simple-http-post-in-go-with-file-upload
func postFileChange(fc dir.FileChange) error {
	body, err := json.Marshal(fc)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", "http://localhost:8090/update", bytes.NewBuffer(body))
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

// readDir takes the absolute path to a dir and constructs a dir entry
// of only the dir in said dir.
func readDirFiles(absolutePath string) (files []dir.File, err error) {
	entries, err := os.ReadDir(absolutePath)
	if err != nil {
		log.Fatal(err)
	}

	// loop through entries & construct current level
	var f []dir.File
	for _, e := range entries {
		if e.IsDir() {
			// in the future we could elaborate this to also
			// recurse into subdirectories
			continue
		}
		fileInfo, err := e.Info()
		if err != nil {
			return f, err
		}
		file, err := os.Open(absolutePath + "/" + fileInfo.Name())
		if err != nil {
			return nil, err
		}
		fileContents, err := io.ReadAll(file)
		if err != nil {
			return nil, err
		}
		f = append(f, dir.File{Name: fileInfo.Name(), Path: absolutePath, LastUpdated: fileInfo.ModTime(), Content: fileContents})
	}
	return f, nil
}
