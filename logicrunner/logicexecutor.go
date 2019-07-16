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

package logicrunner

import (
	"bytes"
	"context"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/logicrunner/artifacts"
)

//go:generate minimock -i github.com/insolar/insolar/logicrunner.LogicExecutor -o ./ -s _mock.go

type LogicExecutor interface {
	Execute(ctx context.Context, transcript *Transcript) (*RequestResult, error)
	ExecuteMethod(ctx context.Context, transcript *Transcript) (*RequestResult, error)
	ExecuteConstructor(ctx context.Context, transcript *Transcript) (*RequestResult, error)
}

type logicExecutor struct {
	MachinesManager  MachinesManager            `inject:""`
	DescriptorsCache artifacts.DescriptorsCache `inject:""`
}

func NewLogicExecutor() LogicExecutor {
	return &logicExecutor{}
}

func (le *logicExecutor) Execute(ctx context.Context, transcript *Transcript) (*RequestResult, error) {
	switch transcript.Request.CallType {
	case record.CTMethod:
		return le.ExecuteMethod(ctx, transcript)
	case record.CTSaveAsChild, record.CTSaveAsDelegate:
		return le.ExecuteConstructor(ctx, transcript)
	default:
		return nil, errors.New("Unknown request call type")
	}
}

func (le *logicExecutor) ExecuteMethod(ctx context.Context, transcript *Transcript) (*RequestResult, error) {
	ctx, span := instracer.StartSpan(ctx, "logicExecutor.ExecuteMethod")
	defer span.End()

	request := transcript.Request

	objDesc := transcript.ObjectDescriptor

	protoDesc, codeDesc, err := le.DescriptorsCache.ByObjectDescriptor(ctx, objDesc)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get descriptors")
	}

	// it's needed to assure that we call method on ref, that has same prototype as proxy, that we import in contract code
	if request.Prototype != nil && !request.Prototype.Equal(*protoDesc.HeadRef()) {
		return nil, errors.New("proxy call error: try to call method of prototype as method of another prototype")
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

	res := NewRequestResult(result)
	if request.Immutable {
		return res, nil
	}

	if transcript.Deactivate {
		res.Deactivate()
	} else if !bytes.Equal(objDesc.Memory(), newData) {
		res.Update(newData)
	}
	return res, nil
}

func (le *logicExecutor) ExecuteConstructor(
	ctx context.Context, transcript *Transcript,
) (
	*RequestResult, error,
) {
	ctx, span := instracer.StartSpan(ctx, "LogicRunner.executeConstructorCall")
	defer span.End()

	request := transcript.Request

	if request.Caller.IsEmpty() {
		return nil, errors.New("Call constructor from nowhere")
	}

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

	newData, err := executor.CallConstructor(ctx, transcript.LogicContext, *codeDesc.Ref(), request.Method, request.Arguments)
	if err != nil {
		return nil, errors.Wrap(err, "executor error")
	}

	res := NewRequestResult(nil)
	res.Activate(newData)
	return res, nil
}

func (le *logicExecutor) genLogicCallContext(
	ctx context.Context,
	transcript *Transcript,
	protoDesc artifacts.ObjectDescriptor,
	codeDesc artifacts.CodeDescriptor,
) *insolar.LogicCallContext {
	request := transcript.Request
	reqRef := transcript.RequestRef
	res := &insolar.LogicCallContext{
		Mode: insolar.ExecuteCallMode,

		Request: &reqRef,

		Callee:    nil, // below
		Prototype: protoDesc.HeadRef(),
		Code:      codeDesc.Ref(),

		Caller:          &request.Caller,
		CallerPrototype: &request.CallerPrototype,

		TraceID: inslogger.TraceID(ctx),
	}

	if oDesc := transcript.ObjectDescriptor; oDesc != nil {
		res.Parent = oDesc.Parent()
		// should be the same as request.Object
		res.Callee = oDesc.HeadRef()
	} else {
		res.Callee = &reqRef
	}

	return res
}
