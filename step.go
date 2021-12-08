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
	"os/exec"

	"github.com/zc2638/aide/stage"
)

type Step struct {
	instance *stage.Instance
	stage    *Stage

	name string
	num  int
	srf  StepFunc
}

func (s *Step) RelyOn(names ...string) *Step {
	s.instance.RelyOn(names...)
	return s
}

func (s *Step) SetFunc(srf StepFunc) *Step {
	s.srf = srf
	s.instance.SetSubFunc(s.execute)
	return s
}

func (s *Step) execute(sc stage.Context) error {
	stepCtx, ok := sc.Value(StepCtxKey).(*StepContext)
	if !ok {
		stepCtx = &StepContext{
			ctx:    sc,
			logger: s.stage.logger,
		}
	}
	s.run(stepCtx)

	if stepCtx.err != nil {
		level := stepCtx.level

		switch level {
		case ErrorLevel:
		case WarnLevel:
		case InfoLevel:
		default:
			level = ErrorLevel
		}
		if stepCtx.err != stage.ErrStageEnd {
			s.stage.logger.Logf(level, "%s", standardMessage(stepCtx.err.Error()))
		}
	}
	return stepCtx.err
}

func (s *Step) run(sc *StepContext) {
	if s.srf == nil {
		return
	}

	defer func() {
		if v := recover(); v != nil {
			sc, ok := v.(*StepContext)
			if !ok {
				sc.err = fmt.Errorf("unexpected error: %v", v)
			}
			if sc.err == nil {
				sc.err = stage.ErrStageEnd
			}
		}
	}()

	s.srf(sc)
}

type StepFunc func(sc *StepContext)

func (f StepFunc) Step(name string) *Step {
	step := &Step{
		name:     name,
		instance: stage.New(name),
	}
	step.SetFunc(f)
	return step
}

type StepContext struct {
	ctx    stage.Context
	logger LogInterface

	level LogLevel
	// err defines the error when an exception exits.
	err error
}

func (c *StepContext) clear() {
	c.level = Unknown
	c.err = nil
}

func (c *StepContext) Log(args ...interface{}) {
	c.Logl(InfoLevel, args...)
}

func (c *StepContext) Logf(format string, args ...interface{}) {
	c.Logfl(InfoLevel, format, args...)
}

func (c *StepContext) Logl(level LogLevel, args ...interface{}) {
	c.logger.Log(level, args...)
}

func (c *StepContext) Logfl(level LogLevel, format string, args ...interface{}) {
	c.logger.Logf(level, format, args...)
}

// ErrorStr exits all execution and return a error by string.
func (c *StepContext) ErrorStr(s string) {
	c.Error(errors.New(s))
}

// Error exits all execution and return a error.
func (c *StepContext) Error(err error) {
	c.err = err
	c.Exit()
}

// Errorf exits all execution and return a error by format.
func (c *StepContext) Errorf(format string, args ...interface{}) {
	c.err = fmt.Errorf(format, args...)
	c.Exit()
}

// Exit exits all execution.
func (c *StepContext) Exit() {
	panic(c)
}

// Context returns a stage.Context
func (c *StepContext) Context() stage.Context {
	return c.ctx
}

// WithContext returns a shallow copy of r with its context changed
// to ctx. The provided ctx must be non-nil.
func (c *StepContext) WithContext(ctx context.Context) {
	if ctx == nil {
		return
	}
	c.ctx.WithCtx(ctx)
}

// Shell helps execute shell scripts.
func (c *StepContext) Shell(command string) error {
	cmd := exec.CommandContext(c.Context(), "sh", "-c", command)
	cmd.Stdout = c.logger.Writer()
	cmd.Stderr = c.logger.Writer()
	return cmd.Run()
}
