package pocketlog_test

import (
	"learngo-pockets/gordle/pocketlog"
	"os"
	"strings"
	"testing"
)

func ExampleLogger_Debugf() {
	debugLogger := pocketlog.New(pocketlog.LevelDebug)
	debugLogger.Debugf("Hello, %s", "world")
	// Output: [Debug] Hello, world
}

func ExampleLogger_Infof() {
	infoLogger := pocketlog.New(pocketlog.LevelInfo)
	infoLogger.Infof("Hello, %s", "world")
	// Output: [Info] Hello, world
}

func ExampleLogger_Errorf() {
	errorLogger := pocketlog.New(pocketlog.LevelError, pocketlog.WithOutput(os.Stderr))
	errorLogger.Errorf("Hello, %s", "world")
	// Output:
}

const (
	debugMessage = "Why write I still all one, ever the same,"
	infoMessage  = "And keep invention in a noted weed,"
	errorMessage = "That every word doth almost tell my name,"
)

func TestLogger_DebugfInfofErrorf(t *testing.T) {
	tt := map[string]struct {
		level    pocketlog.Level
		expected string
	}{
		"debug": {
			pocketlog.LevelDebug,
			"[Debug] " + debugMessage + "\n" + "[Info] " + infoMessage + "\n" + "[Error] " + errorMessage + "\n",
		},
		"info":  {pocketlog.LevelInfo, "[Info] " + infoMessage + "\n" + "[Error] " + errorMessage + "\n"},
		"error": {pocketlog.LevelError, "[Error] " + errorMessage + "\n"},
	}
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			output := &strings.Builder{}
			//output := &bytes.Buffer{}
			lgr := pocketlog.New(tc.level, pocketlog.WithOutput(output))
			lgr.Debugf(debugMessage)
			lgr.Infof(infoMessage)
			lgr.Errorf(errorMessage)
			if output.String() != tc.expected {
				t.Errorf("different output, got: %s, want: %s", output.String(), tc.expected)
			}
		})
	}
}
