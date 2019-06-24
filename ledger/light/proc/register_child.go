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

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow/bus"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/object"
)

type RegisterChild struct {
	jet     insolar.JetID
	msg     *message.RegisterChild
	pulse   insolar.PulseNumber
	replyTo chan<- bus.Reply

	Dep struct {
		IndexLocker   object.IndexLocker
		IndexModifier object.IndexModifier
		IndexAccessor object.IndexAccessor

		JetCoordinator jet.Coordinator
		RecordModifier object.RecordModifier
		PCS            insolar.PlatformCryptographyScheme
	}
}

func NewRegisterChild(jet insolar.JetID, msg *message.RegisterChild, pulse insolar.PulseNumber, replyTo chan<- bus.Reply) *RegisterChild {
	return &RegisterChild{
		jet:     jet,
		msg:     msg,
		pulse:   pulse,
		replyTo: replyTo,
	}
}

func (p *RegisterChild) Proceed(ctx context.Context) error {
	err := p.process(ctx)
	if err != nil {
		p.replyTo <- bus.Reply{Err: err}
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

	p.Dep.IndexLocker.Lock(p.msg.Parent.Record())
	defer p.Dep.IndexLocker.Unlock(p.msg.Parent.Record())

	idx, err := p.Dep.IndexAccessor.ForID(ctx, p.pulse, *p.msg.Parent.Record())
	if err != nil {
		return err
	}

	hash := record.HashVirtual(p.Dep.PCS.ReferenceHasher(), virtRec)
	recID := insolar.NewID(p.pulse, hash)

	// Children exist and pointer does not match (preserving chain consistency).
	// For the case when vm can't save or send result to another vm and it tries to update the same record again
	if idx.Lifeline.ChildPointer != nil && !childRec.PrevChild.Equal(*idx.Lifeline.ChildPointer) && idx.Lifeline.ChildPointer != recID {
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

	idx.Lifeline.ChildPointer = child
	if p.msg.AsType != nil {
		idx.Lifeline.SetDelegate(*p.msg.AsType, p.msg.Child)
	}
	idx.Lifeline.LatestUpdate = p.pulse
	idx.Lifeline.JetID = p.jet
	idx.LifelineLastUsed = p.pulse

	err = p.Dep.IndexModifier.SetIndex(ctx, p.pulse, idx)
	if err != nil {
		return err
	}

	p.replyTo <- bus.Reply{Reply: &reply.ID{ID: *child}}
	return nil
}
