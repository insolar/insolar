// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package handle

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar/payload"

	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/ledger/light/proc"
)

type PassState struct {
	dep  *proc.Dependencies
	meta payload.Meta
}

func NewPassState(dep *proc.Dependencies, meta payload.Meta) *PassState {
	return &PassState{
		dep:  dep,
		meta: meta,
	}
}

func (s *PassState) Present(ctx context.Context, f flow.Flow) error {
	// Pass state unmarshal pl
	pl, err := payload.Unmarshal(s.meta.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal payload")
	}
	passState, ok := pl.(*payload.PassState)
	if !ok {
		return fmt.Errorf("unexpected payload type %T", pl)
	}

	// Origin message unmarshal
	pl, err = payload.Unmarshal(passState.Origin)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal origin payload")
	}
	origin, ok := pl.(*payload.Meta)
	if !ok {
		return fmt.Errorf("unexpected payload type %T", pl)
	}

	// Origin message unmarshal pl
	pl, err = payload.Unmarshal(origin.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal payload")
	}
	_, ok = pl.(*payload.GetObject)
	if !ok {
		return fmt.Errorf("unexpected payload type %T", pl)
	}

	state := proc.NewPassState(s.meta, passState.StateID, *origin)
	s.dep.PassState(state)
	return f.Procedure(ctx, state, false)
}
