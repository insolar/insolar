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

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow/bus"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/reply"
)

type GetJet struct {
	msg     *message.GetJet
	replyTo chan<- bus.Reply

	Dep struct {
		Jets jet.Storage
	}
}

func NewGetJet(msg *message.GetJet, rep chan<- bus.Reply) *GetJet {
	return &GetJet{
		msg:     msg,
		replyTo: rep,
	}
}

func (p *GetJet) Proceed(ctx context.Context) error {
	jetID, actual := p.Dep.Jets.ForID(ctx, p.msg.Pulse, p.msg.Object)
	p.replyTo <- bus.Reply{Reply: &reply.Jet{ID: insolar.ID(jetID), Actual: actual}, Err: nil}
	return nil
}
