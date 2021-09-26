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
	"unicode"

	"github.com/sirupsen/logrus"
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

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:            true,
		DisableLevelTruncation: true,
		PadLevelText:           true,
		DisableTimestamp:       true,
		//FullTimestamp:          true,
		//TimestampFormat: "2006/01/02 15:04:05",
	})
}

var DefaultLog LogInterface = &defaultLog{entry: logrus.StandardLogger()}

type defaultLog struct {
	entry *logrus.Logger
}

func (l *defaultLog) Writer() io.Writer {
	return l.entry.Out
}

func (l *defaultLog) Log(level LogLevel, args ...interface{}) {
	switch level {
	case ErrorLevel:
		l.entry.Errorln(args...)
	case WarnLevel:
		l.entry.Warningln(args...)
	case InfoLevel:
		l.entry.Infoln(args...)
	default:
		fmt.Println(args...)
	}
}

func (l *defaultLog) Logf(level LogLevel, format string, args ...interface{}) {
	switch level {
	case ErrorLevel:
		l.entry.Errorf(format, args...)
	case WarnLevel:
		l.entry.Warningf(format, args...)
	case InfoLevel:
		l.entry.Infof(format, args...)
	default:
		fmt.Println(fmt.Sprintf(format, args...))
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
