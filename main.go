package main

import (
	"io/fs"
	"log"
	"os"
)

func readCWD() []fs.DirEntry {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("can't read current directory, error %s", err)
	}
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalf("can't read directory %s, error %s", dir, err)
	}
	return files
}

func main() {
	for _, file := range readCWD() {
		log.Println(file.Name())
	}
}
