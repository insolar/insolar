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
	pl, err := payload.Unmarshal(s.message.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal SetOutgoingRequest message")
	}
	msg, ok := pl.(*payload.SetOutgoingRequest)
	if !ok {
		return fmt.Errorf("wrong request type: %T", pl)
	}

	virtual := msg.Request
	rec := record.Unwrap(&virtual)
	request, ok := rec.(*record.OutgoingRequest)
	if !ok {
		return fmt.Errorf("wrong request type: %T", rec)
	}

	objectID := *request.AffinityRef().GetLocal()
	if objectID.IsEmpty() {
		return errors.New("object is nil")
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
	jet := proc.NewFetchJet(objectID, flow.Pulse(ctx), s.message, passIfNotExecutor)
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
	// If we don't have it and heavy does, SetOutgoingRequest fails because it should update light's index state.
	getIndex := proc.NewEnsureIndex(objectID, objJetID, s.message, flow.Pulse(ctx))
	s.dep.EnsureIndex(getIndex)
	if err := f.Procedure(ctx, getIndex, false); err != nil {
		return err
	}

	setRequest := proc.NewSetRequest(s.message, request, reqID, objJetID)
	s.dep.SetRequest(setRequest)
	return f.Procedure(ctx, setRequest, false)
}
