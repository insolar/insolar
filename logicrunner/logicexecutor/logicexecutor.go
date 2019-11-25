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

package logicexecutor

import (
	"bytes"
	"context"

	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/logicrunner/common"
	"github.com/insolar/insolar/logicrunner/machinesmanager"
	"github.com/insolar/insolar/logicrunner/requestresult"
)

//go:generate minimock -i github.com/insolar/insolar/logicrunner/logicexecutor.LogicExecutor -o ./ -s _mock.go -g
type LogicExecutor interface {
	Execute(ctx context.Context, transcript *common.Transcript) (artifacts.RequestResult, error)
	ExecuteMethod(ctx context.Context, transcript *common.Transcript) (artifacts.RequestResult, error)
	ExecuteConstructor(ctx context.Context, transcript *common.Transcript) (artifacts.RequestResult, error)
}

type logicExecutor struct {
	MachinesManager  machinesmanager.MachinesManager `inject:""`
	DescriptorsCache artifacts.DescriptorsCache      `inject:""`
	PulseAccessor    pulse.Accessor
}

func NewLogicExecutor(pulseAccessor pulse.Accessor) LogicExecutor {
	return &logicExecutor{PulseAccessor: pulseAccessor}
}

func (le *logicExecutor) Execute(ctx context.Context, transcript *common.Transcript) (artifacts.RequestResult, error) {
	ctx, _ = inslogger.WithField(ctx, "name", transcript.Request.Method)

	switch transcript.Request.CallType {
	case record.CTMethod:
		return le.ExecuteMethod(ctx, transcript)
	case record.CTSaveAsChild:
		return le.ExecuteConstructor(ctx, transcript)
	default:
		return nil, errors.New("Unknown request call type")
	}
}

func (le *logicExecutor) ExecuteMethod(ctx context.Context, transcript *common.Transcript) (artifacts.RequestResult, error) {
	ctx, span := instracer.StartSpan(ctx, "logicExecutor.ExecuteMethod")
	defer span.Finish()

	inslogger.FromContext(ctx).Debug("Executing method")

	request := transcript.Request

	objDesc := transcript.ObjectDescriptor

	protoDesc, codeDesc, err := le.DescriptorsCache.ByObjectDescriptor(ctx, objDesc)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get descriptors")
	}

	// it's needed to assure that we call method on ref, that has same prototype as proxy, that we import in contract code
	if request.Prototype != nil && !request.Prototype.Equal(*protoDesc.HeadRef()) {
		err := errors.New("proxy call error: try to call method of prototype as method of another prototype")
		errResBuf, err := foundation.MarshalMethodErrorResult(err)
		if err != nil {
			return nil, errors.Wrap(err, "couldn't marshal result")
		}

		return requestresult.New(errResBuf, *objDesc.HeadRef()), nil
	}

	executor, err := le.MachinesManager.GetExecutor(codeDesc.MachineType())
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get executor")
	}

	transcript.LogicContext = le.genLogicCallContext(ctx, transcript, protoDesc, codeDesc)

	newData, result, err := executor.CallMethod(
		ctx, transcript.LogicContext, *codeDesc.Ref(), objDesc.Memory(), request.Method, request.Arguments,
	)

	if err != nil {
		return nil, errors.Wrap(err, "executor error")
	}
	if len(result) == 0 {
		return nil, errors.New("return of method is empty")
	}
	if len(newData) == 0 {
		return nil, errors.New("object state is empty")
	}

	res := requestresult.New(result, *objDesc.HeadRef())

	if request.Immutable {
		return res, nil
	}

	switch {
	case transcript.Deactivate:
		res.SetDeactivate(objDesc)
	case !bytes.Equal(objDesc.Memory(), newData):
		res.SetAmend(objDesc, newData)
	}

	return res, nil
}

func (le *logicExecutor) ExecuteConstructor(
	ctx context.Context, transcript *common.Transcript,
) (
	artifacts.RequestResult, error,
) {
	ctx, span := instracer.StartSpan(ctx, "LogicRunner.executeConstructorCall")
	defer span.Finish()

	inslogger.FromContext(ctx).Debug("Executing constructor")

	request := transcript.Request

	if request.Prototype == nil {
		return nil, errors.New("prototype reference is required")
	}

	protoDesc, codeDesc, err := le.DescriptorsCache.ByPrototypeRef(ctx, *request.Prototype)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get descriptors")
	}

	executor, err := le.MachinesManager.GetExecutor(codeDesc.MachineType())
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get executor")
	}

	transcript.LogicContext = le.genLogicCallContext(ctx, transcript, protoDesc, codeDesc)

	newData, result, err := executor.CallConstructor(ctx, transcript.LogicContext, *codeDesc.Ref(), request.Method, request.Arguments)
	if err != nil {
		return nil, errors.Wrap(err, "executor error")
	}
	if len(result) == 0 {
		return nil, errors.New("return of constructor is empty")
	}

	res := requestresult.New(result, *transcript.Request.Object)
	if newData != nil {
		res.SetActivate(*request.Base, *request.Prototype, newData)
	}
	return res, nil
}

func (le *logicExecutor) genLogicCallContext(
	ctx context.Context,
	transcript *common.Transcript,
	protoDesc artifacts.PrototypeDescriptor,
	codeDesc artifacts.CodeDescriptor,
) *insolar.LogicCallContext {
	request := transcript.Request
	reqRef := transcript.RequestRef

	res := insolar.NewLogicCallContext(
		insolar.ExecuteCallMode,
		&reqRef,
		nil,
		nil,
		protoDesc.HeadRef(),
		codeDesc.Ref(),
		&request.Caller,
		&request.CallerPrototype,
		inslogger.TraceID(ctx),
		func(ctx context.Context, pn insolar.PulseNumber) (i insolar.Pulse, e error) {
			return le.PulseAccessor.ForPulseNumber(ctx, pn)
		},
	)

	if oDesc := transcript.ObjectDescriptor; oDesc != nil {
		res.Parent = oDesc.Parent()
		// should be the same as request.Object
		res.Callee = oDesc.HeadRef()
	} else {
		res.Callee = transcript.Request.Object
	}

	return res
}
