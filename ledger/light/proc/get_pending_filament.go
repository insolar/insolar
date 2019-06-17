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

	watermillMsg "github.com/ThreeDotsLabs/watermill/message"

	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/ledger/object"
	"github.com/pkg/errors"
)

type GetPendingFilament struct {
	msg     *message.GetPendingFilament
	message *watermillMsg.Message
	Dep     struct {
		PendingAccessor  object.PendingAccessor
		LifelineAccessor object.LifelineAccessor
		Sender           bus.Sender
	}
}

func NewGetPendingFilament(msg *message.GetPendingFilament, wmmessage *watermillMsg.Message) *GetPendingFilament {
	return &GetPendingFilament{
		msg:     msg,
		message: wmmessage,
	}
}

func (p *GetPendingFilament) Proceed(ctx context.Context) error {
	ctx, span := instracer.StartSpan(ctx, fmt.Sprintf("GetPendingFilament"))
	defer span.End()

	records, err := p.Dep.PendingAccessor.Records(ctx, p.msg.PN, p.msg.ObjectID)
	if err != nil {
		return errors.Wrap(err, "[GetPendingFilament] can't fetch pendings")
	}
	msg := bus.ReplyAsMessage(ctx, &reply.PendingFilament{
		ObjID:   p.msg.ObjectID,
		Records: records,
	})
	p.Dep.Sender.Reply(ctx, p.message, msg)
	return nil
}
