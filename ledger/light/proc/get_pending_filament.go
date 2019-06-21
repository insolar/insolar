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
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/ledger/object"
	"github.com/pkg/errors"
)

type GetPendingFilament struct {
	message   payload.Meta
	objID     insolar.ID
	startFrom insolar.PulseNumber
	readUntil insolar.PulseNumber

	Dep struct {
		PendingAccessor object.PendingAccessor
		Sender          bus.Sender
	}
}

func NewGetPendingFilament(msg payload.Meta, objID insolar.ID, startFrom insolar.PulseNumber, readUntil insolar.PulseNumber) *GetPendingFilament {
	return &GetPendingFilament{
		message:   msg,
		objID:     objID,
		startFrom: startFrom,
		readUntil: readUntil,
	}
}

func (p *GetPendingFilament) Proceed(ctx context.Context) error {
	ctx, span := instracer.StartSpan(ctx, fmt.Sprintf("GetPendingFilament"))
	defer span.End()

	inslogger.FromContext(ctx).Debugf("GetPendingFilament objID: %v, startFrom: %v, readUntil: %v", p.objID.DebugString(), p.startFrom, p.readUntil)
	records, err := p.Dep.PendingAccessor.Records(ctx, p.startFrom, p.objID)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("[GetPendingFilament] can't fetch pendings, pn - %v,  %v", p.objID.DebugString(), p.startFrom))
	}

	inslogger.FromContext(ctx).Debugf("GetPendingFilament objID == %v, len of records - %v", p.objID.DebugString(), len(records))
	for _, r := range records {
		inslogger.FromContext(ctx).Debugf("GetPendingFilament records - %v", r.RecordID.DebugString())
	}

	if len(records) == 0 {
		panic("why this happened?")
	}

	inslogger.FromContext(ctx).Debugf("GetPendingFilament objID == %v, records - %v", p.objID.DebugString(), len(records))
	msg, err := payload.NewMessage(&payload.PendingFilament{
		ObjectID: p.objID,
		Records:  records,
	})
	if err != nil {
		return errors.Wrap(err, "failed to create a PendingFilament message")
	}
	go p.Dep.Sender.Reply(ctx, p.message, msg)
	return nil
}
