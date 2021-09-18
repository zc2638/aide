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
)

const (
	stagePrefixFormat = "[+] STAGE %s"
	stepPrefixFormat  = "=> [STEP](%s)"
)

type Instance struct {
	stages []Stage
}

func New() *Instance {
	return &Instance{}
}

func (i *Instance) AddStages(stages ...*Stage) {
	for _, stage := range stages {
		if stage == nil {
			continue
		}
		i.stages = append(i.stages, *stage)
	}
}

func (i *Instance) Run(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}
	for _, stage := range i.stages {
		if err := stage.run(ctx); err != nil {
			return err
		}
		fmt.Println()
	}
	return nil
}
