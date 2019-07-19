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

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/logicrunner/artifacts"
)

// ------------- CheckOurRole

type CheckOurRole struct {
	msg         insolar.Message
	role        insolar.DynamicRole
	pulseNumber insolar.PulseNumber

	lr *LogicRunner
}

var ErrCantExecute = errors.New("can't executeAndReply this object")

func (ch *CheckOurRole) Proceed(ctx context.Context) error {
	ctx, span := instracer.StartSpan(ctx, "CheckOurRole")
	defer span.End()

	// TODO do map of supported objects for pulse, go to jetCoordinator only if map is empty for ref
	target := ch.msg.DefaultTarget()
	isAuthorized, err := ch.lr.JetCoordinator.IsAuthorized(
		ctx, ch.role, *target.Record(), ch.pulseNumber, ch.lr.JetCoordinator.Me(),
	)
	if err != nil {
		return errors.Wrap(err, "authorization failed with error")
	}
	if !isAuthorized {
		return ErrCantExecute
	}
	return nil
}

// ------------- RegisterIncomingRequest

type RegisterIncomingRequest struct {
	request record.IncomingRequest

	result chan *Ref

	ArtifactManager artifacts.Client
}

func NewRegisterIncomingRequest(request record.IncomingRequest, dep *Dependencies) *RegisterIncomingRequest {
	return &RegisterIncomingRequest{
		request:         request,
		ArtifactManager: dep.lr.ArtifactManager,
		result:          make(chan *Ref, 1),
	}
}

func (r *RegisterIncomingRequest) setResult(result *Ref) { // nolint
	r.result <- result
}

// getResult is blocking
func (r *RegisterIncomingRequest) getResult() *Ref { // nolint
	return <-r.result
}

func (r *RegisterIncomingRequest) Proceed(ctx context.Context) error {
	ctx, span := instracer.StartSpan(ctx, "RegisterIncomingRequest.Proceed")
	defer span.End()

	id, err := r.ArtifactManager.RegisterIncomingRequest(ctx, &r.request)
	if err != nil {
		return err
	}

	r.setResult(insolar.NewReference(*id))
	return nil
}

// ------------- ClarifyPendingState

type ClarifyPendingState struct {
	broker  *ExecutionBroker
	request *record.IncomingRequest

	ArtifactManager artifacts.Client
}

func (c *ClarifyPendingState) Proceed(ctx context.Context) error {
	ctx, span := instracer.StartSpan(ctx, "ClarifyPendingState")
	defer span.End()

	es := &c.broker.executionState

	es.Lock()
	if es.pending != insolar.PendingUnknown {
		es.Unlock()
		return nil
	}

	if c.request != nil {
		if c.request.CallType != record.CTMethod {
			// It's considered that we are not pending except someone calls a method.
			es.pending = insolar.NotPending
			es.Unlock()
			return nil
		}
	}

	es.Unlock()

	es.HasPendingCheckMutex.Lock()
	defer es.HasPendingCheckMutex.Unlock()

	es.Lock()
	if es.pending != insolar.PendingUnknown {
		es.Unlock()
		return nil
	}
	es.Unlock()

	has, err := c.ArtifactManager.HasPendingRequests(ctx, c.broker.Ref)
	if err != nil {
		return err
	}

	es.Lock()
	if es.pending == insolar.PendingUnknown {
		if has {
			es.pending = insolar.InPending
		} else {
			es.pending = insolar.NotPending
		}
	}
	es.Unlock()

	return nil
}
