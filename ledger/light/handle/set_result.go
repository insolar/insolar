// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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
