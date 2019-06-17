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

	watermillMsg "github.com/ThreeDotsLabs/watermill/message"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/ledger/light/proc"
)

type GetJet struct {
	dep          *proc.Dependencies
	msg          *message.GetJet
	watermillMsg *watermillMsg.Message
}

func NewGetJet(dep *proc.Dependencies, watermillMsg *watermillMsg.Message, msg *message.GetJet) *GetJet {
	return &GetJet{
		dep:          dep,
		msg:          msg,
		watermillMsg: watermillMsg,
	}
}

func (s *GetJet) Present(ctx context.Context, f flow.Flow) error {
	getJet := proc.NewGetJet(s.msg, s.watermillMsg)
	s.dep.GetJet(getJet)
	return f.Procedure(ctx, getJet, false)
}
