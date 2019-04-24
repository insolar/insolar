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
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/flow/bus"
	"github.com/insolar/insolar/ledger/light/proc"
)

type GetCode struct {
	dep     *proc.Dependencies
	code    insolar.Reference
	replyTo chan<- bus.Reply
}

func NewGetCode(dep *proc.Dependencies, rep chan<- bus.Reply, code insolar.Reference) *GetCode {
	return &GetCode{
		dep:     dep,
		code:    code,
		replyTo: rep,
	}
}

func (s *GetCode) Present(ctx context.Context, f flow.Flow) error {
	code := s.dep.GetCode(&proc.GetCode{
		ReplyTo: s.replyTo,
		Code:    s.code,
	})
	return code.Proceed(ctx)
}
