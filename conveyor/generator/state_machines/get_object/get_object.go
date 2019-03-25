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

type Event struct{}
type Payload struct{}
type TAR struct{}
type TA1 struct{}

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
	gen.RegisterTransitionsMap(
		FetchingJet, gen.Asd{InvokeWaitingHotData, InvokeWaitingHotData},
	)

	gen.AddMachine("GetObjectStateMachine").

		TransitionFuture(InitState, InitFuture, WaitingPresent).
		MigrationFuturePresent(WaitingPresent, MigrateToPresent, CheckingJet).
		Transition(InitState, Init, CheckingJet).

		Transition(CheckingJet, GetJet, WaitingCheckingJet).
		AdapterResponse(CheckingJet, GetJetResponse).

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

func InitFuture(helper slot.SlotElementHelper, input Event, payload interface{}) (*Payload, fsm.ElementState) {
	helper.DeactivateTill(slot.Response)
	return payload.(*Payload), fsm.ElementState(WaitingPresent)
}
func MigrateToPresent(input Event, payload *Payload) (*Payload, fsm.ElementState) {
	return payload, fsm.ElementState(CheckingJet)
}
func Init(helper slot.SlotElementHelper, input Event, payload interface{}) (*Payload, fsm.ElementState) {
	return payload.(*Payload), fsm.ElementState(CheckingJet)
}

func GetJet(helper slot.SlotElementHelper, input Event, payload *Payload, adapterHelper TA1) (*Payload, fsm.ElementState) {
	// todo invoke adapter
	helper.DeactivateTill(slot.Response)
	return payload, fsm.ElementState(WaitingCheckingJet)
}

func GetJetResponse(input Event, payload *Payload, respPayload TAR) (*Payload, fsm.ElementState) {
	// todo if found
	return payload, fsm.ElementState(InvokeWaitingHotData)
	// todo else
	return payload, fsm.ElementState(FetchingJet)
}

func FetchJet(helper slot.SlotElementHelper, input Event, payload *Payload, adapterHelper TA1) (*Payload, fsm.ElementState) {
	// todo invoke adapter
	helper.DeactivateTill(slot.Response)
	return payload, fsm.ElementState(WaitingFetchingJet)
}

func FetchJetResponse(input Event, payload *Payload, respPayload TAR) (*Payload, fsm.ElementState) {
	// todo update payload
	return payload, fsm.ElementState(InvokeWaitingHotData)
}

func WaitHotData(helper slot.SlotElementHelper, input Event, payload *Payload, adapterHelper TA1) (*Payload, fsm.ElementState) {
	// todo invoke adapter
	helper.DeactivateTill(slot.Response)
	return payload, fsm.ElementState(WaitingHotData)
}

func WaitHotDataResponse(input Event, payload *Payload, respPayload TAR) (*Payload, fsm.ElementState) {
	// todo update payload
	return payload, fsm.ElementState(CheckingIndex)
}

func CheckIndex(helper slot.SlotElementHelper, input Event, payload *Payload, adapterHelper TA1) (*Payload, fsm.ElementState) {
	// todo invoke adapter
	helper.DeactivateTill(slot.Response)
	return payload, fsm.ElementState(WaitingCheckingIndex)
}

func WaitCheckIndex(input Event, payload *Payload, respPayload TAR) (*Payload, fsm.ElementState) {
	// todo if found
	return payload, fsm.ElementState(CheckingState)
	// todo else
	return payload, fsm.ElementState(FetchingIndex)
}

func FetchIndex(helper slot.SlotElementHelper, input Event, payload *Payload, adapterHelper TA1) (*Payload, fsm.ElementState) {
	// todo invoke adapter
	helper.DeactivateTill(slot.Response)
	return payload, fsm.ElementState(WaitingFetchingIndex)
}

func WaitFetchIndex(input Event, payload *Payload, respPayload TAR) (*Payload, fsm.ElementState) {
	// todo update payload
	return payload, fsm.ElementState(CheckingState)
}

func CheckState(helper slot.SlotElementHelper, input Event, payload *Payload, adapterHelper TA1) (*Payload, fsm.ElementState) {
	// todo invoke adapter
	helper.DeactivateTill(slot.Response)
	return payload, fsm.ElementState(WaitingCheckingState)
}

func WaitCheckState(input Event, payload *Payload, respPayload TAR) (*Payload, fsm.ElementState) {
	// todo if found
	return payload, fsm.ElementState(Result)
	// todo else
	return payload, fsm.ElementState(CheckingJetForState)
}

func CheckJetForState(helper slot.SlotElementHelper, input Event, payload *Payload, adapterHelper TA1) (*Payload, fsm.ElementState) {
	// todo invoke adapter
	helper.DeactivateTill(slot.Response)
	return payload, fsm.ElementState(WaitingCheckingJetForState)
}

func WaitCheckJetForState(input Event, payload *Payload, respPayload TAR) (*Payload, fsm.ElementState) {
	// todo if found
	return payload, fsm.ElementState(FetchingState)
	// todo else
	return payload, fsm.ElementState(FetchingJetForState)
}

func FetchJetForState(helper slot.SlotElementHelper, input Event, payload *Payload, adapterHelper TA1) (*Payload, fsm.ElementState) {
	// todo invoke adapter
	helper.DeactivateTill(slot.Response)
	return payload, fsm.ElementState(WaitingFetchingJetForState)
}

func WaitFetchJetForState(input Event, payload *Payload, respPayload TAR) (*Payload, fsm.ElementState) {
	// todo update payload
	return payload, fsm.ElementState(FetchingState)
}

func FetchState(helper slot.SlotElementHelper, input Event, payload *Payload, adapterHelper TA1) (*Payload, fsm.ElementState) {
	// todo invoke adapter
	helper.DeactivateTill(slot.Response)
	return payload, fsm.ElementState(WaitingFetchingState)
}

func WaitFetchState(input Event, payload *Payload, respPayload TAR) (*Payload, fsm.ElementState) {
	// todo update payload
	return payload, fsm.ElementState(Result)
}

