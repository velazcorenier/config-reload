package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/go-fsnotify/fsnotify"
)

// main
func main() {

	filePath := "config.txt"
	fmt.Print("Old content: ")
	readFileContent(filePath)

	// creates a new file watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("ERROR", err)
	}
	defer watcher.Close()

	//
	done := make(chan bool)

	//
	go func() {
		for {
			select {
			// watch for events
			case event := <-watcher.Events:
				fmt.Printf("File Changed %#v\n", event.Name)
				fmt.Printf("Event Operation %#v\n", event.Op.String())
				readFileContent(filePath)

				// watch for errors
			case err := <-watcher.Errors:
				fmt.Println("ERROR", err)
			}
		}
	}()

	// out of the box fsnotify can watch a single file, or a single directory
	if err := watcher.Add(filePath); err != nil {
		fmt.Println("ERROR", err)
	}

	<-done
}

func readFileContent(filePath string) {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {

		log.Fatal(err)
	}
	fmt.Print(string(file))
}
