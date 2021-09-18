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

import "github.com/sirupsen/logrus"

type LogLevel int

const (
	_ LogLevel = iota
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
		FullTimestamp:          true,
		TimestampFormat:        "2006/01/02 15:04:05",
	})
}

func Output(level LogLevel, format string, args ...interface{}) {
	switch level {
	case ErrorLevel:
		logrus.Errorf(format, args...)
	case WarnLevel:
		logrus.Warningf(format, args...)
	case InfoLevel:
		logrus.Infof(format, args...)
	default:
		logrus.Infof(format, args...)
	}
}

func OutputErr(level LogLevel, format string, args ...interface{}) {
	switch level {
	case ErrorLevel:
	case WarnLevel:
	case InfoLevel:
	default:
		level = ErrorLevel
	}
	Output(level, format, args...)
}
