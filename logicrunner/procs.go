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
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/logicrunner/artifacts"
)

// ------------- CheckOurRole

type CheckOurRole struct {
	msg  insolar.Message
	role insolar.DynamicRole

	lr *LogicRunner
}

func (ch *CheckOurRole) Proceed(ctx context.Context) error {
	// TODO do map of supported objects for pulse, go to jetCoordinator only if map is empty for ref
	target := ch.msg.DefaultTarget()
	isAuthorized, err := ch.lr.JetCoordinator.IsAuthorized(
		// TODO: change ch.Dep.lr.pulse(ctx).PulseNumber -> flow.Pulse(ctx)
		ctx, ch.role, *target.Record(), ch.lr.pulse(ctx).PulseNumber, ch.lr.JetCoordinator.Me(),
	)
	if err != nil {
		return errors.Wrap(err, "authorization failed with error")
	}
	if !isAuthorized {
		return errors.New("can't execute this object")
	}
	return nil
}

// ------------- RegisterRequest

type RegisterRequest struct {
	parcel insolar.Parcel

	result chan *Ref

	ArtifactManager artifacts.Client
}

func NewRegisterRequest(parcel insolar.Parcel, dep *Dependencies) *RegisterRequest {
	return &RegisterRequest{
		parcel:          parcel,
		ArtifactManager: dep.lr.ArtifactManager,
		result:          make(chan *Ref, 1),
	}
}

func (r *RegisterRequest) setResult(result *Ref) { // nolint
	r.result <- result
}

// getResult is blocking
func (r *RegisterRequest) getResult() *Ref { // nolint
	return <-r.result
}

func (r *RegisterRequest) Proceed(ctx context.Context) error {
	ctx, span := instracer.StartSpan(ctx, "RegisterRequest.Proceed")
	defer span.End()

	obj := r.parcel.Message().(*message.CallMethod).GetReference()
	id, err := r.ArtifactManager.RegisterRequest(ctx, obj, r.parcel)
	if err != nil {
		return err
	}

	res := obj
	res.SetRecord(*id)

	r.setResult(&res)
	return nil
}

// ------------- ClarifyPendingState

type ClarifyPendingState struct {
	es     *ExecutionState
	parcel insolar.Parcel

	ArtifactManager artifacts.Client
}

func (c *ClarifyPendingState) Proceed(ctx context.Context) error {
	es := c.es
	parcel := c.parcel
	es.Lock()
	if es.pending != message.PendingUnknown {
		es.Unlock()
		return nil
	}

	if parcel != nil {
		if parcel.Type() != insolar.TypeCallMethod {
			es.Unlock()
			es.pending = message.NotPending
			return nil
		}

		msg := parcel.Message().(*message.CallMethod)
		if msg.CallType != message.CTMethod {
			es.Unlock()
			es.pending = message.NotPending
			return nil
		}
	}

	es.Unlock()

	es.HasPendingCheckMutex.Lock()
	defer es.HasPendingCheckMutex.Unlock()

	es.Lock()
	if es.pending != message.PendingUnknown {
		es.Unlock()
		return nil
	}
	es.Unlock()

	has, err := c.ArtifactManager.HasPendingRequests(ctx, es.Ref)
	if err != nil {
		return err
	}

	es.Lock()
	if es.pending == message.PendingUnknown {
		if has {
			es.pending = message.InPending
		} else {
			es.pending = message.NotPending
		}
	}
	es.Unlock()

	return nil
}
