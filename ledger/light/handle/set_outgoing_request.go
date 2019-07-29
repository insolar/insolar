/*
 *    Copyright 2019 Insolar Technologies
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package handle

import (
	"context"
	"fmt"

	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/ledger/light/proc"
	"github.com/pkg/errors"
)

type SetOutgoingRequest struct {
	dep     *proc.Dependencies
	message payload.Meta
	passed  bool
}

func NewSetOutgoingRequest(dep *proc.Dependencies, msg payload.Meta, passed bool) *SetOutgoingRequest {
	return &SetOutgoingRequest{
		dep:     dep,
		message: msg,
		passed:  passed,
	}
}

func (s *SetOutgoingRequest) Present(ctx context.Context, f flow.Flow) error {
	msg := payload.SetOutgoingRequest{}
	err := msg.Unmarshal(s.message.Payload)
	if err != nil {
		return errors.Wrap(err, "SetOutgoingRequest.Present: failed to unmarshal SetIncomingRequest message")
	}

	rec := record.Unwrap(&msg.Request)
	request, ok := rec.(*record.OutgoingRequest)
	if !ok {
		return fmt.Errorf("SetOutgoingRequest.Present: wrong request type: %T", rec)
	}

	if request.AffinityRef() == nil {
		return errors.New("SetOutgoingRequest.Present: object is nil")
	}

	buf, err := msg.Request.Marshal()
	if err != nil {
		return err
	}

	calc := proc.NewCalculateID(buf, flow.Pulse(ctx))
	s.dep.CalculateID(calc)
	if err := f.Procedure(ctx, calc, true); err != nil {
		return err
	}
	reqID := calc.Result.ID

	passIfNotExecutor := !s.passed
	jet := proc.NewFetchJet(*request.AffinityRef().Record(), flow.Pulse(ctx), s.message, passIfNotExecutor)
	s.dep.FetchJet(jet)
	if err := f.Procedure(ctx, jet, true); err != nil {
		if err == proc.ErrNotExecutor && passIfNotExecutor {
			return nil
		}
		return err
	}
	objJetID := jet.Result.Jet

	hot := proc.NewWaitHot(objJetID, flow.Pulse(ctx), s.message)
	s.dep.WaitHot(hot)
	if err := f.Procedure(ctx, hot, false); err != nil {
		return err
	}

	// To ensure, that we have the index. Because index can be on a heavy node.
	// If we don't have it and heavy does, SetResult fails because it should update light's index state
	getIndex := proc.NewEnsureIndex(*request.AffinityRef().Record(), objJetID, s.message)
	s.dep.EnsureIndex(getIndex)
	if err := f.Procedure(ctx, getIndex, false); err != nil {
		return err
	}

	setRequest := proc.NewSetRequest(s.message, request, reqID, objJetID)
	s.dep.SetRequest(setRequest)
	return f.Procedure(ctx, setRequest, false)
}
