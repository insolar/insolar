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

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/insolar/insolar/insolar/belt"
)

type Flow struct {
	cancel     <-chan struct{}
	procedures map[belt.Procedure]chan bool
	message    *message.Message
}

func NewFlow(msg *message.Message, cancel <-chan struct{}) *Flow {
	return &Flow{
		cancel:     cancel,
		procedures: map[belt.Procedure]chan bool{},
		message:    msg,
	}
}

func (c *Flow) Jump(to belt.Handle) {
	panic(cancelPanic{migrateTo: to})
}

func (c *Flow) Yield(migrate belt.Handle, p belt.Procedure) bool {
	if p == nil && migrate == nil {
		panic(cancelPanic{})
	}

	if p == nil {
		<-c.cancel
		panic(cancelPanic{migrateTo: migrate})
	}

	var done bool
	select {
	case <-c.cancel:
		panic(cancelPanic{migrateTo: migrate})
	case done = <-c.proceed(p):
		return done
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
				// TODO: should probably log panic and move on (don't re-panic).
				panic(r)
			}
		}
	}()
	h(ctx, c)
}

func (c *Flow) proceed(a belt.Procedure) <-chan bool {
	if d, ok := c.procedures[a]; ok {
		return d
	}

	done := make(chan bool)
	c.procedures[a] = done
	done <- a.Proceed(context.TODO())
	return done
}
