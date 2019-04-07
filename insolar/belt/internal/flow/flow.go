///
// Copyright 2019 Insolar Technologies GmbH
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
///

package flow

import (
	"context"

	"github.com/insolar/insolar/insolar/belt"
	"github.com/insolar/insolar/insolar/belt/bus"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

type Flow struct {
	controller *Controller
	procedures map[belt.Procedure]chan struct{}
	message    bus.Message
	context    context.Context
}

func NewFlow(ctx context.Context, msg bus.Message, controller *Controller) *Flow {
	return &Flow{
		controller: controller,
		procedures: map[belt.Procedure]chan struct{}{},
		message:    msg,
		context:    ctx,
	}
}

func (c *Flow) Jump(to belt.Handle) {
	panic(cancelPanic{migrateTo: to})
}

func (c *Flow) Handle(ctx context.Context, handle belt.Handle) {
	c.handle(ctx, handle)
}

func (c *Flow) Yield(migrate belt.Handle, p belt.Procedure) {
	if p == nil && migrate == nil {
		panic(cancelPanic{})
	}

	if p == nil {
		<-c.controller.Cancel()
		panic(cancelPanic{migrateTo: migrate})
	}

	select {
	case <-c.controller.Cancel():
		panic(cancelPanic{migrateTo: migrate})
	case <-c.proceed(p):
	}
}

// =====================================================================================================================

func (c *Flow) Run(ctx context.Context, h belt.Handle) error {
	c.handle(ctx, h)
	return nil
}

// =====================================================================================================================

type cancelPanic struct {
	migrateTo belt.Handle
}

func (c *Flow) handle(ctx context.Context, h belt.Handle) {
	defer func() {
		if r := recover(); r != nil {
			if cancel, ok := r.(cancelPanic); ok {
				if cancel.migrateTo != nil {
					c.handle(ctx, cancel.migrateTo)
				}
			} else {
				inslogger.FromContext(ctx).Panic(r)
			}
		}
	}()
	h(ctx, c)
}

func (c *Flow) proceed(a belt.Procedure) <-chan struct{} {
	if d, ok := c.procedures[a]; ok {
		return d
	}

	done := make(chan struct{})
	c.procedures[a] = done
	a.Proceed(c.context)
	close(done)
	return done
}
