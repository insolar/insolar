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

package getobject

import (
	"context"

	"github.com/insolar/insolar/conveyor/generator/generator"
	"github.com/insolar/insolar/conveyor/fsm"
)

type CustomEvent struct{}
type CustomPayload struct{}
type CustomAdapterResponsePayload struct{}
type CustomAdapterHelper struct{}

const (
	InitState fsm.ElementState = iota
	WaitingPresent
	CheckingJet
	WaitingCheckingJet
	FetchingJet
	WaitingFetchingJet
	InvokeWaitingHotData
	WaitingHotData
	CheckingIndex
	WaitingCheckingIndex
	FetchingIndex
	WaitingFetchingIndex
	CheckingState
	WaitingCheckingState
	CheckingJetForState
	WaitingCheckingJetForState
	FetchingJetForState
	WaitingFetchingJetForState
	FetchingState
	WaitingFetchingState
	Result
)

func Register() {
	generator.AddMachine("GetObjectStateMachine").

		InitFuture(InitFuture, WaitingPresent).
		MigrationFuturePresent(WaitingPresent, MigrateToPresent, CheckingJet).
		Init(Init, CheckingJet).

		Transition(CheckingJet, GetJet, WaitingCheckingJet).
		AdapterResponse(CheckingJet, GetJetResponse, FetchingJet, InvokeWaitingHotData).
		// AdapterResponsePast(CheckingJet, GetJetResponse1, FetchingJet, InvokeWaitingHotData).

		Transition(FetchingJet, FetchJet, WaitingFetchingJet).
		AdapterResponse(FetchingJet, FetchJetResponse, InvokeWaitingHotData).

		Transition(InvokeWaitingHotData, WaitHotData, WaitingHotData).
		AdapterResponse(InvokeWaitingHotData, WaitHotDataResponse, CheckingIndex).

		Transition(CheckingIndex, CheckIndex, WaitingCheckingIndex).
		AdapterResponse(CheckingIndex, WaitCheckIndex, CheckingState, FetchingIndex).

		Transition(FetchingIndex, FetchIndex, WaitingFetchingIndex).
		AdapterResponse(FetchingIndex, WaitFetchIndex, CheckingState).

		Transition(CheckingState, CheckState, WaitingCheckingState).
		AdapterResponse(CheckingState, WaitCheckState, Result, CheckingJetForState).

		Transition(CheckingJetForState, CheckJetForState, WaitingCheckingJetForState).
		AdapterResponse(CheckingJetForState, WaitCheckJetForState, FetchingState, FetchingJetForState).

		Transition(FetchingJetForState, FetchJetForState, WaitingFetchingJetForState).
		AdapterResponse(FetchingJetForState, WaitFetchJetForState, FetchingState).

		Transition(FetchingState, FetchState, WaitingFetchingState).
		AdapterResponse(FetchingState, WaitFetchState, Result)
}

func InitFuture(ctx context.Context, helper fsm.SlotElementHelper, input CustomEvent, payload interface{}) (fsm.ElementState, *CustomPayload) {
	helper.DeactivateTill(fsm.Response)
	return fsm.ElementState(WaitingPresent), payload.(*CustomPayload)
}
func MigrateToPresent(ctx context.Context, helper fsm.SlotElementHelper, input CustomEvent, payload *CustomPayload) fsm.ElementState {
	return fsm.ElementState(CheckingJet)
}
func Init(ctx context.Context, helper fsm.SlotElementHelper, input CustomEvent, payload interface{}) (fsm.ElementState, *CustomPayload) {
	return fsm.ElementState(CheckingJet), payload.(*CustomPayload)
}

func GetJet(ctx context.Context, helper fsm.SlotElementHelper, input CustomEvent, payload *CustomPayload, adapterHelper CustomAdapterHelper) fsm.ElementState {
	// todo invoke adapter
	helper.DeactivateTill(fsm.Response)
	return fsm.ElementState(WaitingCheckingJet)
}

func GetJetResponse(ctx context.Context, helper fsm.SlotElementHelper, input CustomEvent, payload *CustomPayload, respPayload CustomAdapterResponsePayload) fsm.ElementState {
	// todo if found
	return fsm.ElementState(InvokeWaitingHotData)
	// todo else
	return fsm.ElementState(FetchingJet)
}

func FetchJet(ctx context.Context, helper fsm.SlotElementHelper, input CustomEvent, payload *CustomPayload, adapterHelper CustomAdapterHelper) fsm.ElementState {
	// helper.adapters.X(WaitingFetchingJet)
	// todo invoke adapter
	helper.DeactivateTill(fsm.Response)
	return fsm.ElementState(WaitingFetchingJet)
}

func FetchJetResponse(ctx context.Context, helper fsm.SlotElementHelper, input CustomEvent, payload *CustomPayload, respPayload CustomAdapterResponsePayload) fsm.ElementState {
	// todo update payload
	return fsm.ElementState(InvokeWaitingHotData)
}

func WaitHotData(ctx context.Context, helper fsm.SlotElementHelper, input CustomEvent, payload *CustomPayload, adapterHelper CustomAdapterHelper) fsm.ElementState {
	// todo invoke adapter
	helper.DeactivateTill(fsm.Response)
	return fsm.ElementState(WaitingHotData)
}

func WaitHotDataResponse(ctx context.Context, helper fsm.SlotElementHelper, input CustomEvent, payload *CustomPayload, respPayload CustomAdapterResponsePayload) fsm.ElementState {
	// todo update payload
	return fsm.ElementState(CheckingIndex)
}

func CheckIndex(ctx context.Context, helper fsm.SlotElementHelper, input CustomEvent, payload *CustomPayload, adapterHelper CustomAdapterHelper) fsm.ElementState {
	// todo invoke adapter
	helper.DeactivateTill(fsm.Response)
	return fsm.ElementState(WaitingCheckingIndex)
}

func WaitCheckIndex(ctx context.Context, helper fsm.SlotElementHelper, input CustomEvent, payload *CustomPayload, respPayload CustomAdapterResponsePayload) fsm.ElementState {
	// todo if found
	return fsm.ElementState(CheckingState)
	// todo else
	return fsm.ElementState(FetchingIndex)
}

func FetchIndex(ctx context.Context, helper fsm.SlotElementHelper, input CustomEvent, payload *CustomPayload, adapterHelper CustomAdapterHelper) fsm.ElementState {
	// todo invoke adapter
	helper.DeactivateTill(fsm.Response)
	return fsm.ElementState(WaitingFetchingIndex)
}

func WaitFetchIndex(ctx context.Context, helper fsm.SlotElementHelper, input CustomEvent, payload *CustomPayload, respPayload CustomAdapterResponsePayload) fsm.ElementState {
	// todo update payload
	return fsm.ElementState(CheckingState)
}

func CheckState(ctx context.Context, helper fsm.SlotElementHelper, input CustomEvent, payload *CustomPayload, adapterHelper CustomAdapterHelper) fsm.ElementState {
	// todo invoke adapter
	helper.DeactivateTill(fsm.Response)
	return fsm.ElementState(WaitingCheckingState)
}

func WaitCheckState(ctx context.Context, helper fsm.SlotElementHelper, input CustomEvent, payload *CustomPayload, respPayload CustomAdapterResponsePayload) fsm.ElementState {
	// todo if found
	return fsm.ElementState(Result)
	// todo else
	return fsm.ElementState(CheckingJetForState)
}

func CheckJetForState(ctx context.Context, helper fsm.SlotElementHelper, input CustomEvent, payload *CustomPayload, adapterHelper CustomAdapterHelper) fsm.ElementState {
	// todo invoke adapter
	helper.DeactivateTill(fsm.Response)
	return fsm.ElementState(WaitingCheckingJetForState)
}

func WaitCheckJetForState(ctx context.Context, helper fsm.SlotElementHelper, input CustomEvent, payload *CustomPayload, respPayload CustomAdapterResponsePayload) fsm.ElementState {
	// todo if found
	return fsm.ElementState(FetchingState)
	// todo else
	return fsm.ElementState(FetchingJetForState)
}

func FetchJetForState(ctx context.Context, helper fsm.SlotElementHelper, input CustomEvent, payload *CustomPayload, adapterHelper CustomAdapterHelper) fsm.ElementState {
	// todo invoke adapter
	helper.DeactivateTill(fsm.Response)
	return fsm.ElementState(WaitingFetchingJetForState)
}

func WaitFetchJetForState(ctx context.Context, helper fsm.SlotElementHelper, input CustomEvent, payload *CustomPayload, respPayload CustomAdapterResponsePayload) fsm.ElementState {
	// todo update payload
	return fsm.ElementState(FetchingState)
}

func FetchState(ctx context.Context, helper fsm.SlotElementHelper, input CustomEvent, payload *CustomPayload, adapterHelper CustomAdapterHelper) fsm.ElementState {
	// todo invoke adapter
	helper.DeactivateTill(fsm.Response)
	return fsm.ElementState(WaitingFetchingState)
}

func WaitFetchState(ctx context.Context, helper fsm.SlotElementHelper, input CustomEvent, payload *CustomPayload, respPayload CustomAdapterResponsePayload) fsm.ElementState {
	// todo update payload
	return fsm.ElementState(Result)
}
