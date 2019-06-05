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

	"github.com/insolar/insolar/insolar/flow/bus"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/ledger/object"
	"github.com/pkg/errors"
)

type GetPendingFilament struct {
	msg      *message.GetPendingFilament
	repylyTo chan<- bus.Reply

	Dep struct {
		PendingAccessor  object.PendingAccessor
		LifelineAccessor object.LifelineAccessor
	}
}

func NewGetPendingFilament(msg *message.GetPendingFilament, rep chan<- bus.Reply) *GetPendingFilament {
	return &GetPendingFilament{
		msg:      msg,
		repylyTo: rep,
	}
}

func (p *GetPendingFilament) Proceed(ctx context.Context) error {
	ctx, span := instracer.StartSpan(ctx, fmt.Sprintf("GetPendingFilament"))
	defer span.End()

	isStateCalc, err := p.Dep.PendingAccessor.IsStateCalculated(ctx, p.msg.PN, p.msg.ObjectID)
	if err != nil {
		return errors.Wrap(err, "[GetPendingFilament] can't fetch a pendings meta")
	}
	records, err := p.Dep.PendingAccessor.Records(ctx, p.msg.PN, p.msg.ObjectID)
	if err != nil {
		return errors.Wrap(err, "[GetPendingFilament] can't fetch pendings")
	}
	lfl, err := p.Dep.LifelineAccessor.ForID(ctx, p.msg.PN, p.msg.ObjectID)
	if err != nil {
		return errors.Wrap(err, "[GetPendingFilament] can't fetch lifeline")
	}

	p.repylyTo <- bus.Reply{
		Reply: &reply.PendingFilament{
			ObjID:             p.msg.ObjectID,
			Records:           records,
			HasFullChain:      isStateCalc,
			PreviousPendingPN: lfl.PreviousPendingFilament,
		},
	}

	return nil
}
