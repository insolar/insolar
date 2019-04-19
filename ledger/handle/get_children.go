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

	"github.com/insolar/insolar/insolar"

	"github.com/insolar/insolar/ledger/proc"

	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/flow/bus"
)

// GetChildren Handler
type GetChildren struct {
	dep *proc.Dependencies

	Message bus.Message
}

func (s *GetChildren) Present(ctx context.Context, f flow.Flow) error {
	jet := &WaitJet{
		dep:     s.dep,
		Message: s.Message,
	}
	if err := f.Handle(ctx, jet.Present); err != nil {
		return err
	}

	p := s.dep.GetChildren(&proc.GetChildren{
		Jet:     insolar.ID(jet.Res.Jet),
		Message: s.Message,
	})
	if err := f.Procedure(ctx, p); err != nil {
		return err
	}

	rep := &proc.ReturnReply{
		ReplyTo: s.Message.ReplyTo,
		Reply:   p.Result.Reply,
	}

	return f.Procedure(ctx, rep)
}
