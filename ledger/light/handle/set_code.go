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

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/insolar/insolar/insolar/record"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/ledger/light/proc"
)

type SetCode struct {
	dep     *proc.Dependencies
	message *message.Message
	passed  bool
}

func NewSetCode(dep *proc.Dependencies, msg *message.Message, passed bool) *SetCode {
	return &SetCode{
		dep:     dep,
		message: msg,
		passed:  passed,
	}
}

func (s *SetCode) Present(ctx context.Context, f flow.Flow) error {
	pl, err := payload.UnmarshalFromMeta(s.message.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal payload")
	}
	msg, ok := pl.(*payload.SetCode)
	if !ok {
		return fmt.Errorf("unexpected payload type: %T", pl)
	}

	calc := proc.NewCalculateID(msg.Record, flow.Pulse(ctx))
	s.dep.CalculateID(calc)
	if err := f.Procedure(ctx, calc, true); err != nil {
		return err
	}
	recID := calc.Result.ID

	passIfNotExecutor := !s.passed
	jet := proc.NewCheckJet(recID, flow.Pulse(ctx), s.message, passIfNotExecutor)
	s.dep.CheckJet(jet)
	if err := f.Procedure(ctx, jet, true); err != nil {
		return err
	}

	rec := record.Code{}
	err = rec.Unmarshal(msg.Record)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal record")
	}

	setCode := proc.NewSetCode(s.message, rec, msg.Code, recID, jet.Result.Jet)
	s.dep.SetCode(setCode)
	return f.Procedure(ctx, setCode, false)
}
