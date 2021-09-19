// Package aide
// Created by zc on 2021/9/18.
package aide

import (
	"context"
	"reflect"
	"testing"
)

func TestStepContext_Context(t *testing.T) {
	type fields struct {
		ctx      context.Context
		level    LogLevel
		exitCode int32
		message  string
	}
	tests := []struct {
		name   string
		fields fields
		want   context.Context
	}{
		{
			name: "case",
			fields: fields{
				ctx: context.Background(),
			},
			want: context.Background(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &StepContext{
				ctx:      tt.fields.ctx,
				level:    tt.fields.level,
				exitCode: tt.fields.exitCode,
				message:  tt.fields.message,
			}
			if got := c.Context(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Context() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStepContext_Exit(t *testing.T) {
	type args struct {
		code int32
	}
	tests := []struct {
		name string
		args args
		want int32
	}{
		{
			name: "case 1: nil",
			args: args{},
			want: 0,
		},
		{
			name: "case 2: 1",
			args: args{
				code: 1,
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &StepContext{}
			c.Exit(tt.args.code)
			if c.exitCode != tt.want {
				t.Errorf("StepContext_Exit() = %v, want %v", c.exitCode, tt.want)
			}
		})
	}
}

func TestStepContext_SetLevel(t *testing.T) {
	type args struct {
		level LogLevel
	}
	tests := []struct {
		name string
		args args
		want LogLevel
	}{
		{
			name: "case 1: nil",
			args: args{},
			want: 0,
		},
		{
			name: "case 2: info",
			args: args{
				level: InfoLevel,
			},
			want: InfoLevel,
		},
		{
			name: "case 3: error",
			args: args{
				level: ErrorLevel,
			},
			want: ErrorLevel,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &StepContext{}
			c.WithLevel(tt.args.level)
			if c.level != tt.want {
				t.Errorf("StepContext_SetLevel() = %v, want %v", c.level, tt.want)
			}
		})
	}
}

func TestStepContext_WithContext(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want context.Context
	}{
		{
			name: "case 1: nil",
			args: args{},
			want: nil,
		},
		{
			name: "case 2: empty",
			args: args{
				ctx: context.Background(),
			},
			want: context.Background(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &StepContext{}
			c.WithContext(tt.args.ctx)
			if c.ctx != tt.want {
				t.Errorf("StepContext_WithContext() = %v, want %v", c.ctx, tt.want)
			}
		})
	}
}

func TestStepContext_Write(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "case 1: nil",
			args: args{},
			want: nil,
		},
		{
			name: "case 2: normal",
			args: args{
				b: []byte("test"),
			},
			want: []byte("test"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &StepContext{}
			c.Write(tt.args.b)
			if c.message != string(tt.want) {
				t.Errorf("StepContext_Write() = %v, want %v", c.ctx, tt.want)
			}
		})
	}
}

func TestStepContext_WriteString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "case 1: nil",
			args: args{},
			want: "",
		},
		{
			name: "case 2: normal",
			args: args{
				s: "test",
			},
			want: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &StepContext{}
			c.WriteString(tt.args.s)
			if c.message != tt.want {
				t.Errorf("StepContext_WriteString() = %v, want %v", c.ctx, tt.want)
			}
		})
	}
}

func TestStep_run(t *testing.T) {
	type fields struct {
		name string
		srf  StepFunc
	}
	type args struct {
		sc *StepContext
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "case",
			fields: fields{
				name: "test",
				srf: func(sc *StepContext) {
					sc.WriteString("ok")
				},
			},
			args: args{
				sc: &StepContext{
					ctx: context.Background(),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Step{
				name: tt.fields.name,
				srf:  tt.fields.srf,
			}
			if err := s.run(tt.args.sc); (err != nil) != tt.wantErr {
				t.Errorf("run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
