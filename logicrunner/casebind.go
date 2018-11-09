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
	"github.com/ugorji/go/codec"
	"golang.org/x/crypto/sha3"
)

func HashInterface(in interface{}) []byte {
	var s []byte
	ch := new(codec.CborHandle)
	err := codec.NewEncoderBytes(&s, ch).Encode(in)
	if err != nil {
		panic("Can't marshal: " + err.Error())
	}
	sh := sha3.New224()
	return sh.Sum(s)
}

func (lr *LogicRunner) Validate(ref Ref, p core.Pulse, cr []core.CaseRecord) (int, error) {
	if len(cr) < 1 {
		return 0, errors.New("casebind is empty")
	}

	es := lr.UpsertExecution(ref)
	ctx := context.TODO()
	es.insContext = ctx
	es.validate = true
	es.objectbody = nil
	err := func() error {
		lr.caseBindReplaysMutex.Lock()
		defer lr.caseBindReplaysMutex.Unlock()
		if _, ok := lr.caseBindReplays[ref]; ok {
			return errors.New("already validating this ref")
		}
		lr.caseBindReplays[ref] = core.CaseBindReplay{
			Pulse:      p,
			Records:    cr,
			RecordsLen: len(cr),
			Step:       0,
		}
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
		parcel, err := message.NewParcel(
			ctx, msg, ref, lr.Network.GetPrivateKey(), lr.execution[ref].callContext.Pulse.PulseNumber, nil,
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
		ret, err := lr.Execute(es.insContext, parcel)
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
	passedStepsCount, validationError := lr.Validate(msg.GetReference(), msg.GetPulse(), msg.GetCaseRecords())
	_, err := lr.MessageBus.Send(
		ctx,
		&message.ValidationResults{
			RecordRef:        msg.GetReference(),
			PassedStepsCount: passedStepsCount,
			Error:            validationError.Error(),
		},
	)

	return nil, err
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
	return nil, nil
}

func (lr *LogicRunner) ExecutorResults(ctx context.Context, inmsg core.Parcel) (core.Reply, error) {
	msg, ok := inmsg.Message().(*message.ExecutorResults)
	if !ok {
		return nil, errors.Errorf("ProcessValidationResults got argument typed %t", inmsg)
	}
	c, _ := lr.GetConsensus(ctx, msg.RecordRef)
	c.AddExecutor(ctx, inmsg, msg)
	return nil, nil
}

// ValidationBehaviour is a special object that responsible for validation behavior of other methods.
type ValidationBehaviour interface {
	Begin(refs Ref, record core.CaseRecord)
	End(refs Ref, record core.CaseRecord)
	GetRole() core.JetRole
	ModifyContext(ctx *core.LogicCallContext)
	NeedSave() bool
	RegisterRequest(m message.IBaseLogicMessage) (*Ref, error)
}

type ValidationSaver struct {
	lr *LogicRunner
}

func (vb ValidationSaver) RegisterRequest(m message.IBaseLogicMessage) (*Ref, error) {
	ctx := context.TODO()
	reqid, err := vb.lr.ArtifactManager.RegisterRequest(ctx, m)
	if err != nil {
		return nil, err
	}
	// TODO: use proper conversion
	reqref := Ref{}
	reqref.SetRecord(*reqid)

	vb.lr.addObjectCaseRecord(m.GetReference(), core.CaseRecord{
		Type:   core.CaseRecordTypeRequest,
		ReqSig: HashInterface(m),
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

func (vb ValidationSaver) GetRole() core.JetRole {
	return core.RoleVirtualExecutor
}

type ValidationChecker struct {
	lr *LogicRunner
	cb core.CaseBindReplay
}

func (vb ValidationChecker) RegisterRequest(m message.IBaseLogicMessage) (*Ref, error) {
	cr, _ := vb.lr.nextValidationStep(m.GetReference())
	if core.CaseRecordTypeRequest != cr.Type {
		return nil, errors.New("Wrong validation type on Request")
	}
	if !bytes.Equal(cr.ReqSig, HashInterface(m)) {
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

func (vb ValidationChecker) GetRole() core.JetRole {
	return core.RoleVirtualValidator
}
