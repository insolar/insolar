/*
 *
 *  *    Copyright 2018 Insolar
 *  *
 *  *    Licensed under the Apache License, Version 2.0 (the "License");
 *  *    you may not use this file except in compliance with the License.
 *  *    You may obtain a copy of the License at
 *  *
 *  *        http://www.apache.org/licenses/LICENSE-2.0
 *  *
 *  *    Unless required by applicable law or agreed to in writing, software
 *  *    distributed under the License is distributed on an "AS IS" BASIS,
 *  *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  *    See the License for the specific language governing permissions and
 *  *    limitations under the License.
 *
 */

package logicrunner

import (
	"context"

	"github.com/insolar/insolar/instrumentation/inslogger"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/pkg/errors"
)

func HashInterface(scheme core.PlatformCryptographyScheme, in interface{}) []byte {
	s, err := core.Serialize(in)
	if err != nil {
		panic("Can't marshal: " + err.Error())
	}
	return scheme.IntegrityHasher().Hash(s)
}

func (lr *LogicRunner) Validate(ctx context.Context, ref Ref, p core.Pulse, cb core.CaseBind) (int, error) {
	os := lr.UpsertObjectState(ref)
	vs := os.StartValidation()

	vs.Lock()
	defer vs.Unlock()

	checker := &ValidationChecker{
		lr: lr,
		cb: core.NewCaseBindReplay(cb),
	}
	vs.Behaviour = checker

	for {
		request := checker.NextRequest()
		if request == nil {
			break
		}

		msg := request.Message
		parcel, err := lr.ParcelFactory.Create(ctx, msg, ref, nil, *core.GenesisPulse)
		if err != nil {
			return 0, errors.New("failed to create a parcel")
		}

		traceID := "TODO" // FIXME

		ctx = inslogger.ContextWithTrace(ctx, traceID)

		vs.Current = &CurrentExecution{
			Context: ctx,
			Request: &request.Request,
		}

		reply, err := lr.executeOrValidate(ctx, vs, parcel)

		err = vs.Behaviour.Result(reply, err)
		if err != nil {
			return 0, errors.Wrap(err, "validation step failed")
		}
	}
	return 1, nil
}

func (lr *LogicRunner) ValidateCaseBind(ctx context.Context, inmsg core.Parcel) (core.Reply, error) {
	msg, ok := inmsg.Message().(*message.ValidateCaseBind)
	if !ok {
		return nil, errors.New("Execute( ! message.ValidateCaseBindInterface )")
	}

	err := lr.CheckOurRole(ctx, msg, core.DynamicRoleVirtualValidator)
	if err != nil {
		return nil, errors.Wrap(err, "can't play role")
	}

	passedStepsCount, validationError := lr.Validate(ctx, msg.GetReference(), msg.GetPulse(), msg.CaseBind)
	errstr := ""
	if validationError != nil {
		errstr = validationError.Error()
	}

	currentSlotPulse, err := lr.PulseManager.Current(ctx)
	if err != nil {
		return nil, err
	}
	_, err = lr.MessageBus.Send(
		ctx,
		&message.ValidationResults{
			RecordRef:        msg.GetReference(),
			PassedStepsCount: passedStepsCount,
			Error:            errstr,
		},
		*currentSlotPulse,
		nil,
	)

	return &reply.OK{}, err
}

func (lr *LogicRunner) ProcessValidationResults(ctx context.Context, inmsg core.Parcel) (core.Reply, error) {
	msg, ok := inmsg.Message().(*message.ValidationResults)
	if !ok {
		return nil, errors.Errorf("ProcessValidationResults got argument typed %t", inmsg)
	}

	c := lr.GetConsensus(ctx, msg.RecordRef)
	if err := c.AddValidated(ctx, inmsg, msg); err != nil {
		return nil, err
	}
	return &reply.OK{}, nil
}

func (lr *LogicRunner) ExecutorResults(ctx context.Context, inmsg core.Parcel) (core.Reply, error) {
	msg, ok := inmsg.Message().(*message.ExecutorResults)
	if !ok {
		return nil, errors.Errorf("ProcessValidationResults got argument typed %t", inmsg)
	}
	c := lr.GetConsensus(ctx, msg.RecordRef)
	c.AddExecutor(ctx, inmsg, msg)
	return &reply.OK{}, nil
}

// ValidationBehaviour is a special object that responsible for validation behavior of other methods.
type ValidationBehaviour interface {
	Mode() string
	Result(reply core.Reply, err error) error
}

type ValidationSaver struct {
	lr       *LogicRunner
	caseBind *core.CaseBind
	current  *core.CaseRequest
}

func (vb ValidationSaver) Mode() string {
	return "execution"
}

func (vb ValidationSaver) NewRequest(msg core.Message, request Ref, mb core.MessageBus) {
	vb.current = vb.caseBind.NewRequest(msg, request, mb)
}

func (vb ValidationSaver) Result(reply core.Reply, err error) error {
	vb.current.Reply = reply
	vb.current.Error = err
	return nil
}

type ValidationChecker struct {
	lr *LogicRunner
	cb *core.CaseBindReplay
}

func (vb ValidationChecker) Mode() string {
	return "validation"
}

func (vb ValidationChecker) NextRequest() *core.CaseRequest {
	return vb.cb.NextRequest()
}

func (vb ValidationChecker) Result(reply core.Reply, err error) error {
	return nil
}
