// Package aide

// Copyright © 2021 zc2638 <zc2638@qq.com>.
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

package aide

import (
	"fmt"
	"io"
	"log"
	"os"
	"unicode"
)

type LogInterface interface {
	Log(level LogLevel, args ...interface{})
	Logf(level LogLevel, format string, args ...interface{})
	Writer() io.Writer
}

type LogLevel int

const (
	Unknown LogLevel = iota
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel
)

var _ LogInterface = (*defaultLog)(nil)

type defaultLog struct {
	entry *log.Logger
}

func newLog(verbose bool) LogInterface {
	entry := log.New(os.Stderr, "", 0)
	if !verbose {
		entry.SetOutput(&emptyWriter{})
	}
	return &defaultLog{entry: entry}
}

func (l *defaultLog) Writer() io.Writer {
	return l.entry.Writer()
}

func (l *defaultLog) Log(level LogLevel, args ...interface{}) {
	switch level {
	case ErrorLevel:
		v := append([]interface{}{"ERROR"}, args...)
		l.entry.Println(v...)
	case WarnLevel:
		v := append([]interface{}{"WARN "}, args...)
		l.entry.Println(v...)
	case InfoLevel:
		v := append([]interface{}{"INFO "}, args...)
		l.entry.Println(v...)
	default:
		l.entry.Println(args...)
	}
}

func (l *defaultLog) Logf(level LogLevel, format string, args ...interface{}) {
	switch level {
	case ErrorLevel:
		l.entry.Println("ERROR", fmt.Sprintf(format, args...))
	case WarnLevel:
		l.entry.Println("WARN ", fmt.Sprintf(format, args...))
	case InfoLevel:
		l.entry.Println("INFO ", fmt.Sprintf(format, args...))
	default:
		l.entry.Printf(format, args...)
	}
}

func standardMessage(s string) string {
	if s != "" {
		rs := []rune(s)
		if unicode.IsLetter(rs[0]) && !unicode.IsUpper(rs[0]) {
			rs[0] = unicode.ToUpper(rs[0])
		}
		s = string(rs)
	}
	return s
}

type emptyWriter struct{}

func (w *emptyWriter) Write(b []byte) (int, error) {
	return len(b), nil
}
