package envify

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"io/fs"
	"log"
	"os"
	"regexp"
	"strings"
)

type Text interface {
	[]byte | string
}

const (
	COMMENTPREFIX = "#"
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

func validKey(key []byte) bool {
	reg := regexp.MustCompile("^[a-zA-Z_]+[a-zA-Z0-9_]*$")
	return reg.Match(key)
}

func splitToKeyValue(line []byte) (string, string) {
	pair := bytes.Split(line, []byte("="))
	key, value := pair[0], pair[1]
	if !validKey(key) {
		log.Fatalf("%s is not valid key for environment variable. Variable names must consist solely of letters, digits, and the underscore ( _ ) and must not begin with a digit.", key)
	}
	return string(bytes.Trim(key, " ")), string(bytes.Trim(removeCommant(value), " "))
}

func removeCommant(line []byte) []byte {
	locCommentPrefix := bytes.Index(line, []byte(COMMENTPREFIX))
	if locCommentPrefix == -1 {
		return line
	}
	return line[:locCommentPrefix]
}

func isCommant(line []byte) bool {
	return bytes.HasPrefix(line, []byte(COMMENTPREFIX))
}

func isEmpty[T Text](line T) bool {
	return len(line) == 0
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
		if !isEmpty(line) && !isCommant(line) {
			key, value = splitToKeyValue(line)
			if !isEmpty(key) && !isEmpty(value) {
				keys, values = append(keys, key), append(values, value)
			}
		}
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
