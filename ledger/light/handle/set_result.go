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

	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/ledger/light/proc"
	"github.com/pkg/errors"
)

type SetResult struct {
	dep     *proc.Dependencies
	message payload.Meta
	passed  bool
}

func NewSetResult(dep *proc.Dependencies, msg payload.Meta, passed bool) *SetResult {
	return &SetResult{
		dep:     dep,
		message: msg,
		passed:  passed,
	}
}

func (s *SetResult) Present(ctx context.Context, f flow.Flow) error {
	pl, err := payload.Unmarshal(s.message.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal SetResult message")
	}
	msg, ok := pl.(*payload.SetResult)
	if !ok {
		return fmt.Errorf("wrong request type: %T", pl)
	}

	virtual := record.Virtual{}
	err = virtual.Unmarshal(msg.Result)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal Result record")
	}

	rec := record.Unwrap(&virtual)
	result, ok := rec.(*record.Result)
	if !ok {
		return fmt.Errorf("wrong result type: %T", rec)
	}

	if result.Object.IsEmpty() {
		return errors.New("object is nil")
	}

	passIfNotExecutor := !s.passed
	jet := proc.NewFetchJet(result.Object, flow.Pulse(ctx), s.message, passIfNotExecutor)
	s.dep.FetchJet(jet)
	if err := f.Procedure(ctx, jet, true); err != nil {
		if err == proc.ErrNotExecutor && passIfNotExecutor {
			return nil
		}
		return err
	}
	jetID := jet.Result.Jet

	hot := proc.NewWaitHot(jetID, flow.Pulse(ctx), s.message)
	s.dep.WaitHot(hot)
	if err := f.Procedure(ctx, hot, false); err != nil {
		return err
	}

	// To ensure, that we have the index. Because index can be on a heavy node.
	// If we don't have it and heavy does, SetResult fails because it should update light's index state.
	idx := proc.NewEnsureIndex(result.Object, jetID, s.message, flow.Pulse(ctx))
	s.dep.EnsureIndex(idx)
	if err := f.Procedure(ctx, idx, false); err != nil {
		return errors.Wrap(err, "can't get index")
	}

	setResult := proc.NewSetResult(s.message, jetID, *result, nil)
	s.dep.SetResult(setResult)
	return f.Procedure(ctx, setResult, false)
}
