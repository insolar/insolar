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

	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/ledger/light/proc"
	"github.com/pkg/errors"
)

type ActivateObject struct {
	dep     *proc.Dependencies
	message payload.Meta
	passed  bool
}

func NewActivateObject(dep *proc.Dependencies, msg payload.Meta, passed bool) *ActivateObject {
	return &ActivateObject{
		dep:     dep,
		message: msg,
		passed:  passed,
	}
}

func (s *ActivateObject) Present(ctx context.Context, f flow.Flow) error {
	msg := payload.Activate{}
	err := msg.Unmarshal(s.message.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal Activate message")
	}

	activateVirt := record.Virtual{}
	err = activateVirt.Unmarshal(msg.Record)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal Activate.Record record")
	}

	act := record.Unwrap(&activateVirt)
	activate, ok := act.(*record.Activate)
	if !ok {
		return fmt.Errorf("wrong activate record type: %T", act)
	}

	if activate.Request.IsEmpty() {
		return errors.New("request is nil")
	}

	calcAct := proc.NewCalculateID(msg.Record, flow.Pulse(ctx))
	s.dep.CalculateID(calcAct)
	if err := f.Procedure(ctx, calcAct, true); err != nil {
		return err
	}
	activateID := calcAct.Result.ID

	resultVirt := record.Virtual{}
	err = resultVirt.Unmarshal(msg.Result)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal Activate.Result record")
	}

	res := record.Unwrap(&resultVirt)
	result, ok := res.(*record.Result)
	if !ok {
		return fmt.Errorf("wrong result record type: %T", res)
	}

	calcRes := proc.NewCalculateID(msg.Result, flow.Pulse(ctx))
	s.dep.CalculateID(calcRes)
	if err := f.Procedure(ctx, calcRes, true); err != nil {
		return err
	}
	resultID := calcRes.Result.ID

	passIfNotExecutor := !s.passed
	jet := proc.NewCheckJet(*activate.Request.Record(), flow.Pulse(ctx), s.message, passIfNotExecutor)
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

	activateObject := proc.NewActivateObject(s.message, *activate, activateID, *result, resultID, objJetID)
	s.dep.ActivateObject(activateObject)
	return f.Procedure(ctx, activateObject, false)
}
