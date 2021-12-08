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

// contextKey is a value for use with context.WithValue. It's used as
// a pointer, so it fits in an interface{} without allocation. This technique
// for defining context keys was copied from Go 1.7's new use of context in net/http.
type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "aide context value " + k.name
}

var (
	// StepCtxKey is the context.Context key to store the step context.
	StepCtxKey = &contextKey{"StepContext"}
	// StepTotalKey is the context.Context key to store the total number of steps.
	StepTotalKey = &contextKey{"StepTotal"}
	// StepLoggerKey is the context.Context key to store the logger.
	StepLoggerKey = &contextKey{"StepLogger"}
)
