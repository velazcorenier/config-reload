package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/go-fsnotify/fsnotify"
)

//
var watcher *fsnotify.Watcher

// main
func main() {

	// creates a new file watcher
	watcher, _ = fsnotify.NewWatcher()
	defer watcher.Close()

	// starting at the root of the project, walk each file/directory searching for
	// directories
	if err := filepath.Walk("/Users/renier.velazco@ibm.com/Desktop/config-reload/config-reload-premium/test", watchDir); err != nil {
		fmt.Println("ERROR", err)
	}

	done := make(chan bool)

	go func() {
		for {
			select {
			// watch for events
			case event := <-watcher.Events:
				fmt.Printf("File Changed %#v\n", event.Name)
				fmt.Printf("Event Operation %#v\n", event.Op.String())
				readFileContent(event.Name)

			// watch for errors
			case err := <-watcher.Errors:
				fmt.Println("ERROR", err)
			}
		}
	}()

	<-done
}

// watchDir gets run as a walk func, searching for directories to add watchers to
func watchDir(path string, fi os.FileInfo, err error) error {

	// since fsnotify can watch all the files in a directory, watchers only need
	// to be added to each nested directory
	if fi.Mode().IsDir() {
		return watcher.Add(path)
	}

	return nil
}

func readFileContent(filePath string) {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {

		log.Fatal(err)
	}
	fmt.Print("New Content: " + string(file))
}
