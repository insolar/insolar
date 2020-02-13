// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package handle

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/light/proc"
)

type GetObject struct {
	dep *proc.Dependencies

	meta   payload.Meta
	passed bool
}

func NewGetObject(dep *proc.Dependencies, meta payload.Meta, passed bool) *GetObject {
	return &GetObject{
		dep:    dep,
		meta:   meta,
		passed: passed,
	}
}

func (s *GetObject) Present(ctx context.Context, f flow.Flow) error {
	pl, err := payload.Unmarshal(s.meta.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal GetObject message")
	}
	msg, ok := pl.(*payload.GetObject)
	if !ok {
		return fmt.Errorf("wrong request type: %T", pl)
	}

	ctx, _ = inslogger.WithField(ctx, "object", msg.ObjectID.DebugString())

	passIfNotExecutor := !s.passed
	jet := proc.NewFetchJet(msg.ObjectID, flow.Pulse(ctx), s.meta, passIfNotExecutor)
	s.dep.FetchJet(jet)
	if err := f.Procedure(ctx, jet, false); err != nil {
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

	// To ensure, that we have the index. Because index can be on a heavy node.
	// If we don't have it and heavy does, GetObject fails because it should update light's index state.
	ensureIdx := proc.NewEnsureIndex(msg.ObjectID, objJetID, s.meta, flow.Pulse(ctx))
	s.dep.EnsureIndex(ensureIdx)
	if err := f.Procedure(ctx, ensureIdx, false); err != nil {
		return err
	}

	send := proc.NewSendObject(s.meta, msg.ObjectID, msg.RequestID)
	s.dep.SendObject(send)
	return f.Procedure(ctx, send, false)
}
