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

	"github.com/insolar/insolar/insolar/message"

	"github.com/insolar/insolar/ledger/light/proc"

	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/flow/bus"
)

type RegisterChild struct {
	dep     *proc.Dependencies
	replyTo chan<- bus.Reply
	message *message.RegisterChild
}

func NewRegisterChild(dep *proc.Dependencies, rep chan<- bus.Reply, msg *message.RegisterChild) *RegisterChild {
	return &RegisterChild{
		dep:     dep,
		replyTo: rep,
		message: msg,
	}
}

func (s *RegisterChild) Present(ctx context.Context, f flow.Flow) error {
	jet := proc.NewFetchJet(*s.message.DefaultTarget().Record() /* TODO is it right? */, flow.Pulse(ctx), s.replyTo)
	s.dep.FetchJet(jet)
	if err := f.Procedure(ctx, jet, false); err != nil {
		return err
	}

	code := proc.NewRegisterChild(s.message, s.replyTo)
	s.dep.RegisterChild(code) // TODO: figure out what is that for
	return f.Procedure(ctx, code, false)
}
