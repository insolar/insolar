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
	"encoding/gob"
	"reflect"

	"github.com/insolar/insolar/instrumentation/inslogger"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/pkg/errors"
)

type CaseRequest struct {
	Message    core.Message
	Request    core.RecordRef
	MessageBus core.MessageBus
	Reply      core.Reply
	Error      error
}

// CaseBinder is a whole result of executor efforts on every object it seen on this pulse
type CaseBind struct {
	Requests []CaseRequest
}

func NewCaseBind() *CaseBind {
	return &CaseBind{Requests: make([]CaseRequest, 0)}
}

func NewCaseBindFromValidateMessage(ctx context.Context, mb core.MessageBus, msg *message.ValidateCaseBind) *CaseBind {
	res := &CaseBind{
		Requests: make([]CaseRequest, len(msg.Requests)),
	}
	for i, req := range msg.Requests {
		mb, err := mb.NewPlayer(ctx, bytes.NewReader(req.MessageBusTape))
		if err != nil {
			panic("couldn't read tape: " + err.Error())
		}
		res.Requests[i] = CaseRequest{
			Message:    req.Message,
			Request:    req.Request,
			MessageBus: mb,
			Reply:      req.Reply,
			Error:      req.Error,
		}
	}
	return res
}

func NewCaseBindFromExecutorResultsMessage(msg *message.ExecutorResults) *CaseBind {
	panic("not implemented")
}

func (cb *CaseBind) getCaseBindForMessage(ctx context.Context) []message.CaseBindRequest {
	if cb == nil {
		return make([]message.CaseBindRequest, 0)
	}

	requests := make([]message.CaseBindRequest, len(cb.Requests))

	for i, req := range cb.Requests {
		var tape bytes.Buffer
		err := req.MessageBus.WriteTape(ctx, &tape)
		if err != nil {
			panic("couldn't write tape: " + err.Error())
		}
		requests[i] = message.CaseBindRequest{
			Message:        req.Message,
			Request:        req.Request,
			MessageBusTape: tape.Bytes(),
			Reply:          req.Reply,
			Error:          req.Error,
		}
	}

	return requests
}

func (cb *CaseBind) ToValidateMessage(ctx context.Context, ref Ref, pulse core.Pulse) *message.ValidateCaseBind {
	res := &message.ValidateCaseBind{
		RecordRef: ref,
		Requests:  cb.getCaseBindForMessage(ctx),
		Pulse:     pulse,
	}
	return res
}

func (cb *CaseBind) ToExecutorResultsMessage(ctx context.Context, ref Ref, pending bool, queue []message.ExecutionQueueElement) *message.ExecutorResults {
	res := &message.ExecutorResults{
		RecordRef: ref,
		Pending:   pending,
		Requests:  cb.getCaseBindForMessage(ctx),
		Queue:     queue,
	}

	return res
}

func (cb *CaseBind) NewRequest(msg core.Message, request Ref, mb core.MessageBus) *CaseRequest {
	res := CaseRequest{
		Message:    msg,
		Request:    request,
		MessageBus: mb,
	}
	cb.Requests = append(cb.Requests, res)
	return &cb.Requests[len(cb.Requests)-1]
}

type CaseBindReplay struct {
	Pulse    core.Pulse
	CaseBind CaseBind
	Request  int
	Record   int
	Steps    int
	Fail     int
}

func NewCaseBindReplay(cb CaseBind) *CaseBindReplay {
	return &CaseBindReplay{
		CaseBind: cb,
		Request:  -1,
		Record:   -1,
	}
}

func (r *CaseBindReplay) NextRequest() *CaseRequest {
	if r.Request+1 >= len(r.CaseBind.Requests) {
		return nil
	}
	r.Request++
	return &r.CaseBind.Requests[r.Request]
}

func HashInterface(scheme core.PlatformCryptographyScheme, in interface{}) []byte {
	s, err := core.Serialize(in)
	if err != nil {
		panic("Can't marshal: " + err.Error())
	}
	return scheme.IntegrityHasher().Hash(s)
}

func (lr *LogicRunner) Validate(ctx context.Context, ref Ref, p core.Pulse, cb CaseBind) (int, error) {
	os := lr.UpsertObjectState(ref)
	vs := os.StartValidation()
	vs.ArtifactManager = lr.ArtifactManager

	vs.Lock()
	defer vs.Unlock()

	checker := &ValidationChecker{
		lr: lr,
		cb: NewCaseBindReplay(cb),
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
		ctx = core.ContextWithMessageBus(ctx, request.MessageBus)

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

	passedStepsCount, validationError := lr.Validate(
		ctx, msg.GetReference(), msg.GetPulse(), *NewCaseBindFromValidateMessage(ctx, lr.MessageBus, msg),
	)
	errstr := ""
	if validationError != nil {
		errstr = validationError.Error()
	}

	currentSlotPulse, err := lr.PulseStorage.Current(ctx)
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

	// now we have 2 different types of data in message.ExecutorResults
	// one part of it is about consensus
	// another one is about prepare state on new executor after pulse
	// TODO make it in different goroutines

	// prepare state after previous executor
	err := lr.prepareObjectState(ctx, msg)
	if err != nil {
		return &reply.Error{}, err
	}

	// validation things
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
	caseBind *CaseBind
	current  *CaseRequest
}

func (vb *ValidationSaver) Mode() string {
	return "execution"
}

func (vb *ValidationSaver) NewRequest(msg core.Message, request Ref, mb core.MessageBus) {
	vb.current = vb.caseBind.NewRequest(msg, request, mb)
}

func (vb *ValidationSaver) Result(reply core.Reply, err error) error {
	if vb.current == nil {
		return errors.New("result call without request registered")
	}
	vb.current.Reply = reply
	vb.current.Error = err
	return nil
}

type ValidationChecker struct {
	lr      *LogicRunner
	cb      *CaseBindReplay
	current *CaseRequest
}

func (vb *ValidationChecker) Mode() string {
	return "validation"
}

func (vb *ValidationChecker) NextRequest() *CaseRequest {
	vb.current = vb.cb.NextRequest()
	return vb.current
}

func (vb *ValidationChecker) Result(reply core.Reply, err error) error {
	if vb.current == nil {
		return errors.New("result call without request registered")
	}
	// TODO: reflect.DeepEqual is not what we want to go with, we should
	// go with HASH comparision
	if !reflect.DeepEqual(vb.current.Reply, reply) {
		return errors.Errorf("replies arn't equal: expected: %+v, got: %+v, err: %+v", vb.current.Reply, reply, err)
	}
	if !reflect.DeepEqual(vb.current.Error, err) {
		return errors.New("errors arn't equal")
	}
	return nil
}

func init() {
	gob.Register(&CaseRequest{})
	gob.Register(&CaseBind{})
}
