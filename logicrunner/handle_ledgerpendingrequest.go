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

	watermillMsg "github.com/ThreeDotsLabs/watermill/message"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
)

type GetLedgerPendingRequest struct {
	dep *Dependencies

	Message *watermillMsg.Message
}

func (p *GetLedgerPendingRequest) Present(ctx context.Context, f flow.Flow) error {
	ctx, span := instracer.StartSpan(ctx, "LogicRunner.getLedgerPendingRequest")
	defer span.End()

	lr := p.dep.lr
	ref := Ref{}.FromSlice(p.Message.Payload)

	broker := lr.StateStorage.GetExecutionState(ref)
	if broker == nil {
		return nil
	}

	broker.executionState.getLedgerPendingMutex.Lock()
	defer broker.executionState.getLedgerPendingMutex.Unlock()

	proc := &UnsafeGetLedgerPendingRequest{
		broker:     broker,
		dep:        p.dep,
		hasPending: false,
	}

	err := f.Procedure(ctx, proc, true)
	if err != nil {
		inslogger.FromContext(ctx).Debug("GetLedgerPendingRequest.Present err: ", err)
		return nil
	}

	if !proc.hasPending {
		return nil
	}

	broker.StartProcessorIfNeeded(ctx)
	return nil
}

type UnsafeGetLedgerPendingRequest struct {
	dep        *Dependencies
	broker     *ExecutionBroker
	hasPending bool
}

func (u *UnsafeGetLedgerPendingRequest) Proceed(ctx context.Context) error {
	lr := u.dep.lr
	broker := u.broker

	broker.executionState.Lock()
	if broker.HasLedgerRequest(ctx) != nil || !broker.executionState.LedgerHasMoreRequests {
		broker.executionState.Unlock()
		return nil
	}
	broker.executionState.Unlock()

	id := *broker.Ref.Record()

	requestRef, request, err := lr.ArtifactManager.GetPendingRequest(ctx, id)
	if err != nil {
		if err != insolar.ErrNoPendingRequest {
			inslogger.FromContext(ctx).Debug("GetPendingRequest failed with error")
			return nil
		}
		broker.executionState.Lock()
		defer broker.executionState.Unlock()

		select {
		case <-ctx.Done():
			inslogger.FromContext(ctx).Debug("UnsafeGetLedgerPendingRequest: pulse changed. Do nothing")
			return nil
		default:
		}

		broker.executionState.LedgerHasMoreRequests = false
		return nil
	}
	broker.executionState.Lock()
	defer broker.executionState.Unlock()

	pulse := lr.pulse(ctx).PulseNumber
	authorized, err := lr.JetCoordinator.IsAuthorized(
		ctx, insolar.DynamicRoleVirtualExecutor, id, pulse, lr.JetCoordinator.Me(),
	)
	if err != nil {
		inslogger.FromContext(ctx).Debug("Authorization failed with error in getLedgerPendingRequest")
		return nil
	}

	if !authorized {
		inslogger.FromContext(ctx).Debug("pulse changed, can't process abandoned messages for this object")
		return nil
	}

	select {
	case <-ctx.Done():
		inslogger.FromContext(ctx).Debug("UnsafeGetLedgerPendingRequest: pulse changed. Do nothing")
		return nil
	default:
	}

	u.hasPending = true
	broker.executionState.LedgerHasMoreRequests = true

	t := NewTranscript(ctx, *requestRef, *request)
	t.FromLedger = true
	broker.Prepend(ctx, true, t)

	return nil
}
