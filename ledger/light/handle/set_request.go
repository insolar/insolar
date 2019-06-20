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
	"fmt"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
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

	virtual := record.Virtual{}
	err = virtual.Unmarshal(msg.Request)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal Request record")
	}

	rec := record.Unwrap(&virtual)
	request, ok := rec.(*record.Request)
	if !ok {
		return fmt.Errorf("wrong request type: %T", rec)
	}
	// This is a workaround. VM should not register such requests.
	// TODO: remove after INS-1939
	if request.CallType != record.CTMethod {
		inslogger.FromContext(ctx).Error("request is not registered")
		return s.setActivationRequest(ctx, reqID, virtual, f)
	}

	if request.Object == nil {
		return errors.New("object is nil")
	}

	passIfNotExecutor := !s.passed
	jet := proc.NewCheckJet(*request.Object.Record(), flow.Pulse(ctx), s.message, passIfNotExecutor)
	s.dep.CheckJet(jet)
	if err := f.Procedure(ctx, jet, true); err != nil {
		if err == proc.ErrNotExecutor && passIfNotExecutor {
			return nil
		}
		return err
	}
	objJetID := jet.Result.Jet

	hot := proc.NewWaitHotWM(objJetID, flow.Pulse(ctx), s.message)
	s.dep.WaitHotWM(hot)
	if err := f.Procedure(ctx, hot, false); err != nil {
		return err
	}

	setRequest := proc.NewSetRequest(s.message, virtual, reqID, objJetID)
	s.dep.SetRequest(setRequest)
	return f.Procedure(ctx, setRequest, false)
}

func (s *SetRequest) setActivationRequest(ctx context.Context, reqID insolar.ID, request record.Virtual, f flow.Flow) error {
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

	setActivationRequest := proc.NewSetActivationRequest(s.message, request, reqID, reqJetID)

	s.dep.SetActivationRequest(setActivationRequest)
	return f.Procedure(ctx, setActivationRequest, false)
}
