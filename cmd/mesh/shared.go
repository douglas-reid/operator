// Copyright 2017 Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package mesh contains types and functions that are used across the full
// set of mixer commands.
package mesh

import (
	"fmt"
	"io"
	"os"

	"istio.io/pkg/log"
)

func initLogsOrExit(args *rootArgs) {
	if err := configLogs(args.logToStdErr); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Could not configure logs: %s", err)
		os.Exit(1)
	}
}

func configLogs(logToStdErr bool) error {
	opt := log.DefaultOptions()
	if !logToStdErr {
		opt.SetOutputLevel(log.OverrideScopeName, log.NoneLevel)
	}
	return log.Configure(opt)
}

//Logger is the struct used for mesh command
type Logger struct {
	logToStdErr bool
	stdOut      io.Writer
	stdErr      io.Writer
}

// NewLogger creates a new logger and returns a pointer to it.
// stdOut and stdErr can be used to capture output for testing.
func NewLogger(logToStdErr bool, stdOut, stdErr io.Writer) *Logger {
	return &Logger{
		logToStdErr: logToStdErr,
		stdOut:      stdOut,
		stdErr:      stdErr,
	}
}

// TODO: this really doesn't belong here. Figure out if it's generally needed and possibly move to istio.io/pkg/log.
func (l *Logger) logAndPrint(v ...interface{}) {
	if len(v) == 0 {
		return
	}
	s := fmt.Sprint(v...)
	if !l.logToStdErr {
		l.print(s + "\n")
	} else {
		log.Infof(s)
	}
}

func (l *Logger) logAndError(v ...interface{}) {
	if len(v) == 0 {
		return
	}
	s := fmt.Sprint(v...)
	if !l.logToStdErr {
		l.printErr(s + "\n")
	} else {
		log.Infof(s)
	}
}

func (l *Logger) logAndFatal(a ...interface{}) {
	l.logAndError(a...)
	os.Exit(-1)
}

func (l *Logger) logAndPrintf(format string, a ...interface{}) {
	s := fmt.Sprintf(format, a...)
	if !l.logToStdErr {
		l.print(s + "\n")
	} else {
		log.Infof(s)
	}
}

func (l *Logger) logAndErrorf(format string, a ...interface{}) {
	s := fmt.Sprintf(format, a...)
	if !l.logToStdErr {
		l.printErr(s + "\n")
	} else {
		log.Infof(s)
	}
}

func (l *Logger) logAndFatalf(format string, a ...interface{}) {
	l.logAndErrorf(format, a...)
	os.Exit(-1)
}

func (l *Logger) print(s string) {
	_, _ = l.stdOut.Write([]byte(s))
}

func (l *Logger) printErr(s string) {
	_, _ = l.stdErr.Write([]byte(s))
}

func refreshGoldenFiles() bool {
	return os.Getenv("REFRESH_GOLDEN") == "true"
}
