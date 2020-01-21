// Copyright 2020 Insolar Network Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
