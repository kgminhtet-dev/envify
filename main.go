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

func splitToKeyValue(line string) (string, string) {
	pair := strings.Split(line, "=")
	return pair[0], pair[1]
}

func getKeyValuePair(file *os.File) ([]string, []string) {
	key, value := "", ""
	keys, values := make([]string, 0), make([]string, 0)
	envReader := bufio.NewReader(file)
	for {
		line, err := envReader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				key, value = splitToKeyValue(line)
				keys, values = append(keys, key), append(values, value)
				break
			}
			log.Fatalf("error reading file %s, error %s", file.Name(), err)
		}
		key, value = splitToKeyValue(line)
		keys, values = append(keys, key), append(values, value)
	}
	return keys, values
}

func readDotEnv(files []fs.DirEntry) map[string]string {
	result := make(map[string]string) // key-value pair
	for _, file := range files {
		envFile, err := os.Open(file.Name())
		if err != nil {
			log.Fatalf("can't read file %s, error %s", file.Name(), err)
		}
		defer envFile.Close()
		keys, values := getKeyValuePair(envFile)
		for i, key := range keys {
			result[key] = values[i]
		}
	}
	return result
}

func main() {
	for key, value := range readDotEnv(filterDotEnv(readCWD())) {
		log.Printf("key: %s, value: %s\n", key, value)
	}
}
