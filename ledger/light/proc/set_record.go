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

	"github.com/insolar/insolar/ledger/light/executor"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/flow/bus"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/ledger/light/hot"
	"github.com/insolar/insolar/ledger/object"
)

type SetRecord struct {
	replyTo chan<- bus.Reply
	record  []byte
	jet     insolar.JetID

	Dep struct {
		Bus insolar.MessageBus

		Coordinator jet.Coordinator

		PCS                  insolar.PlatformCryptographyScheme
		RecordModifier       object.RecordModifier
		WriteAccessor        hot.WriteAccessor
		PendingRequestsLimit int
		Filaments            executor.FilamentModifier
	}
}

func NewSetRecord(jetID insolar.JetID, replyTo chan<- bus.Reply, record []byte) *SetRecord {
	return &SetRecord{
		record:  record,
		replyTo: replyTo,
		jet:     jetID,
	}
}

func (p *SetRecord) Proceed(ctx context.Context) error {
	p.replyTo <- p.reply(ctx)
	return nil
}

func (p *SetRecord) reply(ctx context.Context) bus.Reply {
	done, err := p.Dep.WriteAccessor.Begin(ctx, flow.Pulse(ctx))
	if err == hot.ErrWriteClosed {
		return bus.Reply{Err: flow.ErrCancelled}
	}
	if err != nil {
		return bus.Reply{Err: errors.Wrap(err, "failed to start write")}
	}
	defer done()

	virtual := record.Virtual{}
	err = virtual.Unmarshal(p.record)
	if err != nil {
		return bus.Reply{Err: errors.Wrap(err, "can't deserialize record")}
	}

	hash := record.HashVirtual(p.Dep.PCS.ReferenceHasher(), virtual)
	id := insolar.NewID(flow.Pulse(ctx), hash)

	result, ok := record.Unwrap(&virtual).(*record.Result)
	if ok {
		err := p.Dep.Filaments.SetResult(ctx, *id, p.jet, *result)
		if err != nil {
			return bus.Reply{Err: errors.Wrap(err, "failed to save result")}
		}
	} else {
		err := p.Dep.RecordModifier.Set(ctx, *id, record.Material{Virtual: &virtual, JetID: p.jet})
		if err != nil {
			return bus.Reply{Err: errors.Wrap(err, "failed to save record")}
		}
	}

	return bus.Reply{Reply: &reply.ID{ID: *id}}
}
