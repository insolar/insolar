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

type GetPulse struct {
	dep  *proc.Dependencies
	meta payload.Meta
}

func NewGetPulse(dep *proc.Dependencies, meta payload.Meta) *GetPulse {
	return &GetPulse{
		dep:  dep,
		meta: meta,
	}
}

func (s *GetPulse) Present(ctx context.Context, f flow.Flow) error {
	pl, err := payload.Unmarshal(s.meta.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal GetPulse message")
	}
	msg, ok := pl.(*payload.GetPulse)
	if !ok {
		return fmt.Errorf("wrong request type: %T", pl)
	}

	getPulse := proc.NewGetPulse(s.meta, msg.PulseNumber)
	s.dep.GetPulse(getPulse)
	return f.Procedure(ctx, getPulse, false)
}
