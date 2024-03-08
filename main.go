package main

import (
	"bufio"
	"errors"
	"io"
	"io/fs"
	"log"
	"os"
	"strings"
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

func filterDotEnv(files []fs.DirEntry) []fs.DirEntry {
	dotEnvFiles := make([]fs.DirEntry, 0)
	for _, file := range files {
		if strings.HasPrefix(file.Name(), ".env") {
			dotEnvFiles = append(dotEnvFiles, file)
		}
	}
	return dotEnvFiles
}

func getKeyValue(line string) (string, string) {
	pair := strings.Split(line, "=")
	return pair[0], pair[1]
}

func readDotEnv(files []fs.DirEntry) map[string]string {
	result := make(map[string]string) // key-value pair
	for _, file := range files {
		envFile, err := os.Open(file.Name())
		if err != nil {
			log.Fatalf("can't read file %s, error %s", file.Name(), err)
		}
		defer envFile.Close()
		envReader := bufio.NewReader(envFile)
		for {
			line, err := envReader.ReadString('\n')
			if err != nil {
				if errors.Is(err, io.EOF) {
					key, value := getKeyValue(line)
					result[key] = value
					break
				}
				log.Fatalf("error reading file %s, error %s", file.Name(), err)
			}
			key, value := getKeyValue(line)
			result[key] = value
		}
	}
	return result
}

func main() {
	for key, value := range readDotEnv(filterDotEnv(readCWD())) {
		log.Printf("key: %s, value: %s\n", key, value)
	}
}
