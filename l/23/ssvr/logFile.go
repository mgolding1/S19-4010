package main

// Setup of log file
// This file is MIT licensed.

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pschlump/filelib"
)

// TODO - log rotation every X hours/days etc.

// LogFile sets the output log file to an open file.  This will turn on logging of SQL statments.
func LogFile(fn string) {
	logFileName = fn
	dir := filepath.Dir(logFileName)
	if !filelib.Exists(dir) {
		os.MkdirAll(dir, 0755)
	}
	f, err := filelib.Fopen(logFileName, "a")
	if err != nil {
		fmt.Fprintf(os.Stderr, "log file confiured, but unable to open, file[%s] error[%s]\n", gCfg.LogFileName, err)
		os.Exit(1)
	}
	logFilePtr = f
}
