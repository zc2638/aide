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
	"reflect"
	"runtime"

	"github.com/99nil/go/stage"
)

type Stage struct {
	name     string
	total    int
	instance *stage.Instance
}

func NewStage(name string) *Stage {
	s := &Stage{name: name}
	s.instance = stage.New(name).
		SetPreFunc(func(sc stage.Context) error {
			sc.WithValue(StepTotalKey, s.total)
			stageName := stage.ContextName(sc)
			logf(Unknown, "%s STAGE %s", stageSymbol, stageName)
			return nil
		}).
		SetSubFunc(sub)
	return s
}

func sub(_ stage.Context) error {
	fmt.Println()
	return nil
}

// TODO asynchronous support in the future
//func (s *Stage) SetAsync(async bool) *Stage {
//	s.instance.SetAsync(async)
//	return s
//}

func (s *Stage) SetRely(names ...string) *Stage {
	s.instance.SetRely(names...)
	return s
}

func (s *Stage) AddSteps(steps ...*Step) *Stage {
	for _, step := range steps {
		if step == nil {
			continue
		}
		s.total++
		step.num = s.total
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
