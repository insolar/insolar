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

type UpdateObject struct {
	dep     *proc.Dependencies
	message payload.Meta
	passed  bool
}

func NewUpdateObject(dep *proc.Dependencies, msg payload.Meta, passed bool) *UpdateObject {
	return &UpdateObject{
		dep:     dep,
		message: msg,
		passed:  passed,
	}
}

func (s *UpdateObject) Present(ctx context.Context, f flow.Flow) error {
	msg := payload.Update{}
	err := msg.Unmarshal(s.message.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal Update message")
	}

	updateVirt := record.Virtual{}
	err = updateVirt.Unmarshal(msg.Record)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal Update.Record record")
	}

	upd := record.Unwrap(&updateVirt)
	update, ok := upd.(*record.Amend)
	if !ok {
		return fmt.Errorf("wrong update record type: %T", upd)
	}

	resultVirt := record.Virtual{}
	err = resultVirt.Unmarshal(msg.Result)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal Update.Result record")
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
	// If we don't have it and heavy does, UpdateObject fails because it should update light's index state
	getIndex := proc.NewEnsureIndex(obj, objJetID, s.message, flow.Pulse(ctx))
	s.dep.EnsureIndex(getIndex)
	if err := f.Procedure(ctx, getIndex, false); err != nil {
		return err
	}

	setResult := proc.NewSetResult(s.message, objJetID, *result, update)
	s.dep.SetResult(setResult)
	return f.Procedure(ctx, setResult, false)
}
