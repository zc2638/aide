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
	"context"
	"errors"
	"fmt"
	"unicode"
)

type Step struct {
	name string
	srf  StepResultFunc
}

func (s *Step) SetResultFunc(srf StepResultFunc) *Step {
	s.srf = srf
	return s
}

func (s *Step) run(sc *StepContext) error {
	prefix := fmt.Sprintf(stepPrefixFormat, s.name)
	if s.srf == nil {
		Output(InfoLevel, "%s Nothing to run.", prefix)
		return nil
	}

	s.srf(sc)

	message := sc.message
	if message != "" {
		rs := []rune(message)
		if unicode.IsLetter(rs[0]) && !unicode.IsUpper(rs[0]) {
			rs[0] = unicode.ToUpper(rs[0])
		}
		message = string(rs)
	}

	if sc.exitCode > 0 {
		OutputErr(sc.level, "%s %s", prefix, message)
		return errors.New(message)
	}
	Output(sc.level, "%s %s", prefix, message)
	return nil
}

type StepContext struct {
	ctx context.Context

	level LogLevel
	// exitCode defines the state when an exception exits.
	exitCode int32
	// message describes execution results
	message string
}

func (c *StepContext) SetLevel(level LogLevel) {
	c.level = level
}

func (c *StepContext) Exit(code int32) {
	c.exitCode = code
}

func (c *StepContext) Write(b []byte) {
	c.WriteString(string(b))
}

func (c *StepContext) WriteString(s string) {
	c.message = s
}

func (c *StepContext) Context() context.Context {
	if c.ctx != nil {
		return c.ctx
	}
	return context.Background()
}

func (c *StepContext) WithContext(ctx context.Context) {
	if ctx == nil {
		panic("nil context")
	}
	c.ctx = ctx
}

type StepResultFunc func(ctx *StepContext)

func (f StepResultFunc) Step(name string) *Step {
	return &Step{name: name, srf: f}
}
