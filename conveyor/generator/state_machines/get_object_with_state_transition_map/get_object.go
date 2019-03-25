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

package get_object

import (
	"github.com/insolar/insolar/conveyor/generator/generator/gen"

	"github.com/insolar/insolar/conveyor/interfaces/fsm"
	"github.com/insolar/insolar/conveyor/interfaces/slot"
)

type CustomEvent struct{}
type CustomPayload struct{}
type CustomAdapterResponsePayload struct{}
type CustomAdapterHelper struct{}

const (
	InitState gen.ElState = iota
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
	gen.AddMachine("GetObjectStateMachine").
		RegisterTransitionsMap(
			InitState,            gen.Asd{WaitingPresent, CheckingJet},
			WaitingPresent,       gen.Asd{CheckingJet},
			CheckingJet,          gen.Asd{WaitingCheckingJet, FetchingJet},
			FetchingJet,          gen.Asd{WaitingFetchingJet, InvokeWaitingHotData},
			InvokeWaitingHotData, gen.Asd{WaitingHotData, CheckingIndex},
			CheckingIndex,        gen.Asd{WaitingCheckingIndex, CheckingState, FetchingIndex},
			FetchingIndex,        gen.Asd{WaitingFetchingIndex, CheckingState},
			CheckingState,        gen.Asd{WaitingCheckingState, Result, CheckingJetForState},
			CheckingJetForState,  gen.Asd{WaitingCheckingJetForState, FetchingState, FetchingJetForState},
			FetchingJetForState,  gen.Asd{WaitingFetchingJetForState, FetchingState},
			FetchingState,        gen.Asd{WaitingFetchingState, Result},
		).

		InitFuture(InitState, InitFuture).
		MigrationFuturePresent(WaitingPresent, MigrateToPresent).
		Init(InitState, Init).

		Transition(CheckingJet, GetJet).
		AdapterResponse(CheckingJet, GetJetResponse).

		Transition(FetchingJet, FetchJet).
		AdapterResponse(FetchingJet, FetchJetResponse).

		Transition(InvokeWaitingHotData, WaitHotData).
		AdapterResponse(InvokeWaitingHotData, WaitHotDataResponse).

		Transition(CheckingIndex, CheckIndex).
		AdapterResponse(CheckingIndex, WaitCheckIndex).

		Transition(FetchingIndex, FetchIndex).
		AdapterResponse(FetchingIndex, WaitFetchIndex).

		Transition(CheckingState, CheckState).
		AdapterResponse(CheckingState, WaitCheckState).

		Transition(CheckingJetForState, CheckJetForState).
		AdapterResponse(CheckingJetForState, WaitCheckJetForState).

		Transition(FetchingJetForState, FetchJetForState).
		AdapterResponse(FetchingJetForState, WaitFetchJetForState).

		Transition(FetchingState, FetchState).
		AdapterResponse(FetchingState, WaitFetchState)
}

func InitFuture(helper slot.SlotElementHelper, input CustomEvent, payload interface{}) (*CustomPayload, fsm.ElementState) {
	helper.DeactivateTill(slot.Response)
	return payload.(*CustomPayload), fsm.ElementState(WaitingPresent)
}
func MigrateToPresent(input CustomEvent, payload *CustomPayload) (*CustomPayload, fsm.ElementState) {
	return payload, fsm.ElementState(CheckingJet)
}
func Init(helper slot.SlotElementHelper, input CustomEvent, payload interface{}) (*CustomPayload, fsm.ElementState) {
	return payload.(*CustomPayload), fsm.ElementState(CheckingJet)
}

func GetJet(helper slot.SlotElementHelper, input CustomEvent, payload *CustomPayload, adapterHelper CustomAdapterHelper) (*CustomPayload, fsm.ElementState) {
	// todo invoke adapter
	helper.DeactivateTill(slot.Response)
	return payload, fsm.ElementState(WaitingCheckingJet)
}

func GetJetResponse(input CustomEvent, payload *CustomPayload, respPayload CustomAdapterResponsePayload) (*CustomPayload, fsm.ElementState) {
	// todo if found
	return payload, fsm.ElementState(InvokeWaitingHotData)
	// todo else
	return payload, fsm.ElementState(FetchingJet)
}

func FetchJet(helper slot.SlotElementHelper, input CustomEvent, payload *CustomPayload, adapterHelper CustomAdapterHelper) (*CustomPayload, fsm.ElementState) {
	// todo invoke adapter
	helper.DeactivateTill(slot.Response)
	return payload, fsm.ElementState(WaitingFetchingJet)
}

func FetchJetResponse(input CustomEvent, payload *CustomPayload, respPayload CustomAdapterResponsePayload) (*CustomPayload, fsm.ElementState) {
	// todo update payload
	return payload, fsm.ElementState(InvokeWaitingHotData)
}

func WaitHotData(helper slot.SlotElementHelper, input CustomEvent, payload *CustomPayload, adapterHelper CustomAdapterHelper) (*CustomPayload, fsm.ElementState) {
	// todo invoke adapter
	helper.DeactivateTill(slot.Response)
	return payload, fsm.ElementState(WaitingHotData)
}

func WaitHotDataResponse(input CustomEvent, payload *CustomPayload, respPayload CustomAdapterResponsePayload) (*CustomPayload, fsm.ElementState) {
	// todo update payload
	return payload, fsm.ElementState(CheckingIndex)
}

func CheckIndex(helper slot.SlotElementHelper, input CustomEvent, payload *CustomPayload, adapterHelper CustomAdapterHelper) (*CustomPayload, fsm.ElementState) {
	// todo invoke adapter
	helper.DeactivateTill(slot.Response)
	return payload, fsm.ElementState(WaitingCheckingIndex)
}

func WaitCheckIndex(input CustomEvent, payload *CustomPayload, respPayload CustomAdapterResponsePayload) (*CustomPayload, fsm.ElementState) {
	// todo if found
	return payload, fsm.ElementState(CheckingState)
	// todo else
	return payload, fsm.ElementState(FetchingIndex)
}

func FetchIndex(helper slot.SlotElementHelper, input CustomEvent, payload *CustomPayload, adapterHelper CustomAdapterHelper) (*CustomPayload, fsm.ElementState) {
	// todo invoke adapter
	helper.DeactivateTill(slot.Response)
	return payload, fsm.ElementState(WaitingFetchingIndex)
}

func WaitFetchIndex(input CustomEvent, payload *CustomPayload, respPayload CustomAdapterResponsePayload) (*CustomPayload, fsm.ElementState) {
	// todo update payload
	return payload, fsm.ElementState(CheckingState)
}

func CheckState(helper slot.SlotElementHelper, input CustomEvent, payload *CustomPayload, adapterHelper CustomAdapterHelper) (*CustomPayload, fsm.ElementState) {
	// todo invoke adapter
	helper.DeactivateTill(slot.Response)
	return payload, fsm.ElementState(WaitingCheckingState)
}

func WaitCheckState(input CustomEvent, payload *CustomPayload, respPayload CustomAdapterResponsePayload) (*CustomPayload, fsm.ElementState) {
	// todo if found
	return payload, fsm.ElementState(Result)
	// todo else
	return payload, fsm.ElementState(CheckingJetForState)
}

func CheckJetForState(helper slot.SlotElementHelper, input CustomEvent, payload *CustomPayload, adapterHelper CustomAdapterHelper) (*CustomPayload, fsm.ElementState) {
	// todo invoke adapter
	helper.DeactivateTill(slot.Response)
	return payload, fsm.ElementState(WaitingCheckingJetForState)
}

func WaitCheckJetForState(input CustomEvent, payload *CustomPayload, respPayload CustomAdapterResponsePayload) (*CustomPayload, fsm.ElementState) {
	// todo if found
	return payload, fsm.ElementState(FetchingState)
	// todo else
	return payload, fsm.ElementState(FetchingJetForState)
}

func FetchJetForState(helper slot.SlotElementHelper, input CustomEvent, payload *CustomPayload, adapterHelper CustomAdapterHelper) (*CustomPayload, fsm.ElementState) {
	// todo invoke adapter
	helper.DeactivateTill(slot.Response)
	return payload, fsm.ElementState(WaitingFetchingJetForState)
}

func WaitFetchJetForState(input CustomEvent, payload *CustomPayload, respPayload CustomAdapterResponsePayload) (*CustomPayload, fsm.ElementState) {
	// todo update payload
	return payload, fsm.ElementState(FetchingState)
}

func FetchState(helper slot.SlotElementHelper, input CustomEvent, payload *CustomPayload, adapterHelper CustomAdapterHelper) (*CustomPayload, fsm.ElementState) {
	// todo invoke adapter
	helper.DeactivateTill(slot.Response)
	return payload, fsm.ElementState(WaitingFetchingState)
}

func WaitFetchState(input CustomEvent, payload *CustomPayload, respPayload CustomAdapterResponsePayload) (*CustomPayload, fsm.ElementState) {
	// todo update payload
	return payload, fsm.ElementState(Result)
}

