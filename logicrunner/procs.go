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
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/pkg/errors"
)

// ------------- CheckOurRole

type CheckOurRole struct {
	msg  insolar.Message
	role insolar.DynamicRole

	Dep struct {
		lr *LogicRunner
	}
}

func (ch *CheckOurRole) Proceed(ctx context.Context) error {
	// TODO do map of supported objects for pulse, go to jetCoordinator only if map is empty for ref
	target := ch.msg.DefaultTarget()
	isAuthorized, err := ch.Dep.lr.JetCoordinator.IsAuthorized(
		// TODO: change ch.Dep.lr.pulse(ctx).PulseNumber -> flow.Pulse(ctx)
		ctx, ch.role, *target.Record(), ch.Dep.lr.pulse(ctx).PulseNumber, ch.Dep.lr.JetCoordinator.Me(),
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

	Dep struct {
		ArtifactManager artifacts.Client
	}
}

func NewRegisterRequest(parcel insolar.Parcel, dep *Dependencies) *RegisterRequest {
	return &RegisterRequest{
		parcel: parcel,
		Dep: struct{ ArtifactManager artifacts.Client }{
			ArtifactManager: dep.lr.ArtifactManager,
		},
		result: make(chan *Ref, 1),
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

	obj := r.parcel.Message().(message.IBaseLogicMessage).GetReference()
	id, err := r.Dep.ArtifactManager.RegisterRequest(ctx, obj, r.parcel)
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

	Dep struct {
		ArtifactManager artifacts.Client
	}
}

func (c *ClarifyPendingState) Proceed(ctx context.Context) error {
	c.es.Lock()
	if c.es.pending != message.PendingUnknown {
		c.es.Unlock()
		return nil
	}

	if c.parcel != nil && c.parcel.Type() != insolar.TypeCallMethod {
		c.es.Unlock()
		c.es.pending = message.NotPending
		return nil
	}
	c.es.Unlock()

	c.es.HasPendingCheckMutex.Lock()
	defer c.es.HasPendingCheckMutex.Unlock()

	c.es.Lock()

	if c.es.pending != message.PendingUnknown {
		c.es.Unlock()
		return nil
	}

	c.es.Unlock()

	has, err := c.Dep.ArtifactManager.HasPendingRequests(ctx, c.es.Ref)
	if err != nil {
		return err
	}

	c.es.Lock()
	if c.es.pending == message.PendingUnknown {
		if has {
			c.es.pending = message.InPending
		} else {
			c.es.pending = message.NotPending
		}
	}
	c.es.Unlock()

	return nil
}
