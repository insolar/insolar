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

	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/object"
)

type RegisterChild struct {
	jet     insolar.JetID
	msg     *message.RegisterChild
	pulse   insolar.PulseNumber
	idx     object.Lifeline
	message payload.Meta

	Dep struct {
		IDLocker              object.IDLocker
		LifelineIndex         object.LifelineIndex
		JetCoordinator        jet.Coordinator
		RecordModifier        object.RecordModifier
		LifelineStateModifier object.LifelineStateModifier
		PCS                   insolar.PlatformCryptographyScheme
		Sender                bus.Sender
	}
}

func NewRegisterChild(jet insolar.JetID, msg *message.RegisterChild, pulse insolar.PulseNumber, index object.Lifeline, message payload.Meta) *RegisterChild {
	return &RegisterChild{
		jet:     jet,
		msg:     msg,
		pulse:   pulse,
		idx:     index,
		message: message,
	}
}

func (p *RegisterChild) Proceed(ctx context.Context) error {
	err := p.process(ctx)
	if err != nil {
		msg := bus.ErrorAsMessage(ctx, err)
		p.Dep.Sender.Reply(ctx, p.message, msg)
	}
	return err
}

func (p *RegisterChild) process(ctx context.Context) error {
	virtRec := record.Virtual{}
	err := virtRec.Unmarshal(p.msg.Record)
	if err != nil {
		return errors.Wrap(err, "can't deserialize record")
	}
	concreteRec := record.Unwrap(&virtRec)
	childRec, ok := concreteRec.(*record.Child)
	if !ok {
		return errors.New("wrong child record")
	}

	p.Dep.IDLocker.Lock(p.msg.Parent.Record())
	defer p.Dep.IDLocker.Unlock(p.msg.Parent.Record())

	hash := record.HashVirtual(p.Dep.PCS.ReferenceHasher(), virtRec)
	recID := insolar.NewID(p.pulse, hash)

	// Children exist and pointer does not match (preserving chain consistency).
	// For the case when vm can't save or send result to another vm and it tries to update the same record again
	if p.idx.ChildPointer != nil && !childRec.PrevChild.Equal(*p.idx.ChildPointer) && p.idx.ChildPointer != recID {
		return errors.New("invalid child record")
	}

	hash = record.HashVirtual(p.Dep.PCS.ReferenceHasher(), virtRec)
	child := insolar.NewID(p.pulse, hash)
	rec := record.Material{
		Virtual: &virtRec,
		JetID:   p.jet,
	}

	err = p.Dep.RecordModifier.Set(ctx, *child, rec)

	if err == object.ErrOverride {
		inslogger.FromContext(ctx).WithField("type", fmt.Sprintf("%T", virtRec)).Warn("set record override (#2)")
		child = recID
	} else if err != nil {
		return errors.Wrap(err, "can't save record into storage")
	}

	p.idx.ChildPointer = child
	if p.msg.AsType != nil {
		p.idx.SetDelegate(*p.msg.AsType, p.msg.Child)
	}
	p.idx.LatestUpdate = p.pulse
	p.idx.JetID = p.jet

	err = p.Dep.LifelineIndex.Set(ctx, p.pulse, *p.msg.Parent.Record(), p.idx)
	if err != nil {
		return err
	}
	err = p.Dep.LifelineStateModifier.SetLifelineUsage(ctx, p.pulse, *p.msg.Parent.Record())
	if err != nil {
		return errors.Wrap(err, "can't update a lifeline status")
	}

	msg := bus.ReplyAsMessage(ctx, &reply.ID{ID: *child})
	p.Dep.Sender.Reply(ctx, p.message, msg)
	return nil
}
