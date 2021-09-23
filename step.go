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
	"os"
	"os/exec"

	"github.com/99nil/go/stage"
)

type Step struct {
	name     string
	num      int
	srf      StepFunc
	instance *stage.Instance
}

func (s *Step) SetRely(names ...string) *Step {
	s.instance.SetRely(names...)
	return s
}

func (s *Step) SetFunc(srf StepFunc) *Step {
	s.srf = srf
	s.instance.SetPreFunc(func(sc stage.Context) error {
		total, ok := sc.Value(StepTotalKey).(int)
		if !ok {
			total = s.num
		}
		name := stage.ContextName(sc)
		logf(Unknown, "%s [%d/%d] %s", stepSymbol, s.num, total, name)
		return nil
	})
	s.instance.SetSubFunc(func(sc stage.Context) error {
		stepCtx, ok := sc.Value(StepCtxKey).(*StepContext)
		if !ok {
			stepCtx = &StepContext{ctx: sc}
		}
		return s.run(stepCtx)
	})
	return s
}

func (s *Step) run(sc *StepContext) error {
	//prefix := fmt.Sprintf(stepPrefixFormat, s.name)
	if s.srf == nil {
		logf(InfoLevel, "%s", "Nothing to run.")
		return nil
	}

	defer func() {
		if v := recover(); v != nil {
			sc, ok := v.(*StepContext)
			if !ok {
				return
			}
			if sc.exitCode > 0 {
				level := sc.level
				switch level {
				case ErrorLevel:
				case WarnLevel:
				case InfoLevel:
				default:
					level = ErrorLevel
				}
				logf(level, "%s", standardMessage(sc.message))
				// Forced exit according to exit code
				os.Exit(sc.exitCode)
			}
		}
		if sc.level == Unknown {
			sc.level = InfoLevel
		}
		if sc.message != "" {
			log(sc.level, standardMessage(sc.message))
		}
	}()

	s.srf(sc)
	return nil
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
	ctx stage.Context

	level LogLevel
	// exitCode defines the state when an exception exits.
	exitCode int
	// message describes execution results
	message string
}

func (c *StepContext) clear() {
	c.level = Unknown
	c.message = ""
	c.exitCode = 0
}

func (c *StepContext) Log(args ...interface{}) {
	c.Logl(InfoLevel, args...)
}

func (c *StepContext) Logf(format string, args ...interface{}) {
	c.Logfl(InfoLevel, format, args...)
}

func (c *StepContext) Logl(level LogLevel, args ...interface{}) {
	log(level, args...)
}

func (c *StepContext) Logfl(level LogLevel, format string, args ...interface{}) {
	logf(level, format, args...)
}

func (c *StepContext) Message(s string) *StepContext {
	c.message = s
	return c
}

func (c *StepContext) Return(s string) {
	c.Message(s).Exit(0)
}

func (c *StepContext) Error(err error) {
	c.Break(err.Error())
}

func (c *StepContext) Break(s string) {
	c.Message(s).Exit(1)
}

func (c *StepContext) Exit(code int) {
	c.exitCode = code
	panic(c)
}

func (c *StepContext) Context() stage.Context {
	return c.ctx
}

func (c *StepContext) WithContext(ctx context.Context) {
	if ctx == nil {
		return
	}
	c.ctx.WithCtx(ctx)
}

func (c *StepContext) Shell(command string) error {
	cmd := exec.Command("sh", "-c", command)
	return cmd.Run()
}
