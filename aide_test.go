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
	"reflect"
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

func TestInstance_Run(t *testing.T) {
	stage1 := NewStage("test1")
	stage2 := NewStage("test2")
	stage3 := NewStage("test3").AddSteps(
		StepFunc(func(sc *StepContext) {
			sc.Log("OK")
		}).Step("step1"),
	)
	stage4 := NewStage("test4").AddSteps(
		StepFunc(func(sc *StepContext) {
			sc.Exit()
			sc.Log("failed")
		}).Step("step1"),
	)

	type fields struct {
		stages []Stage
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "case 1: nil",
			fields: fields{
				stages: nil,
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
		{
			name: "case 2: empty",
			fields: fields{
				stages: []Stage{
					*stage1,
					*stage2,
				},
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
		{
			name: "case 3: normal",
			fields: fields{
				stages: []Stage{
					*stage1,
					*stage2,
					*stage3,
				},
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
		{
			name: "case 3: step fail",
			fields: fields{
				stages: []Stage{
					*stage1,
					*stage2,
					*stage4,
				},
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//i := &Instance{
			//	instance: tt.fields.stages,
			//}
			//if err := i.Run(tt.args.ctx); (err != nil) != tt.wantErr {
			//	t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			//}
		})
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *Instance
	}{
		{
			name: "case",
			want: &Instance{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
