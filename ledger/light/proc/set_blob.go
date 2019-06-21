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

	watermillMsg "github.com/ThreeDotsLabs/watermill/message"
	"github.com/insolar/insolar/insolar"
	wmBus "github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/flow/bus"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/ledger/blob"
	"github.com/insolar/insolar/ledger/light/hot"
	"github.com/insolar/insolar/ledger/object"
	"github.com/pkg/errors"
)

type SetBlob struct {
	message payload.Meta
	msg     *message.SetBlob
	jet     insolar.JetID

	Dep struct {
		BlobAccessor  blob.Accessor
		BlobModifier  blob.Modifier
		PCS           insolar.PlatformCryptographyScheme
		WriteAccessor hot.WriteAccessor
		Sender        wmBus.Sender
	}
}

func NewSetBlob(jetID insolar.JetID, message payload.Meta, msg *message.SetBlob) *SetBlob {
	return &SetBlob{
		msg:     msg,
		message: message,
		jet:     jetID,
	}
}

func (p *SetBlob) Proceed(ctx context.Context) error {
	r := p.reply(ctx)
	var msg *watermillMsg.Message
	if r.Err != nil {
		msg = wmBus.ErrorAsMessage(ctx, r.Err)
	} else {
		msg = wmBus.ReplyAsMessage(ctx, r.Reply)
	}
	p.Dep.Sender.Reply(ctx, p.message, msg)
	return nil
}

func (p *SetBlob) reply(ctx context.Context) bus.Reply {
	done, err := p.Dep.WriteAccessor.Begin(ctx, flow.Pulse(ctx))
	if err == hot.ErrWriteClosed {
		return bus.Reply{Err: flow.ErrCancelled}
	}
	if err != nil {
		return bus.Reply{Err: errors.Wrap(err, "failed to start write")}
	}
	defer done()
	msg := p.msg

	calculatedID := object.CalculateIDForBlob(p.Dep.PCS, flow.Pulse(ctx), msg.Memory)

	_, err = p.Dep.BlobAccessor.ForID(ctx, *calculatedID)
	if err == nil {
		return bus.Reply{Reply: &reply.ID{ID: *calculatedID}}
	}
	if err != blob.ErrNotFound {
		return bus.Reply{Err: err}
	}

	err = p.Dep.BlobModifier.Set(ctx, *calculatedID, blob.Blob{Value: msg.Memory, JetID: p.jet})
	if err == nil {
		return bus.Reply{Reply: &reply.ID{ID: *calculatedID}}
	}
	if err == blob.ErrOverride {
		return bus.Reply{Reply: &reply.ID{ID: *calculatedID}}
	}

	return bus.Reply{Err: err}
}
