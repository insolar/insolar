// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package proc

import (
	"context"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/light/executor"
)

type GetPendings struct {
	message  payload.Meta
	objectID insolar.ID
	count    int
	skip     []insolar.ID

	dep struct {
		filaments executor.FilamentCalculator
		sender    bus.Sender
	}
}

func NewGetPendings(msg payload.Meta, objectID insolar.ID, count int, skip []insolar.ID) *GetPendings {
	return &GetPendings{
		message:  msg,
		objectID: objectID,
		count:    count,
		skip:     skip,
	}
}

func (p *GetPendings) Dep(
	f executor.FilamentCalculator,
	s bus.Sender,
) {
	p.dep.filaments = f
	p.dep.sender = s
}

func (p *GetPendings) Proceed(ctx context.Context) error {
	pendings, err := p.dep.filaments.OpenedRequests(ctx, flow.Pulse(ctx), p.objectID, true)
	if err != nil {
		return errors.Wrap(err, "failed to calculate pending")
	}
	logger := inslogger.FromContext(ctx)
	if len(pendings) == 0 {
		errMsg, errErr := payload.NewMessage(&payload.Error{
			Text: insolar.ErrNoPendingRequest.Error(),
			Code: payload.CodeNoPendings,
		})
		if errErr != nil {
			logger.Error("Failed to return error reply: ", errErr.Error())
			return errErr
		}
		p.dep.sender.Reply(ctx, p.message, errMsg)
		return nil
	}

	var skipMap map[insolar.ID]struct{}
	if len(p.skip) > 0 {
		skipMap = make(map[insolar.ID]struct{}, len(p.skip))
		for _, id := range p.skip {
			skipMap[id] = struct{}{}
		}
	}

	var ids []insolar.ID
	for _, pend := range pendings {
		if len(ids) >= p.count {
			break
		}
		if skipMap != nil {
			if _, ok := skipMap[pend.RecordID]; ok {
				continue
			}
		}
		ids = append(ids, pend.RecordID)
	}

	if len(ids) == 0 {
		errMsg, errErr := payload.NewMessage(&payload.Error{
			Text: insolar.ErrNoPendingRequest.Error(),
			Code: payload.CodeNoPendings,
		})
		if errErr != nil {
			logger.Error("Failed to return error reply: ", errErr.Error())
			return errErr
		}
		p.dep.sender.Reply(ctx, p.message, errMsg)
		return nil
	}

	msg, err := payload.NewMessage(&payload.IDs{
		IDs: ids,
	})
	if err != nil {
		return errors.Wrap(err, "failed to create reply")
	}

	p.dep.sender.Reply(ctx, p.message, msg)
	return nil
}
