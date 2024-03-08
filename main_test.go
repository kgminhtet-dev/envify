package main

import (
	"os"
	"testing"
)

func TestEnvVaraibles(t *testing.T) {
	testcases := map[string]string{
		"USER":     "me",
		"HOST":     "lcalhost",
		"DNAME":    "you",
		"PASSWORD": "1234567",
		"PNAME":    "I",
	}
	setEnvVariables(readDotEnv(filterDotEnv(readCWD())))
	for key, value := range testcases {
		t.Log(key)
		result := os.Getenv(key)
		if result != value {
			t.Log("vaule", value, "result", result, "\n")
		}
	}
}
