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
	idx     object.Lifeline
	replyTo chan<- bus.Reply

	Dep struct {
		IDLocker                   object.IDLocker
		IndexStorage               object.IndexStorage
		JetCoordinator             jet.Coordinator
		RecordModifier             object.RecordModifier
		IndexStateModifier         object.ExtendedIndexModifier
		PlatformCryptographyScheme insolar.PlatformCryptographyScheme
	}
}

func NewRegisterChild(jet insolar.JetID, msg *message.RegisterChild, pulse insolar.PulseNumber, index object.Lifeline, replyTo chan<- bus.Reply) *RegisterChild {
	return &RegisterChild{
		jet:     jet,
		msg:     msg,
		pulse:   pulse,
		idx:     index,
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
	r, err := object.DecodeVirtual(p.msg.Record)
	if err != nil {
		return errors.Wrap(err, "can't deserialize record")
	}
	childRec, ok := r.(*object.ChildRecord)
	if !ok {
		return errors.New("wrong child record")
	}

	p.Dep.IDLocker.Lock(p.msg.Parent.Record())
	defer p.Dep.IDLocker.Unlock(p.msg.Parent.Record())
	p.Dep.IndexStateModifier.SetUsageForPulse(ctx, *p.msg.Parent.Record(), p.pulse)
	recID := object.NewRecordIDFromRecord(p.Dep.PlatformCryptographyScheme, p.pulse, childRec)

	// Children exist and pointer does not match (preserving chain consistency).
	// For the case when vm can't save or send result to another vm and it tries to update the same record again
	if p.idx.ChildPointer != nil && !childRec.PrevChild.Equal(*p.idx.ChildPointer) && p.idx.ChildPointer != recID {
		return errors.New("invalid child record")
	}

	child := object.NewRecordIDFromRecord(p.Dep.PlatformCryptographyScheme, p.pulse, childRec)
	rec := record.MaterialRecord{
		Record: childRec,
		JetID:  insolar.JetID(p.jet),
	}

	err = p.Dep.RecordModifier.Set(ctx, *child, rec)

	if err == object.ErrOverride {
		inslogger.FromContext(ctx).WithField("type", fmt.Sprintf("%T", r)).Warn("set record override (#2)")
		child = recID
	} else if err != nil {
		return errors.Wrap(err, "can't save record into storage")
	}

	p.idx.ChildPointer = child
	if p.msg.AsType != nil {
		p.idx.Delegates[*p.msg.AsType] = p.msg.Child
	}
	p.idx.LatestUpdate = p.pulse
	p.idx.JetID = insolar.JetID(p.jet)
	err = p.Dep.IndexStorage.Set(ctx, *p.msg.Parent.Record(), p.idx)
	if err != nil {
		return err
	}

	p.replyTo <- bus.Reply{Reply: &reply.ID{ID: *child}}
	return nil
}
