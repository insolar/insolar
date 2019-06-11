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

type ExpirePending struct {
	replyTo   chan<- bus.Reply
	expiredPN insolar.PulseNumber
	objID     insolar.ID

	Dep struct {
		Coordinator           jet.Coordinator
		MessageBus            insolar.MessageBus
		FilamentStateModifier object.PendingFilamentStateModifier
	}
}

func NewExpirePending(replyTo chan<- bus.Reply, expiredPN insolar.PulseNumber, objID insolar.ID) *ExpirePending {
	return &ExpirePending{
		replyTo:   replyTo,
		expiredPN: expiredPN,
		objID:     objID,
	}
}

func (p *ExpirePending) Proceed(ctx context.Context) error {
	err := p.process(ctx)
	if err != nil {
		p.replyTo <- bus.Reply{Err: err}
	}
	return err
}

func (p *ExpirePending) process(ctx context.Context) error {
	logger := inslogger.FromContext(ctx)

	isBeyond, err := p.Dep.Coordinator.IsBeyondLimit(ctx, flow.Pulse(ctx), p.expiredPN)
	if err != nil {
		logger.Error(errors.Wrapf(err, "[ExpirePending] failed to save index - %v", p.objID))
		return err
	}
	if !isBeyond {
		return nil
	}

	heavy, err := p.Dep.Coordinator.Heavy(ctx, p.expiredPN)
	if err != nil {
		logger.Errorf("expireRequests failed with: %v", err)
	}
	genericReact, err := p.Dep.MessageBus.Send(ctx,
		&message.GetOpenRequests{ObjID: p.objID, PN: p.expiredPN},
		&insolar.MessageSendOptions{
			Receiver: heavy,
		})
	if err != nil {
		logger.Errorf("fetching expireRequests failed with: %v", err)
	}
	switch rep := genericReact.(type) {
	case *reply.OpenRequestsOnHeavy:
		return p.Dep.FilamentStateModifier.ExpireRequests(ctx, flow.Pulse(ctx), p.objID, rep.Requests)
	case *reply.Error:
		logger.Errorf("expireRequests failed with: %v", rep.Error())
		return rep.Error()
	default:
		return fmt.Errorf("expireRequests failed with unexpected reply: %v", rep)
	}
}
