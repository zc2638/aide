// Package stage

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

package stage

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/sync/errgroup"

	"github.com/99nil/go/cycle"
	"github.com/99nil/go/sets"
)

type InstanceFunc func(c Context) error

type Instance struct {
	name string
	// Whether to enable asynchronous processing.
	async bool
	// Stage subset.
	cs []*Instance
	// Calling method before executing cs.
	pre InstanceFunc
	// Calling method after executing cs.
	sub InstanceFunc
	// The names of other stages that need to be relied upon before execution.
	relies []string
	// After executing the current stage, the name of the next stage that needs to be executed.
	// Note: other stages between the current stage and the next stage will not be executed.
	next string

	skip     bool
	skipFunc func() bool
}

func New(name string) *Instance {
	return &Instance{name: name}
}

func (ins *Instance) rename() {
	ns := sets.NewString()
	for _, c := range ins.cs {
		i := 0
		for ns.Has(c.name) {
			c.name = fmt.Sprintf("%s_%d", c.name, i)
			i++
		}
		ns.Add(c.name)
	}
}

func (ins *Instance) Len() int {
	return len(ins.cs)
}

func (ins *Instance) Skip(is bool) *Instance {
	ins.skip = is
	return ins
}

func (ins *Instance) SkipFunc(f func() bool) *Instance {
	ins.skipFunc = f
	return ins
}

// Add adds subsets
func (ins *Instance) Add(cs ...*Instance) *Instance {
	for _, c := range cs {
		ins.cs = append(ins.cs, c)
	}
	return ins
}

// SetAsync sets whether the current stage is executed asynchronously.
func (ins *Instance) SetAsync(async bool) *Instance {
	ins.async = async
	return ins
}

// SetPreFunc sets the execution method before executing the subset.
func (ins *Instance) SetPreFunc(f InstanceFunc) *Instance {
	ins.pre = f
	return ins
}

// SetSubFunc sets the execution method after executing the subset.
func (ins *Instance) SetSubFunc(f InstanceFunc) *Instance {
	ins.sub = f
	return ins
}

// RelyOn sets the names of other stages that the current stage needs to depend on.
func (ins *Instance) RelyOn(names ...string) *Instance {
	ins.relies = make([]string, 0, len(names))
	ins.relies = append(ins.relies, names...)
	return ins
}

// Goto sets the name of the next stage to be executed.
// Note: will be invalid in asynchronous stage.
func (ins *Instance) Goto(name string) *Instance {
	ins.next = name
	return ins
}

// getChildNames gets the names of all subsets.
func (ins *Instance) getChildNames() []string {
	csLen := len(ins.cs)
	if csLen == 0 {
		return nil
	}
	res := make([]string, 0, csLen)
	for _, c := range ins.cs {
		res = append(res, c.name)
	}
	return res
}

// hasLoop checks whether there is a circular dependency.
func (ins *Instance) hasLoop() bool {
	graph := cycle.New()
	for _, c := range ins.cs {
		graph.Add(c.name, c.relies...)
	}
	return graph.DetectCycles()
}

func (ins *Instance) Run(ctx context.Context) error {
	if ins.hasLoop() {
		return errors.New("dependency cycle detected")
	}
	// TODO Check for non-existent dependencies.
	sc := NewCtx(ctx)
	if err := ins.run(sc); err != nil && err != ErrStageEnd {
		return err
	}
	return nil
}

func (ins *Instance) run(sc Context) error {
	if ins.skip {
		return nil
	}
	if ins.skipFunc != nil && ins.skipFunc() {
		return nil
	}

	var err error
	if ins.pre != nil {
		sc.WithValue(NameKey, ins.name)
		err = ins.pre(sc)
		// If the exception is ErrStageSkip, skip the execution stage.
		if err == ErrStageSkip {
			return nil
		}
		if err != nil {
			return err
		}
	}

	if ins.async {
		err = ins.runAsync(sc)
	} else {
		err = ins.runSync(sc)
	}
	if err != nil {
		return err
	}

	if ins.sub != nil {
		sc.WithValue(NameKey, ins.name)
		if err := ins.sub(sc); err != nil {
			return err
		}
	}
	return nil
}

// runSync runs synchronous execution of the stage subset.
func (ins *Instance) runSync(sc Context) error {
	doneSet := sets.NewString()
	pending := ins.cs[:]
	for len(pending) > 0 {
		var next string
		wait := make([]*Instance, 0, len(pending))

		for _, c := range pending {
			if doneSet.Has(c.name) {
				continue
			}

			if next != "" && c.name != next {
				doneSet.Add(c.name)
				continue
			}
			if len(c.relies) != 0 {
				// Check whether the dependency has completed running.
				if !doneSet.HasAll(c.relies...) {
					wait = append(wait, c)
					continue
				}
			}
			if err := c.run(sc); err != nil {
				return err
			}
			doneSet.Add(c.name)
			next = c.next
		}

		pending = wait[:]
	}
	return nil
}

// runAsync runs asynchronous execution of the stage subset.
func (ins *Instance) runAsync(sc Context) error {
	doneSet := sets.NewString()
	pending := ins.cs[:]
	for len(pending) > 0 {
		wait := make([]*Instance, 0, len(pending))

		ctx := sc.Ctx()
		eg, cancelCtx := errgroup.WithContext(ctx)
		sc.WithCtx(cancelCtx)

		scCopySet := make([]Context, 0, len(pending))
		for _, c := range pending {
			if doneSet.Has(c.name) {
				continue
			}
			if len(c.relies) != 0 {
				// Check whether the dependency has completed running.
				if !doneSet.HasAll(c.relies...) {
					wait = append(wait, c)
					continue
				}
			}

			scCopy := sc.(*valueCtx).clone()
			func(c *Instance) {
				eg.Go(func() error {
					return c.run(scCopy)
				})
			}(c)
			doneSet.Add(c.name)
			scCopySet = append(scCopySet, scCopy)
		}

		// Parallel processing needs to wait for all processes to end.
		// If you want to control, please use the `Done` method of `Context`.
		if err := eg.Wait(); err != nil {
			return err
		}

		// Combine the context of two processes.
		for _, scCopy := range scCopySet {
			sc.(*valueCtx).combine(scCopy)
		}
		sc.WithCtx(ctx)
		pending = wait[:]
	}
	return nil
}
