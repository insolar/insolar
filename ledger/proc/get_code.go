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
	"github.com/insolar/insolar/insolar/flow/bus"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/ledger/storage/blob"
	"github.com/insolar/insolar/ledger/storage/object"
	"github.com/pkg/errors"
)

type GetCode struct {
	JetID   insolar.JetID
	Message bus.Message
	Code    insolar.Reference

	Result struct {
		CodeRec *object.CodeRecord
	}

	Dep struct {
		Bus                    insolar.MessageBus
		DelegationTokenFactory insolar.DelegationTokenFactory
		RecordAccessor         object.RecordAccessor
		Coordinator            insolar.JetCoordinator
		Accessor               blob.Accessor
		BlobModifier           blob.Modifier
	}
}

func (p *GetCode) Proceed(ctx context.Context) error {
	// ctx = contextWithJet(ctx, insolar.ID(p.JetID))
	r := bus.Reply{}
	r.Reply, r.Err = p.handle(ctx)
	p.Message.ReplyTo <- r
	return nil
}

func (p *GetCode) handle(ctx context.Context) (insolar.Reply, error) {
	parcel := p.Message.Parcel
	rec, err := p.Dep.RecordAccessor.ForID(ctx, *p.Code.Record())
	if err == object.ErrNotFound {
		// We don't have code record. Must be on another node.
		node, err := p.Dep.Coordinator.NodeForJet(ctx, insolar.ID(p.JetID), parcel.Pulse(), p.Code.Record().Pulse())
		if err != nil {
			return nil, err
		}
		return reply.NewGetCodeRedirect(p.Dep.DelegationTokenFactory, parcel, node)
	}
	if err != nil {
		return nil, err
	}

	virtRec := rec.Record
	codeRec, ok := virtRec.(*object.CodeRecord)
	if !ok {
		return nil, errors.Wrap(ErrInvalidRef, "failed to retrieve code record")
	}

	code, err := p.Dep.Accessor.ForID(ctx, *codeRec.Code)
	if err == blob.ErrNotFound {
		hNode, err := p.Dep.Coordinator.Heavy(ctx, parcel.Pulse())
		if err != nil {
			return nil, err
		}
		return p.saveCodeFromHeavy(ctx, p.JetID, p.Code, *codeRec.Code, hNode)
	}

	if err != nil {
		return nil, err
	}

	rep := reply.Code{
		Code:        code.Value,
		MachineType: codeRec.MachineType,
	}

	return &rep, nil
}

func (p *GetCode) saveCodeFromHeavy(
	ctx context.Context, jetID insolar.JetID, code insolar.Reference, blobID insolar.ID, heavy *insolar.Reference,
) (*reply.Code, error) {
	genericReply, err := p.Dep.Bus.Send(ctx, &message.GetCode{
		Code: code,
	}, &insolar.MessageSendOptions{
		Receiver: heavy,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to send")
	}
	rep, ok := genericReply.(*reply.Code)
	if !ok {
		return nil, fmt.Errorf("failed to fetch code: unexpected reply type %T (reply=%+v)", genericReply, genericReply)
	}

	err = p.Dep.BlobModifier.Set(ctx, blobID, blob.Blob{JetID: jetID, Value: rep.Code})
	if err != nil {
		return nil, errors.Wrap(err, "failed to save")
	}
	return rep, nil
}
