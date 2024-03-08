package main

import (
	"bufio"
	"bytes"
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

func splitToKeyValue(line []byte) (string, string) {
	pair := bytes.Split(line, []byte("="))
	return string(pair[0]), string(pair[1])
}

func getKeyValuePair(file *os.File) ([]string, []string) {
	defer file.Close()
	key, value := "", ""
	keys, values := make([]string, 0), make([]string, 0)
	envReader := bufio.NewReader(file)
	for {
		line, _, err := envReader.ReadLine()
		if err != nil {
			if errors.Is(err, io.EOF) {
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
		keys, values := getKeyValuePair(envFile)
		for i, key := range keys {
			result[key] = values[i]
		}
	}
	return result
}

func setEnvVariables(pairs map[string]string) {
	for key, value := range pairs {
		setEnvVariable(key, value)
	}
}

func setEnvVariable(key, value string) {
	if err := os.Setenv(key, value); err != nil {
		log.Fatalf("error occured in setting environment variable %s", err)
	}
}

func Load() {
	setEnvVariables(readDotEnv(filterDotEnv(readCWD())))
}

// func main() {
// 	for key, value := range readDotEnv(filterDotEnv(readCWD())) {
// 		setEnvVariable(key, value)
// 		log.Printf("os.Getenv(%v) = value(%v)", os.Getenv(key), value)
// 	}
// }
