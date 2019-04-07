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
	"github.com/pkg/errors"
)

type Flow struct {
	controller *Controller
	cancel     <-chan struct{}
	procedures map[belt.Procedure]chan error
	message    bus.Message
	migrated   bool
}

func NewFlow(msg bus.Message, controller *Controller) *Flow {
	return &Flow{
		controller: controller,
		cancel:     controller.Cancel(),
		procedures: map[belt.Procedure]chan error{},
		message:    msg,
	}
}

func (f *Flow) Handle(ctx context.Context, handle belt.Handle) error {
	if f.cancelled() {
		return belt.ErrCancelled
	}

	err := handle(ctx, f)
	if err != nil {
		return err
	}

	if f.cancelled() {
		return belt.ErrCancelled
	}

	return nil
}

func (f *Flow) Procedure(ctx context.Context, p belt.Procedure) error {
	if f.cancelled() {
		return belt.ErrCancelled
	}

	if p == nil {
		return errors.New("procedure called with nil procedure")
	}

	ctx, cancel := context.WithCancel(ctx)
	var err error
	select {
	case <-f.cancel:
		cancel()
		return belt.ErrCancelled
	case err = <-f.procedure(ctx, p):
	}

	return err
}

func (f *Flow) Migrate(ctx context.Context, to belt.Handle) error {
	if f.migrated {
		return errors.New("migrate called on migrated flow")
	}

	<-f.cancel
	f.migrated = true
	subFlow := NewFlow(f.message, f.controller)
	return to(ctx, subFlow)
}

// =====================================================================================================================

func (f *Flow) Run(ctx context.Context, h belt.Handle) error {
	return h(ctx, f)
}

// =====================================================================================================================

func (f *Flow) procedure(ctx context.Context, a belt.Procedure) <-chan error {
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

func (f *Flow) cancelled() bool {
	select {
	case <-f.cancel:
		return true
	default:
		return false
	}
}
