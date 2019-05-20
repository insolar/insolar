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

package handle

import (
	"context"

	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/flow/bus"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/ledger/light/proc"
)

type SetRecord struct {
	dep        *proc.Dependencies
	msg        *message.SetRecord
	busMessage bus.Message
}

func NewSetRecord(dep *proc.Dependencies, busMessage bus.Message, msg *message.SetRecord) *SetRecord {
	return &SetRecord{
		dep:        dep,
		msg:        msg,
		busMessage: busMessage,
	}
}

func (s *SetRecord) Present(ctx context.Context, f flow.Flow) error {
	jet := proc.NewFetchJet(*s.msg.TargetRef.Record(), flow.Pulse(ctx), s.busMessage.WatermillMsg)
	s.dep.FetchJet(jet)
	if err := f.Procedure(ctx, jet, true); err != nil {
		return err
	}
	hot := proc.NewWaitHot(jet.Result.Jet, flow.Pulse(ctx), s.busMessage.WatermillMsg)
	s.dep.WaitHot(hot)
	if err := f.Procedure(ctx, hot, true); err != nil {
		return err
	}

	setRecord := proc.NewSetRecord(jet.Result.Jet, s.busMessage.ReplyTo, s.msg.Record)
	s.dep.SetRecord(setRecord)
	return f.Procedure(ctx, setRecord, false)
}
