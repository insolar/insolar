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

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/flow/bus"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/ledger/proc"
	"github.com/pkg/errors"
)

type GetCode struct {
	dep *proc.Dependencies

	Message bus.Message
}

func (s *GetCode) Present(ctx context.Context, f flow.Flow) error {
	msg := s.Message.Parcel.Message().(*message.GetCode)

	jet := s.dep.FetchJet(&proc.FetchJet{Parcel: s.Message.Parcel})
	if err := f.Procedure(ctx, jet); err != nil {
		if err == flow.ErrCancelled {
			return err
		}
		return err
	}

	if jet.Result.Miss {
		rep := &proc.ReturnReply{
			ReplyTo: s.Message.ReplyTo,
			Reply:   &reply.JetMiss{JetID: insolar.ID(jet.Result.Jet)},
		}
		if err := f.Procedure(ctx, rep); err != nil {
			return err
		}
		return errors.New("jet miss")
	}

	codeRec := s.dep.GetCode(&proc.GetCode{
		JetID:   jet.Result.Jet,
		Message: s.Message,
		Code:    msg.Code,
	})
	return f.Procedure(ctx, codeRec)
}
