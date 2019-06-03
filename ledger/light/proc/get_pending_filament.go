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

	"github.com/insolar/insolar/insolar/flow/bus"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/ledger/object"
	"github.com/pkg/errors"
)

type GetPendingFilament struct {
	msg      *message.GetPendingFilament
	repylyTo chan<- bus.Reply

	Dep struct {
		PendingAccessor object.PendingAccessor
	}
}

func NewGetPendingFilament(msg *message.GetPendingFilament, rep chan<- bus.Reply) *GetPendingFilament {
	return &GetPendingFilament{
		msg:      msg,
		repylyTo: rep,
	}
}

func (p *GetPendingFilament) Proceed(ctx context.Context) error {
	meta, err := p.Dep.PendingAccessor.MetaForObjID(ctx, p.msg.PN, p.msg.ObjectID)
	if err != nil {
		return errors.Wrap(err, "can't fetch a pendings meta")
	}
	records, err := p.Dep.PendingAccessor.Records(ctx, p.msg.PN, p.msg.ObjectID)
	if err != nil {
		return errors.Wrap(err, "can't fetch pendings")
	}

	p.repylyTo <- bus.Reply{
		Reply: &reply.PendingFilament{
			ID:                p.msg.ObjectID,
			Records:           records,
			HasFullChain:      meta.IsStateCalculated,
			PreviousPendingPN: meta.PrevSegmentPN,
		},
	}

	return nil
}
