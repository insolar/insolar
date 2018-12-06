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
	"bytes"
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
	if len(cb.Requests) < 1 {
		return 0, errors.New("casebind is empty")
	}

	es := lr.UpsertExecution(ref)
	es.insContext = ctx
	es.validate = true
	es.objectbody = nil
	var cbr core.CaseBindReplay
	err := func() error {
		lr.caseBindReplaysMutex.Lock()
		defer lr.caseBindReplaysMutex.Unlock()
		if _, ok := lr.caseBindReplays[ref]; ok {
			return errors.New("already validating this ref")
		}
		lr.caseBindReplays[ref] = core.CaseBindReplay{
			Pulse:    p,
			CaseBind: cb,
			Request:  0,
			Record:   -1,
			Steps:    0,
		}
		cbr = lr.caseBindReplays[ref]
		return nil
	}()
	if err != nil {
		return 0, err
	}

	defer func() {
		lr.caseBindReplaysMutex.Lock()
		defer lr.caseBindReplaysMutex.Unlock()
		delete(lr.caseBindReplays, ref)
	}()

	for {
		start, step := lr.nextValidationStep(ref)
		if step < 0 {
			return step, errors.New("no validation data")
		} else if start == nil { // finish
			return step, nil
		}
		if start.Type != core.CaseRecordTypeStart {
			return step, errors.New("step between two shores")
		}

		msg := start.Resp.(core.Message)
		parcel, err := lr.ParcelFactory.Create(
			ctx, msg, ref, nil,
		)
		if err != nil {
			return 0, errors.New("failed to create a parcel message")
		}

		traceStep, step := lr.nextValidationStep(ref)
		if traceStep == nil {
			return step, errors.New("trace is missing")
		}

		traceID, ok := traceStep.Resp.(string)
		if !ok {
			return step, errors.New("trace is wrong type")
		}

		es.insContext = inslogger.ContextWithTrace(es.insContext, traceID)
		es.Lock()
		ret, err := lr.executeOrValidate(es.insContext, es, ValidationChecker{lr: lr, cb: cbr}, parcel)
		if err != nil {
			return 0, errors.Wrap(err, "validation step failed")
		}
		stop, step := lr.nextValidationStep(ref)
		if step < 0 {
			return 0, errors.New("validation container broken")
		} else if stop.Type != core.CaseRecordTypeResult {
			return step, errors.New("Validation stoped not on result")
		}

		switch need := stop.Resp.(type) {
		case *reply.CallMethod:
			if got, ok := ret.(*reply.CallMethod); !ok {
				return step, errors.New("not result type callmethod")
			} else if !bytes.Equal(got.Data, need.Data) {
				return step, errors.New("body mismatch")
			} else if !bytes.Equal(got.Result, need.Result) {
				return step, errors.New("result mismatch")
			}
		case *reply.CallConstructor:
			if got, ok := ret.(*reply.CallConstructor); !ok {
				return step, errors.New("not result type callconstructor")
			} else if !got.Object.Equal(*need.Object) {
				return step, errors.New("constructed refs mismatch mismatch")
			}
		default:
			return step, errors.New("unknown result type")
		}
	}
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
	_, err = lr.MessageBus.Send(
		ctx,
		&message.ValidationResults{
			RecordRef:        msg.GetReference(),
			PassedStepsCount: passedStepsCount,
			Error:            errstr,
		},
		nil,
	)

	return &reply.OK{}, err
}

func (lr *LogicRunner) ProcessValidationResults(ctx context.Context, inmsg core.Parcel) (core.Reply, error) {
	msg, ok := inmsg.Message().(*message.ValidationResults)
	if !ok {
		return nil, errors.Errorf("ProcessValidationResults got argument typed %t", inmsg)
	}
	c, _ := lr.GetConsensus(ctx, msg.RecordRef)
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
	c, _ := lr.GetConsensus(ctx, msg.RecordRef)
	c.AddExecutor(ctx, inmsg, msg)
	return &reply.OK{}, nil
}

// ValidationBehaviour is a special object that responsible for validation behavior of other methods.
type ValidationBehaviour interface {
	Begin(refs Ref, record core.CaseRecord)
	End(refs Ref, record core.CaseRecord)
	GetRole() core.DynamicRole
	ModifyContext(ctx *core.LogicCallContext)
	NeedSave() bool
	RegisterRequest(p core.Parcel) (*Ref, error)
}

type ValidationSaver struct {
	lr *LogicRunner
}

func (vb ValidationSaver) RegisterRequest(p core.Parcel) (*Ref, error) {
	ctx := context.TODO()
	m := p.Message().(message.IBaseLogicMessage)
	reqid, err := vb.lr.ArtifactManager.RegisterRequest(ctx, p)
	if err != nil {
		return nil, err
	}
	// TODO: use proper conversion
	reqref := Ref{}
	reqref.SetRecord(*reqid)

	vb.lr.addObjectCaseRecord(m.GetReference(), core.CaseRecord{
		Type:   core.CaseRecordTypeRequest,
		ReqSig: HashInterface(vb.lr.PlatformCryptographyScheme, m),
		Resp:   reqref,
	})
	return &reqref, err
}

func (vb ValidationSaver) NeedSave() bool {
	return true
}

func (vb ValidationSaver) ModifyContext(ctx *core.LogicCallContext) {
	// nothing need
}

func (vb ValidationSaver) Begin(refs Ref, record core.CaseRecord) {
	vb.lr.addObjectCaseRecord(refs, record)
}

func (vb ValidationSaver) End(refs Ref, record core.CaseRecord) {
	vb.lr.addObjectCaseRecord(refs, record)
}

func (vb ValidationSaver) GetRole() core.DynamicRole {
	return core.DynamicRoleVirtualExecutor
}

type ValidationChecker struct {
	lr *LogicRunner
	cb core.CaseBindReplay
}

func (vb ValidationChecker) RegisterRequest(p core.Parcel) (*Ref, error) {
	m := p.Message().(message.IBaseLogicMessage)
	cr, _ := vb.lr.nextValidationStep(m.GetReference())
	if core.CaseRecordTypeRequest != cr.Type {
		return nil, errors.New("Wrong validation type on Request")
	}
	if !bytes.Equal(cr.ReqSig, HashInterface(vb.lr.PlatformCryptographyScheme, m)) {
		return nil, errors.New("Wrong validation sig on Request")
	}
	if req, ok := cr.Resp.(Ref); ok {
		return &req, nil
	}
	return nil, errors.Errorf("wrong validation, request contains %t", cr.Resp)

}

func (vb ValidationChecker) NeedSave() bool {
	return false
}

func (vb ValidationChecker) ModifyContext(ctx *core.LogicCallContext) {
	ctx.Pulse = vb.cb.Pulse
}

func (vb ValidationChecker) Begin(refs Ref, record core.CaseRecord) {
	// do nothing, everything done in lr.Validate
}

func (vb ValidationChecker) End(refs Ref, record core.CaseRecord) {
	// do nothing, everything done in lr.Validate
}

func (vb ValidationChecker) GetRole() core.DynamicRole {
	return core.DynamicRoleVirtualValidator
}
