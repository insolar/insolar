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

type GetFilament struct {
	dep *proc.Dependencies

	meta payload.Meta
}

func NewGetFilament(dep *proc.Dependencies, meta payload.Meta) *GetFilament {
	return &GetFilament{
		dep:  dep,
		meta: meta,
	}
}

func (s *GetFilament) Present(ctx context.Context, f flow.Flow) error {
	pl, err := payload.Unmarshal(s.meta.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal GetFilament message")
	}
	msg, ok := pl.(*payload.GetFilament)
	if !ok {
		return fmt.Errorf("wrong request type: %T", pl)
	}

	getFilament := proc.NewSendFilament(s.meta, msg.ObjectID, msg.StartFrom, msg.ReadUntil)
	s.dep.SendFilament(getFilament)
	return f.Procedure(ctx, getFilament, false)
}
