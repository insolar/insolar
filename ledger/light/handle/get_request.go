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

type GetRequest struct {
	dep     *proc.Dependencies
	message payload.Meta
	passed  bool
}

func NewGetRequest(dep *proc.Dependencies, msg payload.Meta, passed bool) *GetRequest {
	return &GetRequest{
		dep:     dep,
		message: msg,
		passed:  passed,
	}
}

func (s *GetRequest) Present(ctx context.Context, f flow.Flow) error {
	pl, err := payload.Unmarshal(s.message.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal GetRequest message")
	}
	msg, ok := pl.(*payload.GetRequest)
	if !ok {
		return fmt.Errorf("wrong request type: %T", pl)
	}

	passIfNotFound := !s.passed
	req := proc.NewGetRequest(s.message, msg.ObjectID, msg.RequestID, passIfNotFound)
	s.dep.GetRequest(req)
	return f.Procedure(ctx, req, false)
}
