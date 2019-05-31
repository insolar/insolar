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

package logicrunner

import (
	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/flow/bus"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

type HandleCall struct {
	dep *Dependencies

	Message bus.Message
}

func (h *HandleCall) handleActual(
	ctx context.Context,
	parcel insolar.Parcel,
	msg *message.CallMethod,
	f flow.Flow,
) (insolar.Reply, error) {

	lr := h.dep.lr
	ref := msg.GetReference()
	os := lr.UpsertObjectState(ref)

	os.Lock()
	if os.ExecutionState == nil {
		os.ExecutionState = &ExecutionState{
			Ref:   ref,
			Queue: make([]ExecutionQueueElement, 0),
		}
	}
	es := os.ExecutionState
	os.Unlock()

	es.Lock()

	procCheckRole := CheckOurRole{
		msg:  msg,
		role: insolar.DynamicRoleVirtualExecutor,
		lr:   lr,
	}

	if err := f.Procedure(ctx, &procCheckRole, true); err != nil {
		es.Unlock()
		if err == flow.ErrCancelled {
			return nil, err // message bus will retry on the calling side
		}
		return nil, errors.Wrap(err, "[ handleActual ] can't play role")
	}

	if lr.CheckExecutionLoop(ctx, es, parcel) {
		es.Unlock()
		return nil, os.WrapError(nil, "loop detected")
	}
	es.Unlock()

	// RegisterRequest is an external, slow call to the LME thus we have to
	// unlock ExecutionState during the call.

	procRegisterRequest := NewRegisterRequest(parcel, h.dep)

	if err := f.Procedure(ctx, procRegisterRequest, true); err != nil {
		if err == flow.ErrCancelled {
			return nil, err // message bus will automatically retry on the calling side
		}
		return nil, os.WrapError(err, "[ Execute ] can't create request")
	}
	request := procRegisterRequest.getResult()

	es.Lock()
	qElement := ExecutionQueueElement{
		ctx:     ctx,
		parcel:  parcel,
		request: request,
	}

	es.Queue = append(es.Queue, qElement)
	es.Unlock()

	procClarifyPendingState := ClarifyPendingState{
		es:              es,
		parcel:          parcel,
		ArtifactManager: lr.ArtifactManager,
	}

	if err := f.Procedure(ctx, &procClarifyPendingState, true); err != nil {
		if err == flow.ErrCancelled {
			// If the flow has canceled during ClarifyPendingState there are two possibilities.
			// 1. It's possible that we started to execute ClarifyPendingState, the pulse has
			// changed and the execution queue was sent to the next executor in OnPulse method.
			// This means that the next executor already has the last queue element, it's OK.
			// 2. It's also possible that the pulse has changed after RegisterRequest but
			// before adding an item to the execution queue. In this case the queue was sent to the
			// next executor without the last item. We could just return ErrCanceled to make the
			// caller to resend the request. However this will cause a slow request deduplication
			// process on the LME side. As an optimization we decided to send a special message
			// type to the next executor. It's possible that while we send the message the pulse
			// will change once again and the receiver will be not an executor of the object anymore.
			// However in this case MessageBus will automatically resend the message to the right VE.

			// TODO: it's done to support current logic. Do it correctly when go to flow
			f.Continue(ctx)
		} else {
			return nil, err
		}
	}

	s := StartQueueProcessorIfNeeded{
		es:  es,
		dep: h.dep,
		ref: &ref,
	}
	if err := f.Handle(ctx, s.Present); err != nil {
		inslogger.FromContext(ctx).Warn("[ handleActual ] StartQueueProcessorIfNeeded returns error: ", err)
	}

	return &reply.RegisterRequest{
		Request: *request,
	}, nil

}

func (h *HandleCall) Present(ctx context.Context, f flow.Flow) error {
	parcel := h.Message.Parcel
	ctx = loggerWithTargetID(ctx, parcel)
	inslogger.FromContext(ctx).Debug("HandleCall.Present starts ...")

	msg, ok := parcel.Message().(*message.CallMethod)
	if !ok {
		return errors.New("is not CallMethod message")
	}

	ctx, span := instracer.StartSpan(ctx, "LogicRunner.Execute")
	span.AddAttributes(
		trace.StringAttribute("msg.Type", msg.Type().String()),
	)
	defer span.End()

	r := bus.Reply{}
	r.Reply, r.Err = h.handleActual(ctx, parcel, msg, f)

	h.Message.ReplyTo <- r
	return nil

}
