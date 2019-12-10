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

package s_contract_runner

import (
	"bytes"
	"context"
	"fmt"
	"sync"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/conveyor/smachine"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/logicrunner/common"
	"github.com/insolar/insolar/logicrunner/goplugin/rpctypes"
	"github.com/insolar/insolar/logicrunner/logicexecutor"
	"github.com/insolar/insolar/logicrunner/machinesmanager"
	"github.com/insolar/insolar/logicrunner/requestresult"
	"github.com/insolar/insolar/logicrunner/s_contract_runner/outgoing"
)

type ContractCallType uint8

const (
	ContractCallUnknown ContractCallType = iota
	ContractCallMutable
	ContractCallImmutable
	ContractCallSaga
)

type CallResult interface{}

type ContractRunnerService interface {
	ContractRunnerRPCMethods

	ClassifyCall(request *record.IncomingRequest) ContractCallType
	ExecutionStart(ctx context.Context, transcript *common.Transcript) (*ContractExecutionStateUpdate, error)
	ExecutionContinue(ctx context.Context, requestReference insolar.Reference, result interface{}) (*ContractExecutionStateUpdate, error)
}

type ContractRunnerRPCMethods interface {
	GetCode(rpctypes.UpGetCodeReq, *rpctypes.UpGetCodeResp) error
	RouteCall(rpctypes.UpRouteReq, *rpctypes.UpRouteResp) error
	SaveAsChild(rpctypes.UpSaveAsChildReq, *rpctypes.UpSaveAsChildResp) error
	DeactivateObject(rpctypes.UpDeactivateObjectReq, *rpctypes.UpDeactivateObjectResp) error
}

type ContractRunnerServiceAdapter struct {
	svc  ContractRunnerService
	exec smachine.ExecutionAdapter
}

func (a *ContractRunnerServiceAdapter) PrepareSync(ctx smachine.ExecutionContext, fn func(svc ContractRunnerService)) smachine.SyncCallRequester {
	return a.exec.PrepareSync(ctx, func(_ interface{}) smachine.AsyncResultFunc {
		fn(a.svc)
		return nil
	})
}

func (a *ContractRunnerServiceAdapter) PrepareAsync(ctx smachine.ExecutionContext, fn func(svc ContractRunnerService) smachine.AsyncResultFunc) smachine.AsyncCallRequester {
	return a.exec.PrepareAsync(ctx, func(_ interface{}) smachine.AsyncResultFunc {
		return fn(a.svc)
	})
}

type ExecutionContext struct {
	transcript *common.Transcript

	deactivate bool
	output     chan *ContractExecutionStateUpdate
	input      chan interface{}
}

func NewExecution(transcript *common.Transcript) *ExecutionContext {
	return &ExecutionContext{
		transcript: transcript,
		output:     make(chan *ContractExecutionStateUpdate, 1),
		input:      make(chan interface{}, 1),
	}
}

func (c *ExecutionContext) Error(err error) bool {
	c.output <- &ContractExecutionStateUpdate{
		Type:  ContractError,
		Error: err,
	}

	return true
}
func (c *ExecutionContext) ErrorString(text string) bool {
	return c.Error(errors.New(text))
}
func (c *ExecutionContext) ErrorWrapped(err error, text string) bool {
	return c.Error(errors.Wrap(err, text))
}

func (c *ExecutionContext) ExternalCall(event outgoing.RPCEvent) bool {
	c.output <- &ContractExecutionStateUpdate{
		Type:     ContractOutgoingCall,
		Outgoing: event,
	}
	return true
}
func (c *ExecutionContext) Result(result *requestresult.RequestResult) bool {
	c.output <- &ContractExecutionStateUpdate{
		Type:   ContractDone,
		Result: result,
	}
	return true
}

func (c *ExecutionContext) Stop() {
	close(c.input)
	close(c.output)
}

type contractRunnerService struct {
	LogicExecutor    logicexecutor.LogicExecutor
	MachinesManager  machinesmanager.MachinesManager
	DescriptorsCache artifacts.DescriptorsCache

	executionsLock sync.Mutex
	executions     map[insolar.Reference]*ExecutionContext
}

func CreateContractRunner(
	executor logicexecutor.LogicExecutor,
	manager machinesmanager.MachinesManager,
	artifactManager artifacts.Client,
) ContractRunnerService {
	return &contractRunnerService{
		LogicExecutor:    executor,
		MachinesManager:  manager,
		DescriptorsCache: artifacts.NewDescriptorsCache(artifactManager),

		executions: make(map[insolar.Reference]*ExecutionContext),
	}
}

func CreateContractRunnerService(
	contractRunner ContractRunnerService,
) *ContractRunnerServiceAdapter {
	ctx := context.Background()

	ae, ch := smachine.NewCallChannelExecutor(ctx, -1, false, 16)
	smachine.StartDynamicChannelWorker(ctx, ch, nil)

	return &ContractRunnerServiceAdapter{
		svc:  contractRunner,
		exec: smachine.NewExecutionAdapter("ArtifactClientService", ae),
	}
}

func (c contractRunnerService) ClassifyCall(request *record.IncomingRequest) ContractCallType {
	switch {
	// case request.ReturnMode == recornd.ReturnSaga && false:
	// 	if !request.Immutable && request.CallType == record.CTMethod {
	// 		return ContractCallSaga
	// 	} else {
	// 		return ContractCallUnknown
	// 	}
	// case request.Immutable:
	// 	return ContractCallImmutable
	default:
		return ContractCallMutable
	}
}

func (c *contractRunnerService) getExecution(request insolar.Reference) *ExecutionContext {
	c.executionsLock.Lock()
	defer c.executionsLock.Unlock()

	return c.executions[request]
}

type ExecuteFuncType func(context.Context, insolar.Reference) bool

func (c *contractRunnerService) executeConstructor(ctx context.Context, requestReference insolar.Reference) bool {
	var (
		executionContext = c.getExecution(requestReference)
		request          = executionContext.transcript.Request
	)

	protoDesc, codeDesc, err := c.DescriptorsCache.ByPrototypeRef(ctx, *request.Prototype)
	if err != nil {
		return executionContext.ErrorWrapped(err, "couldn't get descriptors")
	}

	executor, err := c.MachinesManager.GetExecutor(codeDesc.MachineType())
	if err != nil {
		return executionContext.ErrorWrapped(err, "couldn't get executor")
	}

	logicContext := logicexecutor.GenerateCallContext(ctx, executionContext.transcript, protoDesc, codeDesc)

	newData, result, err := executor.CallConstructor(ctx, logicContext, *codeDesc.Ref(), request.Method, request.Arguments)
	if err != nil {
		return executionContext.ErrorWrapped(err, "execution error")
	}
	if len(result) == 0 {
		return executionContext.ErrorString("return of constructor is empty")
	}

	// form and return result
	res := requestresult.New(result, executionContext.transcript.RequestRef)
	if newData != nil {
		res.SetActivate(*request.Base, *request.Prototype, newData)
	}

	return executionContext.Result(res)
}

func (c *contractRunnerService) executeMethod(ctx context.Context, requestReference insolar.Reference) bool {
	var (
		executionContext = c.getExecution(requestReference)
		request          = executionContext.transcript.Request
		objectDescriptor = executionContext.transcript.ObjectDescriptor

		codeDescriptor      artifacts.CodeDescriptor
		prototypeDescriptor artifacts.PrototypeDescriptor
	)

	prototypeReference, err := objectDescriptor.Prototype()
	if err != nil {
		return executionContext.ErrorWrapped(err, "couldn't get prototype reference")
	}

	prototypeDescriptor, codeDescriptor, err = c.DescriptorsCache.ByPrototypeRef(ctx, *prototypeReference)
	if err != nil {
		return executionContext.ErrorWrapped(err, "couldn't get descriptors")
	}

	executor, err := c.MachinesManager.GetExecutor(codeDescriptor.MachineType())
	if err != nil {
		return executionContext.ErrorWrapped(err, "couldn't get executor")
	}

	logicContext := logicexecutor.GenerateCallContext(ctx, executionContext.transcript, prototypeDescriptor, codeDescriptor)

	newData, result, err := executor.CallMethod(
		ctx, logicContext, *codeDescriptor.Ref(), objectDescriptor.Memory(), request.Method, request.Arguments,
	)
	if err != nil {
		return executionContext.ErrorWrapped(err, "execution error")
	}
	if len(result) == 0 {
		return executionContext.ErrorString("return of constructor is empty")
	}
	if len(newData) == 0 {
		return executionContext.ErrorString("object state is empty")
	}

	// form and return result
	res := requestresult.New(result, *objectDescriptor.HeadRef())

	if !request.Immutable {
		switch {
		case executionContext.deactivate:
			res.SetDeactivate(objectDescriptor)
		case !bytes.Equal(objectDescriptor.Memory(), newData):
			res.SetAmend(objectDescriptor, newData)
		}
	}

	return executionContext.Result(res)
}

type ContractExecutionStateUpdateType int

const (
	_ ContractExecutionStateUpdateType = iota
	ContractError
	ContractOutgoingCall
	ContractDone
)

type ContractExecutionStateUpdate struct {
	Type  ContractExecutionStateUpdateType
	Error error

	Result   *requestresult.RequestResult
	Outgoing outgoing.RPCEvent
}

func (c *contractRunnerService) getExecutionContext(requestReference insolar.Reference) *ExecutionContext {
	c.executionsLock.Lock()
	defer c.executionsLock.Unlock()

	return c.executions[requestReference]
}

func (c *contractRunnerService) createExecutionContext(transcript *common.Transcript) (*ExecutionContext, bool) {
	c.executionsLock.Lock()
	defer c.executionsLock.Unlock()

	if val, ok := c.executions[transcript.RequestRef]; ok {
		return val, ok
	}
	c.executions[transcript.RequestRef] = NewExecution(transcript)

	return c.executions[transcript.RequestRef], false
}

func (c *contractRunnerService) stopExecution(requestReference insolar.Reference) error {
	c.executionsLock.Lock()
	defer c.executionsLock.Unlock()

	if val, ok := c.executions[requestReference]; ok {
		delete(c.executions, requestReference)
		val.Stop()
	}

	return nil
}

func (c *contractRunnerService) waitForReply(requestReference insolar.Reference) (*ContractExecutionStateUpdate, error) {
	executionContext := c.getExecutionContext(requestReference)
	if executionContext == nil {
		panic("failed to find ExecutionContext")
	}

	switch update := <-executionContext.output; update.Type {
	case ContractDone:
		_ = c.stopExecution(requestReference)
		fallthrough
	case ContractError, ContractOutgoingCall:
		return update, nil
	default:
		panic(fmt.Sprintf("unknown return type %v", update.Type))
	}
}

func (c *contractRunnerService) executionRecover(ctx context.Context, requestReference insolar.Reference) {
	if r := recover(); r != nil {
		// replace with custom error, not RecoverSlotPanicWithStack
		err := smachine.RecoverSlotPanicWithStack("ContractRunnerService panic", r, nil)

		executionContext := c.getExecution(requestReference)
		if executionContext == nil {
			inslogger.FromContext(ctx).Errorf("[executionRecover] Failed to find a job execution context %s", requestReference.String())
			inslogger.FromContext(ctx).Errorf("[executionRecover] Failed to execute a job, panic: %v", r)
			return
		}

		executionContext.Error(err)
	}
}

// means - create new Job for execution, that'll barge in two cases:
// 1) create outgoing call
// 2) finished execution
func (c *contractRunnerService) ExecutionStart(ctx context.Context, transcript *common.Transcript) (*ContractExecutionStateUpdate, error) {
	requestReference := transcript.RequestRef

	if _, ok := c.createExecutionContext(transcript); ok {
		return nil, errors.Errorf("request %s already executed", requestReference.String())
	}

	var executeFunc ExecuteFuncType = c.executeMethod
	switch transcript.Request.CallType {
	case record.CTMethod:
		executeFunc = c.executeMethod
	case record.CTSaveAsChild:
		executeFunc = c.executeConstructor
	}

	go func() {
		defer c.executionRecover(ctx, requestReference)

		executeFunc(ctx, requestReference)
	}()

	return c.waitForReply(requestReference)
}

func (c *contractRunnerService) ExecutionContinue(ctx context.Context, requestReference insolar.Reference, result interface{}) (*ContractExecutionStateUpdate, error) {
	executionContext := c.getExecutionContext(requestReference)
	if executionContext == nil {
		return nil, errors.Errorf("request %s not found", requestReference.String())
	}

	executionContext.input <- result

	return c.waitForReply(requestReference)
}

func (c *contractRunnerService) GetCode(in rpctypes.UpGetCodeReq, out *rpctypes.UpGetCodeResp) error {
	requestReference := in.Request
	executionContext := c.getExecutionContext(requestReference)
	if executionContext == nil {
		panic("failed to find ExecutionContext")
	}

	event := outgoing.NewRPCBuilder(in.Request, in.Callee).
		GetCode(in.Code)
	executionContext.ExternalCall(event)

	rawValue := <-executionContext.input

	switch val := rawValue.(type) {
	case []byte:
		out.Code = val
	case error:
		return val
	default:
		panic(fmt.Sprintf("GetCode result unexpected type %T", val))
	}

	return nil
}

func (c *contractRunnerService) RouteCall(in rpctypes.UpRouteReq, out *rpctypes.UpRouteResp) error {
	requestReference := in.Request
	executionContext := c.getExecutionContext(requestReference)
	if executionContext == nil {
		panic("failed to find ExecutionContext")
	}

	event := outgoing.RPCEvent(
		outgoing.NewRPCBuilder(in.Request, in.Callee).
			RouteCall(in.Object, in.Prototype, in.Method, in.Arguments).
			SetImmutable(in.Immutable).
			SetSaga(in.Saga),
	)
	executionContext.ExternalCall(event)

	rawValue := <-executionContext.input

	switch val := rawValue.(type) {
	case insolar.Arguments:
		out.Result = val
	case []uint8:
		out.Result = insolar.Arguments(val)
	case error:
		return val
	default:
		panic(fmt.Sprintf("RouteCall result unexpected type %T", val))

	}

	return nil
}

func (c *contractRunnerService) SaveAsChild(in rpctypes.UpSaveAsChildReq, out *rpctypes.UpSaveAsChildResp) error {
	requestReference := in.Request
	executionContext := c.getExecutionContext(requestReference)
	if executionContext == nil {
		panic("failed to find ExecutionContext")
	}

	event := outgoing.NewRPCBuilder(in.Request, in.Callee).
		SaveAsChild(in.Prototype, in.ConstructorName, in.ArgsSerialized)
	executionContext.ExternalCall(event)

	rawValue := <-executionContext.input

	switch val := rawValue.(type) {
	case insolar.Arguments:
		out.Result = val
	case []uint8:
		out.Result = insolar.Arguments(val)
	case error:
		return val
	default:
		panic(fmt.Sprintf("SaveAsChild result unexpected type %T", val))
	}

	return nil
}

func (c *contractRunnerService) DeactivateObject(in rpctypes.UpDeactivateObjectReq, out *rpctypes.UpDeactivateObjectResp) error {
	requestReference := in.Request
	executionContext := c.getExecutionContext(requestReference)
	if executionContext == nil {
		panic("failed to find ExecutionContext")
	}

	event := outgoing.NewRPCBuilder(in.Request, in.Callee).
		Deactivate()
	executionContext.ExternalCall(event)

	rawValue := <-executionContext.input

	switch val := rawValue.(type) {
	case nil:
		return nil
	case error:
		return val
	default:
		panic(fmt.Sprintf("Deactivate result unexpected type %T", val))
	}

	return nil
}
