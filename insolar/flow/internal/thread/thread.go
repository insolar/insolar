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

	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/flow/bus"
	"github.com/pkg/errors"
)

type Thread struct {
	controller *Controller
	cancel     <-chan struct{}
	procedures map[flow.Procedure]chan error
	message    bus.Message
	migrated   bool
}

// NewThread creates a new Thread instance. Thread implements the Flow interface.
func NewThread(msg bus.Message, controller *Controller) *Thread {
	return &Thread{
		controller: controller,
		cancel:     controller.Cancel(),
		procedures: map[flow.Procedure]chan error{},
		message:    msg,
	}
}

func (f *Thread) Handle(ctx context.Context, handle flow.Handle) error {
	if f.cancelled() {
		return flow.ErrCancelled
	}

	err := handle(ctx, f)
	if err != nil {
		return err
	}

	if f.cancelled() {
		return flow.ErrCancelled
	}

	return nil
}

func (f *Thread) Procedure(ctx context.Context, p flow.Procedure) error {
	if f.cancelled() {
		return flow.ErrCancelled
	}

	if p == nil {
		return errors.New("procedure called with nil procedure")
	}

	ctx, cancel := context.WithCancel(ctx)
	var err error
	select {
	case <-f.cancel:
		cancel()
		return flow.ErrCancelled
	case err = <-f.procedure(ctx, p):
		cancel()
	}

	return err
}

func (f *Thread) Migrate(ctx context.Context, to flow.Handle) error {
	if f.migrated {
		return errors.New("migrate called on migrated flow")
	}

	<-f.cancel
	f.migrated = true
	subFlow := NewThread(f.message, f.controller)
	return to(ctx, subFlow)
}

func (f *Thread) Continue(context.Context) {
	<-f.cancel
	f.cancel = f.controller.Cancel()
}

// =====================================================================================================================

func (f *Thread) Run(ctx context.Context, h flow.Handle) error {
	return h(ctx, f)
}

// =====================================================================================================================

func (f *Thread) procedure(ctx context.Context, a flow.Procedure) <-chan error {
	if d, ok := f.procedures[a]; ok {
		return d
	}

	done := make(chan error, 1)
	f.procedures[a] = done
	go func() {
		done <- a.Proceed(ctx)
	}()
	return done
}

func (f *Thread) cancelled() bool {
	select {
	case <-f.cancel:
		return true
	default:
		return false
	}
}
