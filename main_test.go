package envify

import (
	"os"
	"testing"
)

func TestEnvVaraibles(t *testing.T) {
	testcases := map[string]string{
		"USER":     "me",
		"HOST":     "localhost",
		"DNAME":    "you",
		"PASSWORD": "1234567",
		"PNAME":    "I",
		"PHOST":    "localhost",
		"EMPYT":    "",
	}
	setEnvVariables(readDotEnv(filterDotEnv(readCWD())))
	for key, value := range testcases {
		t.Log(key)
		result := os.Getenv(key)
		if result != value {
			t.Logf("vaule %q, result %q\n", value, result)
		}
	}
}
