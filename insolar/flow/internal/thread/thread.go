//
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
//

package thread

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar/flow"
)

type Thread struct {
	controller *Controller
	cancel     <-chan struct{}
	canBegin   <-chan struct{}
	procedures map[flow.Procedure]*result
	message    *message.Message
	migrated   bool
}

type result struct {
	done chan struct{}
	err  error
}

// NewThread creates a new Thread instance. Thread implements the Flow interface.
func NewThread(msg *message.Message, controller *Controller) *Thread {
	return &Thread{
		controller: controller,
		cancel:     controller.Cancel(),
		canBegin:   controller.CanBegin(),
		procedures: map[flow.Procedure]*result{},
		message:    msg,
	}
}

func (f *Thread) Handle(ctx context.Context, handle flow.Handle) error {
	return handle(ctx, f)
}

func (f *Thread) Procedure(ctx context.Context, proc flow.Procedure, cancel bool) error {
	if proc == nil {
		panic("procedure called with nil procedure")
	}

	if !cancel {
		res := f.procedure(ctx, proc)
		<-res.done
		return res.err
	}

	if f.cancelled() {
		return flow.ErrCancelled
	}

	ctx, cl := context.WithCancel(ctx)
	res := f.procedure(ctx, proc)
	select {
	case <-f.cancel:
		cl()
		return flow.ErrCancelled
	case <-res.done:
		cl()
		return res.err
	}
}

func (f *Thread) Migrate(ctx context.Context, to flow.Handle) error {
	if f.migrated {
		return errors.New("migrate called on migrated flow")
	}

	<-f.canBegin
	f.migrated = true
	subFlow := NewThread(f.message, f.controller)
	return to(ctx, subFlow)
}

func (f *Thread) Continue(context.Context) {
	<-f.canBegin
	f.canBegin = f.controller.CanBegin()
}

func (f *Thread) Run(ctx context.Context, h flow.Handle) error {
	return h(ctx, f)
}

func (f *Thread) procedure(ctx context.Context, proc flow.Procedure) *result {
	if res, ok := f.procedures[proc]; ok {
		return res
	}

	res := &result{
		done: make(chan struct{}),
		err:  nil,
	}
	f.procedures[proc] = res
	go func() {
		res.err = proc.Proceed(ctx)
		close(res.done)
	}()
	return res
}

func (f *Thread) cancelled() bool {
	select {
	case <-f.cancel:
		return true
	default:
		return false
	}
}
