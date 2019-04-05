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

	"github.com/insolar/insolar/conveyor/fsm"
)

type CustomEvent struct{}
type CustomPayload struct{}
type CustomAdapterResponsePayload struct{}

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

type payload struct {
	slotElement
	input CustomEvent
	data  []byte
}

func Init(ctx context.Context, helper fsm.SlotElementHelper, input interface{}) *payload {
	p := &payload{
		input: input.(CustomEvent),
		data:  []byte(""),
	}
	p.
		Transition(CheckingJet, p.GetJet, WaitingCheckingJet).
		AdapterResponse(CheckingJet, p.GetJetResponse, FetchingJet, InvokeWaitingHotData).
		Transition(FetchingJet, p.FetchJet, WaitingFetchingJet).
		AdapterResponse(FetchingJet, p.FetchJetResponse, InvokeWaitingHotData).
		Transition(InvokeWaitingHotData, p.WaitHotData, WaitingHotData).
		AdapterResponse(InvokeWaitingHotData, p.WaitHotDataResponse, CheckingIndex).
		Transition(CheckingIndex, p.CheckIndex, WaitingCheckingIndex).
		AdapterResponse(CheckingIndex, p.WaitCheckIndex, CheckingState, FetchingIndex).
		Transition(FetchingIndex, p.FetchIndex, WaitingFetchingIndex).
		AdapterResponse(FetchingIndex, p.WaitFetchIndex, CheckingState).
		Transition(CheckingState, p.CheckState, WaitingCheckingState).
		AdapterResponse(CheckingState, p.WaitCheckState, Result, CheckingJetForState).
		Transition(CheckingJetForState, p.CheckJetForState, WaitingCheckingJetForState).
		AdapterResponse(CheckingJetForState, p.WaitCheckJetForState, FetchingState, FetchingJetForState).
		Transition(FetchingJetForState, p.FetchJetForState, WaitingFetchingJetForState).
		AdapterResponse(FetchingJetForState, p.WaitFetchJetForState, FetchingState).
		Transition(FetchingState, p.FetchState, WaitingFetchingState).
		AdapterResponse(FetchingState, p.WaitFetchState, Result)
	return &payload{}
}

func (p *payload) GetJet(ctx context.Context, helper fsm.SlotElementHelper) fsm.ElementState {
	// todo invoke adapter
	helper.DeactivateTill(fsm.Response)
	return fsm.ElementState(WaitingCheckingJet)
}

func (p *payload) GetJetResponse(ctx context.Context, helper fsm.SlotElementHelper, respPayload interface{}) fsm.ElementState {
	_ = respPayload.(CustomAdapterResponsePayload)
	// todo if found
	return fsm.ElementState(InvokeWaitingHotData)
	// todo else
	return fsm.ElementState(FetchingJet)
}

func (p *payload) FetchJet(ctx context.Context, helper fsm.SlotElementHelper) fsm.ElementState {
	// helper.adapters.X(WaitingFetchingJet)
	// todo invoke adapter
	helper.DeactivateTill(fsm.Response)
	return fsm.ElementState(WaitingFetchingJet)
}

func (p *payload) FetchJetResponse(ctx context.Context, helper fsm.SlotElementHelper, respPayload interface{}) fsm.ElementState {
	// todo update payload
	return fsm.ElementState(InvokeWaitingHotData)
}

func (p *payload) WaitHotData(ctx context.Context, helper fsm.SlotElementHelper) fsm.ElementState {
	// todo invoke adapter
	helper.DeactivateTill(fsm.Response)
	return fsm.ElementState(WaitingHotData)
}

func (p *payload) WaitHotDataResponse(ctx context.Context, helper fsm.SlotElementHelper, respPayload interface{}) fsm.ElementState {
	// todo update payload
	return fsm.ElementState(CheckingIndex)
}

func (p *payload) CheckIndex(ctx context.Context, helper fsm.SlotElementHelper) fsm.ElementState {
	// todo invoke adapter
	helper.DeactivateTill(fsm.Response)
	return fsm.ElementState(WaitingCheckingIndex)
}

func (p *payload) WaitCheckIndex(ctx context.Context, helper fsm.SlotElementHelper, respPayload interface{}) fsm.ElementState {
	// todo if found
	return fsm.ElementState(CheckingState)
	// todo else
	return fsm.ElementState(FetchingIndex)
}

func (p *payload) FetchIndex(ctx context.Context, helper fsm.SlotElementHelper) fsm.ElementState {
	// todo invoke adapter
	helper.DeactivateTill(fsm.Response)
	return fsm.ElementState(WaitingFetchingIndex)
}

func (p *payload) WaitFetchIndex(ctx context.Context, helper fsm.SlotElementHelper, respPayload interface{}) fsm.ElementState {
	// todo update payload
	return fsm.ElementState(CheckingState)
}

func (p *payload) CheckState(ctx context.Context, helper fsm.SlotElementHelper) fsm.ElementState {
	// todo invoke adapter
	helper.DeactivateTill(fsm.Response)
	return fsm.ElementState(WaitingCheckingState)
}

func (p *payload) WaitCheckState(ctx context.Context, helper fsm.SlotElementHelper, respPayload interface{}) fsm.ElementState {
	// todo if found
	return fsm.ElementState(Result)
	// todo else
	return fsm.ElementState(CheckingJetForState)
}

func (p *payload) CheckJetForState(ctx context.Context, helper fsm.SlotElementHelper) fsm.ElementState {
	// todo invoke adapter
	helper.DeactivateTill(fsm.Response)
	return fsm.ElementState(WaitingCheckingJetForState)
}

func (p *payload) WaitCheckJetForState(ctx context.Context, helper fsm.SlotElementHelper, respPayload interface{}) fsm.ElementState {
	// todo if found
	return fsm.ElementState(FetchingState)
	// todo else
	return fsm.ElementState(FetchingJetForState)
}

func (p *payload) FetchJetForState(ctx context.Context, helper fsm.SlotElementHelper) fsm.ElementState {
	// todo invoke adapter
	helper.DeactivateTill(fsm.Response)
	return fsm.ElementState(WaitingFetchingJetForState)
}

func (p *payload) WaitFetchJetForState(ctx context.Context, helper fsm.SlotElementHelper, respPayload interface{}) fsm.ElementState {
	// todo update payload
	return fsm.ElementState(FetchingState)
}

func (p *payload) FetchState(ctx context.Context, helper fsm.SlotElementHelper) fsm.ElementState {
	// todo invoke adapter
	helper.DeactivateTill(fsm.Response)
	return fsm.ElementState(WaitingFetchingState)
}

func (p *payload) WaitFetchState(ctx context.Context, helper fsm.SlotElementHelper, respPayload interface{}) fsm.ElementState {
	// todo update payload
	return fsm.ElementState(Result)
}
