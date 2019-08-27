//
// Copyright 2019 Insolar Technologies GmbH
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
//

package proc

import (
	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/ledger/light/executor"
	"github.com/pkg/errors"
)

type GetPendings struct {
	message  payload.Meta
	objectID insolar.ID

	dep struct {
		filaments executor.FilamentCalculator
		sender    bus.Sender
	}
}

func NewGetPendings(msg payload.Meta, objectID insolar.ID) *GetPendings {
	return &GetPendings{
		message:  msg,
		objectID: objectID,
	}
}

func (gp *GetPendings) Dep(
	f executor.FilamentCalculator,
	s bus.Sender,
) {
	gp.dep.filaments = f
	gp.dep.sender = s
}

func (gp *GetPendings) Proceed(ctx context.Context) error {
	pendings, err := gp.dep.filaments.OpenedRequests(ctx, flow.Pulse(ctx), gp.objectID, true)
	if err != nil {
		return errors.Wrap(err, "failed to calculate pending")
	}
	if len(pendings) == 0 {
		msg, err := payload.NewMessage(&payload.Error{
			Code: payload.CodeNoPendings,
			Text: insolar.ErrNoPendingRequest.Error(),
		})
		if err != nil {
			return errors.Wrap(err, "failed to create reply")
		}
		gp.dep.sender.Reply(ctx, gp.message, msg)
		return nil
	}

	ids := make([]insolar.ID, len(pendings))
	for i, pend := range pendings {
		ids[i] = pend.RecordID
	}

	msg, err := payload.NewMessage(&payload.IDs{
		IDs: ids,
	})
	if err != nil {
		return errors.Wrap(err, "failed to create reply")
	}

	gp.dep.sender.Reply(ctx, gp.message, msg)
	return nil
}
