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

package handle

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/ledger/light/proc"
)

type DeactivateObject struct {
	dep     *proc.Dependencies
	message payload.Meta
	passed  bool
}

func NewDeactivateObject(dep *proc.Dependencies, msg payload.Meta, passed bool) *DeactivateObject {
	return &DeactivateObject{
		dep:     dep,
		message: msg,
		passed:  passed,
	}
}

func (s *DeactivateObject) Present(ctx context.Context, f flow.Flow) error {
	pl, err := payload.Unmarshal(s.message.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal Deactivate message")
	}
	msg, ok := pl.(*payload.Deactivate)
	if !ok {
		return fmt.Errorf("wrong request type: %T", pl)
	}

	deactivateVirt := record.Virtual{}
	err = deactivateVirt.Unmarshal(msg.Record)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal Deactivate.Record record")
	}

	deact := record.Unwrap(&deactivateVirt)
	deactivate, ok := deact.(*record.Deactivate)
	if !ok {
		return fmt.Errorf("wrong deactivate record type: %T", deact)
	}

	resultVirt := record.Virtual{}
	err = resultVirt.Unmarshal(msg.Result)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal Deactivate.Result record")
	}

	res := record.Unwrap(&resultVirt)
	result, ok := res.(*record.Result)
	if !ok {
		return fmt.Errorf("wrong result record type: %T", res)
	}

	obj := result.Object
	if obj.IsEmpty() {
		return errors.New("object is nil")
	}

	passIfNotExecutor := !s.passed
	jet := proc.NewFetchJet(obj, flow.Pulse(ctx), s.message, passIfNotExecutor)
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
	// If we don't have it and heavy does, DeactivateObject fails because it should update light's index state.
	ensureIdx := proc.NewEnsureIndex(obj, objJetID, s.message, flow.Pulse(ctx))
	s.dep.EnsureIndex(ensureIdx)
	if err := f.Procedure(ctx, ensureIdx, false); err != nil {
		return err
	}

	setResult := proc.NewSetResult(s.message, objJetID, *result, deactivate)
	s.dep.SetResult(setResult)
	return f.Procedure(ctx, setResult, false)
}
