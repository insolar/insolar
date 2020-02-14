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

type GetCode struct {
	dep    *proc.Dependencies
	meta   payload.Meta
	passed bool
}

func NewGetCode(dep *proc.Dependencies, meta payload.Meta, passed bool) *GetCode {
	return &GetCode{
		dep:    dep,
		meta:   meta,
		passed: passed,
	}
}

func (s *GetCode) Present(ctx context.Context, f flow.Flow) error {
	pl, err := payload.Unmarshal(s.meta.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal GetCode message")
	}
	msg, ok := pl.(*payload.GetCode)
	if !ok {
		return fmt.Errorf("wrong request type: %T", pl)
	}

	passIfNotFound := !s.passed
	code := proc.NewGetCode(s.meta, msg.CodeID, passIfNotFound)
	s.dep.GetCode(code)
	return f.Procedure(ctx, code, false)
}
