// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package thread

import (
	"context"
	"fmt"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
	"go.opencensus.io/stats"
	"go.opencensus.io/tag"

	"github.com/insolar/insolar/instrumentation/inslogger"

	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/instrumentation/insmetrics"
	"github.com/insolar/insolar/instrumentation/instracer"
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
		inslogger.FromContext(ctx).Panic("procedure called with nil procedure")
	}

	var procName string
	procStringer, ok := proc.(fmt.Stringer)
	if ok {
		procName = procStringer.String()
	} else {
		procName = fmt.Sprintf("%T", proc)
	}

	ctx, span := instracer.StartSpan(ctx, procName)
	span.SetTag("type", "flow_proc")
	defer span.Finish()

	start := time.Now()
	err := func() error {
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
	}()
	duration := time.Since(start)

	result := "ok"
	if err != nil {
		if err == flow.ErrCancelled {
			result = "cancelled"
		} else {
			result = "error"
		}
		instracer.AddError(span, err)
	}
	mctx := insmetrics.ChangeTags(ctx,
		tag.Insert(tagProcedureName, procName),
		tag.Insert(tagResult, result),
	)
	stats.Record(mctx, procCallTime.M(float64(duration.Nanoseconds())/1e6))

	return err
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
