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

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/ledger/light/proc"
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
	pl, err := payload.Unmarshal(s.message.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal Activate message")
	}
	msg, ok := pl.(*payload.Activate)
	if !ok {
		return fmt.Errorf("wrong request type: %T", pl)
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
		return errors.New("request is empty")
	}

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

	passIfNotExecutor := !s.passed
	jet := proc.NewFetchJet(*activate.Request.GetLocal(), flow.Pulse(ctx), s.message, passIfNotExecutor)
	s.dep.FetchJet(jet)
	if err := f.Procedure(ctx, jet, true); err != nil {
		if err == proc.ErrNotExecutor && passIfNotExecutor {
			return nil
		}
		return err
	}
	objJetID := jet.Result.Jet

	hot := proc.NewWaitHot(objJetID, flow.Pulse(ctx), s.message)
	s.dep.WaitHot(hot)
	if err := f.Procedure(ctx, hot, false); err != nil {
		return err
	}

	setResult := proc.NewSetResult(s.message, objJetID, *result, activate)
	s.dep.SetResult(setResult)
	return f.Procedure(ctx, setResult, false)
}
