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
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

type HandleCall struct {
	dep *Dependencies

	Message bus.Message
}

func (h *HandleCall) executeActual(
	ctx context.Context,
	parcel insolar.Parcel,
	msg message.IBaseLogicMessage,
	f flow.Flow) (insolar.Reply, error) {

	lr := h.dep.lr
	ref := msg.GetReference()
	os := lr.UpsertObjectState(ref)
	log.Info("After UpsertObjectState")

	os.Lock()
	if os.ExecutionState == nil {
		os.ExecutionState = &ExecutionState{
			Ref:       ref,
			Queue:     make([]ExecutionQueueElement, 0),
			Behaviour: &ValidationSaver{lr: lr, caseBind: NewCaseBind()},
		}
	}
	es := os.ExecutionState
	os.Unlock()

	// ExecutionState should be locked between CheckOurRole and
	// appending ExecutionQueueElement to the queue to prevent a race condition.
	// Otherwise it's possible that OnPulse will clean up the queue and set
	// ExecutionState.Pending to NotPending. Execute will add an element to the
	// queue afterwards. In this case cross-pulse execution will break.
	es.Lock()

	procCheckRole := CheckOurRole{
		msg:  msg,
		role: insolar.DynamicRoleVirtualExecutor,
		Dep: struct{ JetCoordinator insolar.JetCoordinator }{
			JetCoordinator: lr.JetCoordinator,
		},
	}

	if err := f.Procedure(ctx, &procCheckRole); err != nil {
		es.Unlock()
		// TODO: check if error is ErrCancelled
		return nil, errors.Wrap(err, "[ Execute ] can't play role")
	}

	if lr.CheckExecutionLoop(ctx, es, parcel) {
		es.Unlock()
		return nil, os.WrapError(nil, "loop detected")
	}
	es.Unlock()

	procRegisterRequest := NewRegisterRequest(parcel, h.dep)

	if err := f.Procedure(ctx, procRegisterRequest); err != nil {
		// TODO: check if error is ErrCancelled
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
		es:     es,
		parcel: parcel,
		Dep: struct{ ArtifactManager artifacts.Client }{
			ArtifactManager: lr.ArtifactManager,
		},
	}

	if err := f.Procedure(ctx, &procClarifyPendingState); err != nil {
		// TODO: check if error is ErrCancelled
		return nil, err
	}

	s := StartQueueProcessorIfNeeded{
		es:  es,
		dep: h.dep,
		ref: &ref,
	}
	if err := f.Handle(ctx, s.Present); err != nil {
		inslogger.FromContext(ctx).Warn("[ executeActual ] StartQueueProcessorIfNeeded returns error: ", err)
	}

	return &reply.RegisterRequest{
		Request: *request,
	}, nil

}

func (h *HandleCall) Present(ctx context.Context, f flow.Flow) error {
	parcel := h.Message.Parcel
	ctx = loggerWithTargetID(ctx, parcel)
	inslogger.FromContext(ctx).Debug("LogicRunner.Execute starts ...")

	msg, ok := parcel.Message().(message.IBaseLogicMessage)
	if !ok {
		return errors.New("Execute( ! message.IBaseLogicMessage )")
	}

	ctx, span := instracer.StartSpan(ctx, "LogicRunner.Execute")
	span.AddAttributes(
		trace.StringAttribute("msg.Type", msg.Type().String()),
	)
	defer span.End()

	r := bus.Reply{}
	r.Reply, r.Err = h.executeActual(ctx, parcel, msg, f)

	h.Message.ReplyTo <- r
	return nil

}
