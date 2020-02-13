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
	"github.com/insolar/insolar/ledger/light/proc"
	"github.com/pkg/errors"
)

type GetPendings struct {
	dep    *proc.Dependencies
	meta   payload.Meta
	passed bool
}

func NewGetPendings(dep *proc.Dependencies, meta payload.Meta, passed bool) *GetPendings {
	return &GetPendings{
		dep:    dep,
		meta:   meta,
		passed: passed,
	}
}

func (s *GetPendings) Present(ctx context.Context, f flow.Flow) error {
	pl, err := payload.Unmarshal(s.meta.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal GetPendings message")
	}
	msg, ok := pl.(*payload.GetPendings)
	if !ok {
		return fmt.Errorf("wrong request type: %T", pl)
	}

	passIfNotExecutor := !s.passed
	jet := proc.NewFetchJet(msg.ObjectID, flow.Pulse(ctx), s.meta, passIfNotExecutor)
	s.dep.FetchJet(jet)
	if err := f.Procedure(ctx, jet, true); err != nil {
		if err == proc.ErrNotExecutor && passIfNotExecutor {
			return nil
		}
		return err
	}

	objJetID := jet.Result.Jet

	hot := proc.NewWaitHot(objJetID, flow.Pulse(ctx), s.meta)
	s.dep.WaitHot(hot)
	if err := f.Procedure(ctx, hot, false); err != nil {
		return err
	}

	getPendings := proc.NewGetPendings(s.meta, msg.ObjectID, int(msg.Count), msg.SkipRequestRefs)
	s.dep.GetPendings(getPendings)
	if err := f.Procedure(ctx, getPendings, false); err != nil {
		return err
	}

	return nil
}
