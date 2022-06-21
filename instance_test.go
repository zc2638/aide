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
	"testing"
)

func TestInstance_AddStages(t *testing.T) {
	stage1 := NewStage("test1")
	stage2 := NewStage("test2")

	type args struct {
		stages []*Stage
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "case 1: nil",
			args: args{
				stages: nil,
			},
		},
		{
			name: "case 2: one",
			args: args{
				stages: []*Stage{stage1},
			},
		},
		{
			name: "case 3: multi",
			args: args{
				stages: []*Stage{stage1, stage2},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ins := New()
			ins.AddStages(tt.args.stages...)
		})
	}
}
