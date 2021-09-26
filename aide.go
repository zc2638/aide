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

	"github.com/zc2638/aide/stage"
)

const (
	stageSymbol = "[+]"
	stepSymbol  = "=>"
)

type Instance struct {
	// TODO stage pre logs define
	// TODO step pre logs define
	instance *stage.Instance
}

func New() *Instance {
	return &Instance{
		instance: stage.New(""),
	}
}

func (i *Instance) AddStages(stages ...*Stage) *Instance {
	for _, s := range stages {
		if s == nil {
			continue
		}
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
