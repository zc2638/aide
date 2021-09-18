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

func TestNewStage(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want *Stage
	}{
		{
			name: "case 1: nil",
			args: args{
				name: "",
			},
			want: &Stage{},
		},
		{
			name: "case 2: name only",
			args: args{
				name: "test",
			},
			want: &Stage{
				name: "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStage(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStage_AddSteps(t *testing.T) {
	step1 := &Step{name: "step1"}
	step2 := &Step{name: "step2"}

	type args struct {
		name  string
		steps []*Step
	}
	tests := []struct {
		name string
		args args
		want *Stage
	}{
		{
			name: "case 1: nil",
			args: args{
				steps: nil,
			},
			want: NewStage(""),
		},
		{
			name: "case 2: one",
			args: args{
				name:  "one",
				steps: []*Step{step1},
			},
			want: &Stage{
				name:  "one",
				steps: []Step{*step1},
			},
		},
		{
			name: "case 3: multi",
			args: args{
				name:  "multi",
				steps: []*Step{step1, step2},
			},
			want: &Stage{
				name:  "multi",
				steps: []Step{*step1, *step2},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewStage(tt.args.name)
			if got := s.AddSteps(tt.args.steps...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddSteps() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStage_run(t *testing.T) {
	step1 := &Step{name: "step1"}
	step2 := &Step{name: "step2"}
	step3 := StepFunc(func(sc *StepContext) {
		sc.WriteString("ok")
	}).Step("step1")
	step4 := StepFunc(func(sc *StepContext) {
		sc.Exit(1)
		sc.WriteString("failed")
	}).Step("step1")

	type fields struct {
		name  string
		steps []Step
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
				steps: nil,
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
		{
			name: "case 2: empty",
			fields: fields{
				steps: []Step{
					*step1,
					*step2,
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
				steps: []Step{
					*step1,
					*step2,
					*step3,
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
				steps: []Step{
					*step1,
					*step2,
					*step4,
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
			s := &Stage{
				name:  tt.fields.name,
				steps: tt.fields.steps,
			}
			if err := s.run(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
