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
	"reflect"
	"runtime"
	"strings"

	"github.com/zc2638/aide/stage"
)

type Stage struct {
	instance *stage.Instance
	logger   LogInterface

	name     string
	symbol   string
	total    int
	skip     bool
	skipFunc func() bool
}

func NewStage(name string) *Stage {
	s := &Stage{name: name}
	s.instance = stage.New(name)
	return s
}

func (s *Stage) SetSymbol(symbol string) {
	s.symbol = symbol
}

func (s *Stage) SetLogger(logger LogInterface) {
	s.logger = logger
}

func (s *Stage) RelyOn(names ...string) *Stage {
	s.instance.RelyOn(names...)
	return s
}

func (s *Stage) AddStepFunc(name string, sf StepFunc) *Stage {
	step := sf.Step(name)
	s.AddSteps(step)
	return s
}

func (s *Stage) AddSteps(steps ...*Step) *Stage {
	for _, step := range steps {
		if step == nil {
			continue
		}
		s.total++
		step.num = s.total
		step.stage = s
		step.instance.SetPreFunc(func(sc stage.Context) error {
			if len(s.symbol) == 0 {
				return nil
			}
			if strings.Count(s.symbol, "%s") > 0 {
				name := stage.ContextName(sc)
				s.logger.Logf(Unknown, s.symbol, name)
			} else {
				s.logger.Logf(Unknown, s.symbol)
			}
			return nil
		})
		s.instance.Add(step.instance)
	}
	return s
}

func (s *Stage) AddStepFuncs(sfs ...StepFunc) *Stage {
	for _, sf := range sfs {
		fc := runtime.FuncForPC(reflect.ValueOf(sf).Pointer())
		s.AddSteps(sf.Step(fc.Name()))
	}
	return s
}

func (s *Stage) Skip(is bool) *Stage {
	s.skip = is
	return s
}

func (s *Stage) SkipFunc(f func() bool) *Stage {
	s.skipFunc = f
	return s
}
