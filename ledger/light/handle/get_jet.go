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
	"github.com/insolar/insolar/ledger/light/proc"
)

type GetJet struct {
	dep    *proc.Dependencies
	meta   payload.Meta
	passed bool
}

func NewGetJet(dep *proc.Dependencies, meta payload.Meta, passed bool) *GetJet {
	return &GetJet{
		dep:    dep,
		meta:   meta,
		passed: passed,
	}
}

func (h *GetJet) Present(ctx context.Context, f flow.Flow) error {
	pl, err := payload.Unmarshal(h.meta.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal GetJet message")
	}
	msg, ok := pl.(*payload.GetJet)
	if !ok {
		return fmt.Errorf("wrong request type: %T", pl)
	}

	getJet := proc.NewGetJet(h.meta, msg.ObjectID, msg.PulseNumber)
	h.dep.GetJet(getJet)
	return f.Procedure(ctx, getJet, false)
}
