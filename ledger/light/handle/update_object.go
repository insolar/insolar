///
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
///

package handle

import (
	"context"

	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/flow/bus"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/ledger/light/proc"
)

type UpdateObject struct {
	dep *proc.Dependencies

	Message bus.Message
	replyTo chan<- bus.Reply
}

func (s *UpdateObject) Present(ctx context.Context, f flow.Flow) error {
	msg := s.Message.Parcel.Message().(*message.UpdateObject)

	jet := proc.NewFetchJet(*msg.Object.Record(), flow.Pulse(ctx), s.replyTo)
	s.dep.FetchJet(jet)
	if err := f.Procedure(ctx, jet, false); err != nil {
		return err
	}

	/*jet := &WaitJet{
		dep:     s.dep,
		Message: s.Message,
	}
	if err := f.Handle(ctx, jet.Present); err != nil {
		return err
	}*/

	updateProc := &proc.UpdateObject{
		JetID:   jet.Result.Jet,
		BusMsg:  s.Message,
		Message: msg,
		Parcel:  s.Message.Parcel,
	}
	s.dep.UpdateObject(updateProc)

	if err := f.Procedure(ctx, updateProc, false); err != nil {
		return err
	}
	// return updateProc.Proceed(ctx)

	return nil
}
