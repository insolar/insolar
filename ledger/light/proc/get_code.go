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
	"fmt"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/flow/bus"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/ledger/blob"
	"github.com/insolar/insolar/ledger/object"
	"github.com/pkg/errors"
)

type GetCode struct {
	ReplyTo chan<- bus.Reply
	Code    insolar.Reference

	Dep struct {
		Bus            insolar.MessageBus
		RecordAccessor object.RecordAccessor
		Coordinator    jet.Coordinator
		BlobAccessor   blob.Accessor
		CheckJet       CheckJet
	}
}

func (p *GetCode) Proceed(ctx context.Context) error {
	p.ReplyTo <- p.reply(ctx)
	return nil
}

func (p *GetCode) reply(ctx context.Context) bus.Reply {
	codeID := *p.Code.Record()
	jetID, mine, err := p.Dep.CheckJet(ctx, codeID, codeID.Pulse())
	if err != nil {
		return bus.Reply{Err: errors.Wrap(err, "failed to check jet")}
	}
	if !mine {
		return bus.Reply{Reply: &reply.JetMiss{JetID: insolar.ID(jetID), Pulse: codeID.Pulse()}}
	}

	rec, err := p.Dep.RecordAccessor.ForID(ctx, codeID)
	if err == object.ErrNotFound {
		heavy, err := p.Dep.Coordinator.Heavy(ctx, flow.Pulse(ctx))
		if err != nil {
			return bus.Reply{Err: errors.Wrap(err, "failed to calculate heavy")}
		}
		genericReply, err := p.Dep.Bus.Send(ctx, &message.GetCode{
			Code: p.Code,
		}, &insolar.MessageSendOptions{
			Receiver: heavy,
		})
		if err != nil {
			return bus.Reply{Err: errors.Wrap(err, "failed to fetch code from heavy")}
		}
		rep, ok := genericReply.(*reply.Code)
		if !ok {
			err := fmt.Errorf(
				"failed to fetch code from heavy: unexpected reply type %T",
				genericReply,
			)
			return bus.Reply{Err: err}
		}
		return bus.Reply{Reply: rep}
	}
	if err != nil {
		return bus.Reply{Err: errors.Wrap(err, "failed to fetch code")}
	}

	codeRec, ok := rec.Record.(*object.CodeRecord)
	if !ok {
		return bus.Reply{Err: errors.Wrap(ErrInvalidRef, "failed to retrieve code record")}
	}

	code, err := p.Dep.BlobAccessor.ForID(ctx, *codeRec.Code)
	if err != nil {
		return bus.Reply{Err: errors.Wrap(err, "failed to fetch code blob")}
	}

	rep := &reply.Code{
		Code:        code.Value,
		MachineType: codeRec.MachineType,
	}
	return bus.Reply{Reply: rep}
}
