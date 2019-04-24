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
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/flow/bus"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/object"
	"github.com/pkg/errors"
)

type GetIndex struct {
	object  insolar.Reference
	jet     insolar.JetID
	replyTo chan<- bus.Reply

	Result struct {
		Index object.Lifeline
	}

	Dep struct {
		IndexState  object.ExtendedIndexModifier
		Locker      object.IDLocker
		Storage     object.IndexStorage
		Coordinator jet.Coordinator
		Bus         insolar.MessageBus
	}
}

func NewGetIndex(obj insolar.Reference, jetID insolar.JetID, rep chan<- bus.Reply) *GetIndex {
	return &GetIndex{
		object:  obj,
		jet:     jetID,
		replyTo: rep,
	}
}

func (p *GetIndex) Proceed(ctx context.Context) error {
	err := p.process(ctx)
	if err != nil {
		p.replyTo <- bus.Reply{Err: err}
	}
	return err
}

func (p *GetIndex) process(ctx context.Context) error {
	objectID := *p.object.Record()
	logger := inslogger.FromContext(ctx)

	p.Dep.Locker.Lock(&objectID)
	defer p.Dep.Locker.Unlock(&objectID)

	idx, err := p.Dep.Storage.ForID(ctx, objectID)
	if err == nil {
		p.Result.Index = idx
		return nil
	}
	if err != object.ErrIndexNotFound {
		return errors.Wrap(err, "failed to fetch index")
	}
	p.Dep.IndexState.SetUsageForPulse(ctx, objectID, flow.Pulse(ctx))

	logger.Debug("failed to fetch index (fetching from heavy)")
	heavy, err := p.Dep.Coordinator.Heavy(ctx, flow.Pulse(ctx))
	if err != nil {
		return errors.Wrap(err, "failed to calculate heavy")
	}
	genericReply, err := p.Dep.Bus.Send(ctx, &message.GetObjectIndex{
		Object: p.object,
	}, &insolar.MessageSendOptions{
		Receiver: heavy,
	})
	if err != nil {
		return errors.Wrap(err, "failed to fetch index from heavy")
	}
	rep, ok := genericReply.(*reply.ObjectIndex)
	if !ok {
		return fmt.Errorf("failed to fetch index from heavy: unexpected reply type %T", genericReply)
	}

	p.Result.Index, err = object.DecodeIndex(rep.Index)
	if err != nil {
		return errors.Wrap(err, "failed to decode index")
	}

	p.Result.Index.JetID = p.jet
	err = p.Dep.Storage.Set(ctx, objectID, p.Result.Index)
	if err != nil {
		return errors.Wrap(err, "failed to save index")
	}

	return nil
}
