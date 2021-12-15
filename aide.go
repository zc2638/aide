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
	"fmt"
	"strings"

	"github.com/zc2638/aide/stage"
)

const (
	stageSymbol = "[+] STAGE %s"
	stepSymbol  = "=> %s"
)

type Instance struct {
	instance *stage.Instance
	logger   LogInterface

	stageSymbol string
	stepSymbol  string
	verbose     bool
}

func New(opts ...InstanceOption) *Instance {
	ins := &Instance{
		instance:    stage.New(""),
		stageSymbol: stageSymbol,
		stepSymbol:  stepSymbol,
		verbose:     true,
	}
	for _, opt := range opts {
		opt(ins)
	}
	ins.logger = newLog(ins.verbose)
	return ins
}

func (i *Instance) AddStages(stages ...*Stage) *Instance {
	for _, s := range stages {
		if s == nil {
			continue
		}
		if s.preFunc == nil {
			s.instance.SetPreFunc(i.buildPre(s))
		}
		if s.subFunc == nil {
			s.instance.SetSubFunc(sub)
		}
		s.SetLogger(i.logger)
		s.SetSymbol(i.stepSymbol)
		s.instance.Skip(s.skip)
		s.instance.SkipFunc(s.skipFunc)
		i.instance.Add(s.instance)
	}
	return i
}

func (i *Instance) Run(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}
	return i.instance.Run(ctx)
}

func (i *Instance) buildPre(s *Stage) func(sc stage.Context) error {
	return func(sc stage.Context) error {
		sc.WithValue(StepTotalKey, s.total)
		stageName := stage.ContextName(sc)
		if len(i.stageSymbol) == 0 {
			return nil
		}
		if strings.Count(i.stageSymbol, "%s") > 0 {
			i.logger.Logf(Unknown, i.stageSymbol, stageName)
		} else {
			i.logger.Log(Unknown, i.stageSymbol)
		}
		return nil
	}
}

func sub(_ stage.Context) error {
	fmt.Println()
	return nil
}

type InstanceOption func(i *Instance)

func WithSymbolOption(stage string, step string) InstanceOption {
	return func(i *Instance) {
		i.stageSymbol = stage
		i.stepSymbol = step
	}
}

func WithVerboseOption(verbose bool) InstanceOption {
	return func(i *Instance) {
		i.verbose = verbose
	}
}

func WithLogOption(log LogInterface) InstanceOption {
	return func(i *Instance) {
		i.logger = log
	}
}
