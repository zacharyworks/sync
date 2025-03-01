package main

import (
	"fmt"
	"github.com/zacharyworks/sync/packages/dir"
	"github.com/zacharyworks/sync/packages/persistance"
	"log"
	"os"
	"time"
)

func main() {
	store := persistance.New()
	for {
		files, _ := readDirFiles("/Users/zacharybriggs/sync/from")
		changes := store.Update(files)
		fmt.Printf("%+v\n", changes)
		time.Sleep(5 * time.Second)
	}
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
		f = append(f, dir.File{Name: fileInfo.Name(), LastUpdated: fileInfo.ModTime()})
	}
	return f, nil
}
