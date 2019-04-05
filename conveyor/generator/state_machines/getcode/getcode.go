/*
 *    Copyright 2019 Insolar Technologies
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package getcode

import (
	"context"
	"fmt"

	"github.com/insolar/insolar/conveyor/adapter"
	"github.com/insolar/insolar/conveyor/fsm"
	"github.com/insolar/insolar/conveyor/generator/generator"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/artifactmanager"
)

const (
	InitState fsm.ElementState = iota
	WaitingPresent
	GettingCode
	WaitingGettingCode
	ReturningResult
	WaitingReturningResult
)

type GetCodePayload struct {
	parcel insolar.Parcel
	future insolar.ConveyorFuture
	reply  insolar.Reply
	err    error
}

// Register adds get code state machine to generator
func Register(g *generator.Generator) {
	g.AddMachine("GetCode").
		InitFuture(InitFuture, WaitingPresent, 0).
		MigrationFuturePresent(WaitingPresent, MigrateToPresent, GettingCode).
		Init(Init, GettingCode, 0).
		Transition(GettingCode, GetCode, WaitingGettingCode, 0).
		AdapterResponse(GettingCode, GetCodeResponse, ReturningResult).
		Transition(ReturningResult, ReturnResult, WaitingReturningResult, 0).
		AdapterResponse(ReturningResult, ReturnResultResponse, 0)
}

func InitFuture(ctx context.Context, helper fsm.SlotElementHelper, input interface{}, payload interface{}) (fsm.ElementState, *GetCodePayload) {
	getCodePayload := GetCodePayload{}
	parcel, ok := helper.GetInputEvent().(insolar.Parcel)
	if !ok {
		inslogger.FromContext(ctx).Warnf("[ ParseInputEvent ] InputEvent must be insolar.Parcel. Actual: %+v", helper.GetInputEvent())
		return 0, nil
	}
	getCodePayload.parcel = parcel
	helper.DeactivateTill(fsm.Response)
	return fsm.ElementState(WaitingPresent), &getCodePayload
}
func MigrateToPresent(ctx context.Context, helper fsm.SlotElementHelper, input interface{}, payload *GetCodePayload) fsm.ElementState {
	return fsm.ElementState(GettingCode)
}

func Init(ctx context.Context, helper fsm.SlotElementHelper, input interface{}, payload interface{}) (fsm.ElementState, *GetCodePayload) {
	fmt.Println("hi love, we are here")
	getCodePayload := GetCodePayload{}
	parcel, ok := helper.GetInputEvent().(insolar.Parcel)
	if !ok {
		inslogger.FromContext(ctx).Warnf("[ GetCode.Init ] InputEvent must be insolar.Parcel. Actual: %+v", helper.GetInputEvent())
		return 0, nil
	}
	getCodePayload.parcel = parcel
	return fsm.ElementState(GettingCode), &getCodePayload
}

func GetCode(ctx context.Context, helper fsm.SlotElementHelper, input interface{}, payload *GetCodePayload, adapterHelper artifactmanager.GetCodeHelper) fsm.ElementState {
	fmt.Println("hi love, GetCode")
	err := adapterHelper.GetCode(helper, payload.parcel, uint32(GettingCode))
	if err != nil {
		fmt.Println("hi love, GetCode err", err)
		return 0
	}
	// helper.DeactivateTill(fsm.Response) - this should be automatically done in helper
	return fsm.ElementState(WaitingGettingCode)
}

func GetCodeResponse(ctx context.Context, helper fsm.SlotElementHelper, input interface{}, payload *GetCodePayload, respPayload artifactmanager.GetCodeResp) fsm.ElementState {
	fmt.Println("hi love, GetCodeResponse")
	payload.err = respPayload.Err
	payload.reply = respPayload.Reply
	fmt.Println("hi love, GetCodeResponse all well")
	return fsm.ElementState(ReturningResult)
}

func ReturnResult(ctx context.Context, helper fsm.SlotElementHelper, input interface{}, payload *GetCodePayload, adapterHelper adapter.SendResponseHelper) fsm.ElementState {
	var err error
	fmt.Println("hi love, ReturnResult")
	if payload.err != nil {
		// TODO: return error to future
		err = adapterHelper.SendResponse(helper, nil, helper.GetResponseFuture(), uint32(ReturningResult))
	} else {
		err = adapterHelper.SendResponse(helper, payload.reply, helper.GetResponseFuture(), uint32(ReturningResult))
	}
	if err != nil {
		fmt.Println("hi love, ReturnResult err", err)
		return 0
	}
	// helper.DeactivateTill(fsm.Response) - this should be automatically done in helper
	fmt.Println("hi love, ReturnResult all well")
	return fsm.ElementState(WaitingReturningResult)
}

func ReturnResultResponse(ctx context.Context, helper fsm.SlotElementHelper, input interface{}, payload *GetCodePayload, respPayload interface{}) fsm.ElementState {
	fmt.Println("hi love, ReturnResultResponse")
	switch res := respPayload.(type) {
	case string, error:
	default:
		fmt.Println("hi love, ReturnResultResponse err")
		inslogger.FromContext(ctx).Errorf("GetCode: unexpected reply: %T", res)
	}
	fmt.Println("hi love, ReturnResultResponse all well")
	return 0
}
