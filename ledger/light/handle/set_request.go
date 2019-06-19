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
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/ledger/light/proc"
	"github.com/pkg/errors"
)

type SetRequest struct {
	dep     *proc.Dependencies
	message payload.Meta
	passed  bool
}

func NewSetRequest(dep *proc.Dependencies, msg payload.Meta, passed bool) *SetRequest {
	return &SetRequest{
		dep:     dep,
		message: msg,
		passed:  passed,
	}
}

func (s *SetRequest) Present(ctx context.Context, f flow.Flow) error {
	msg := payload.SetRequest{}
	err := msg.Unmarshal(s.message.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal SetRequest message")
	}

	calc := proc.NewCalculateID(msg.Request, flow.Pulse(ctx))
	s.dep.CalculateID(calc)
	if err := f.Procedure(ctx, calc, true); err != nil {
		return err
	}
	reqID := calc.Result.ID

	passIfNotExecutor := !s.passed
	jet := proc.NewCheckJet(reqID, flow.Pulse(ctx), s.message, passIfNotExecutor)
	s.dep.CheckJet(jet)
	if err := f.Procedure(ctx, jet, true); err != nil {
		if err == proc.ErrNotExecutor && passIfNotExecutor {
			return nil
		}
		return err
	}
	reqJetID := jet.Result.Jet

	hot := proc.NewWaitHotWM(reqJetID, flow.Pulse(ctx), s.message)
	s.dep.WaitHotWM(hot)
	if err := f.Procedure(ctx, hot, false); err != nil {
		return err
	}

	request := record.Request{}
	err = request.Unmarshal(msg.Request)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal Request record")
	}

	setRequest := proc.NewSetRequest(s.message, request, reqID, reqJetID)
	s.dep.SetRequest(setRequest)
	return f.Procedure(ctx, setRequest, false)
}
