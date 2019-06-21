/*
 *    Copyright 2019 Insolar Technologies
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package proc

import (
	"context"
	"fmt"

	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/ledger/object"
	"github.com/pkg/errors"
)

type GetPendingFilament struct {
	meta payload.Meta

	Dep struct {
		PendingAccessor object.HeavyPendingAccessor
		Sender          bus.Sender
	}
}

func NewGetPendingFilament(meta payload.Meta) *GetPendingFilament {
	return &GetPendingFilament{
		meta: meta,
	}
}

func (p *GetPendingFilament) Proceed(ctx context.Context) error {
	ctx, span := instracer.StartSpan(ctx, fmt.Sprintf("GetPendingFilament"))
	defer span.End()

	getPFil := payload.GetPendingFilament{}
	err := getPFil.Unmarshal(p.meta.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to decode PassState payload")
	}

	inslogger.FromContext(ctx).Debugf("GetPendingFilament objID == %v, startFrom %v, ReadUntil %v", getPFil.ObjectID.DebugString(), getPFil.StartFrom, getPFil.ReadUntil)
	records, err := p.Dep.PendingAccessor.Records(ctx, getPFil.StartFrom, getPFil.ReadUntil, getPFil.ObjectID)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("[GetPendingFilament] can't fetch pendings, pn - %v,  %v", getPFil.ObjectID.DebugString(), getPFil.StartFrom))
	}

	inslogger.FromContext(ctx).Debugf("GetPendingFilament objID == %v, records - %v, records == nil - %v, , startFrom %v, ReadUntil %v", getPFil.ObjectID.DebugString(), len(records), records == nil, getPFil.StartFrom, getPFil.ReadUntil)
	msg, err := payload.NewMessage(&payload.PendingFilament{
		ObjectID: getPFil.ObjectID,
		Records:  records,
	})
	if err != nil {
		return errors.Wrap(err, "failed to create a PendingFilament message")
	}
	go p.Dep.Sender.Reply(ctx, p.meta, msg)
	return nil
}
