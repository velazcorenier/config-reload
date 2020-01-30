package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/go-fsnotify/fsnotify"
)

var watcher *fsnotify.Watcher

func main() {

	// creates a new file watcher
	watcher, _ = fsnotify.NewWatcher()
	defer watcher.Close()

	// Starting at the root of the project
	if err := filepath.Walk("/Users/renier.velazco@ibm.com/Desktop/config-reload/config-reload-premium/test", watchDir); err != nil {
		fmt.Println("ERROR", err)
	}

	done := make(chan bool)

	go func() {
		for {
			select {
			// watch for events
			case event := <-watcher.Events:
				operation := event.Op.String()
				path := event.Name

				eventInformation(path, operation)

				switch operation {
				case "CHMOD":
					//We are going to filter this Operation
				case "REMOVE":
					watcher.Remove(path)
				default:
					readFileContent(path)
				}

			// watch for errors
			case err := <-watcher.Errors:
				fmt.Println("ERROR", err)
			}
		}
	}()

	<-done
}

// Gets run as a walk func, searching for directories to add watchers to...
func watchDir(path string, fi os.FileInfo, err error) error {

	// since fsnotify can watch all the files in a directory, watchers only need to be added to nested directories
	if fi.Mode().IsDir() {
		return watcher.Add(path)
	}

	return nil
}

//Print file data
func readFileContent(filePath string) {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {

		log.Fatal(err)
	}
	if len(file) > 0 {
		fmt.Print("New Content: " + string(file))
	}
}

//Print event summary
func eventInformation(filepath string, event string) {
	fmt.Printf("File Changed %#v\n", filepath)
	fmt.Printf("Event Operation %#v\n", event)
}
