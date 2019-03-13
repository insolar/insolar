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

package sample

import (
	"github.com/insolar/insolar/conveyor/generator/common"
	"github.com/insolar/insolar/conveyor/interfaces/adapter"
	"github.com/insolar/insolar/conveyor/interfaces/slot"
	"github.com/insolar/insolar/conveyor/interfaces/statemachine"

	"errors"
)

type BaseTestStateMachine struct {}

func (*BaseTestStateMachine) GetTypeID() statemachine.ID {
    return 1
}

func (*BaseTestStateMachine) stateFirst() statemachine.StateID {
    return 1
}
func (*BaseTestStateMachine) stateSecond() statemachine.StateID {
    return 2
}

type RawTestStateMachine struct {
    cleanStateMachine TestStateMachine
}

func RawTestStateMachineFactory() [3]common.StateMachine {
    m := RawTestStateMachine{
        cleanStateMachine: &CleanTestStateMachine{},
    }

    var x = [3][]common.State{}
    // future state machine
    x[0] = append(x[0], common.State{
        Transition: m.initFutureHandler,
        ErrorState: m.errorFutureInit,
    },
    common.State{
        Transition: m.transitFutureFirst,
        AdapterResponse: m.responseFutureFirst,
        ErrorState: m.errorFutureFirst,
        AdapterResponseError: m.errorResponseFutureFirst,
    },
    common.State{
        Transition: m.transitFutureFirstState,
        AdapterResponse: m.responseFutureSecond,
        ErrorState: m.errorFutureSecond,
        AdapterResponseError: m.errorResponseFutureSecond,
    },)

    // present state machine
    x[1] = append(x[1], common.State{
        Transition: m.initPresentHandler,
        ErrorState: m.errorPresentInit,
    },
    common.State{
        Migration: m.migrateFromFutureFirst,
        Transition: m.transitPresentFirst,
        AdapterResponse: m.responsePresentFirst,
        ErrorState: m.errorPresentFirst,
        AdapterResponseError: m.errorResponsePresentFirst,
    },
    common.State{
        Migration: m.migrateFromFutureSecond,
        Transition: m.transitPresentFirstState,
        AdapterResponse: m.responsePresentSecond,
        ErrorState: m.errorPresentSecond,
        AdapterResponseError: m.errorResponsePresentSecond,
    },)

    // past state machine
    x[2] = append(x[2], common.State{
        Transition: m.initPastHandler,
        ErrorState: m.errorPastInit,
    },
    common.State{
        Transition: m.transitPastFirst,
        AdapterResponse: m.responsePastFirst,
        ErrorState: m.errorPastFirst,
        AdapterResponseError: m.errorResponsePastFirst,
    },
    common.State{
        Transition: m.transitPastFirstState,
        AdapterResponse: m.responsePastSecond,
        ErrorState: m.errorPastSecond,
        AdapterResponseError: m.errorResponsePastSecond,
    },)

    return [3]common.StateMachine{
        {
            ID: m.cleanStateMachine.(TestStateMachine).GetTypeID(),
            States: x[0],
        },
        {
            ID: m.cleanStateMachine.(TestStateMachine).GetTypeID(),
            States: x[0],
        },
        {
            ID: m.cleanStateMachine.(TestStateMachine).GetTypeID(),
            States: x[0],
        },
    }
}

func (s *RawTestStateMachine) initPresentHandler(element slot.SlotElementHelper) (interface{}, statemachine.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    payload, state, err := s.cleanStateMachine.initPresentHandler(aInput, element.GetPayload())
    return payload, state, err
}
func (s *RawTestStateMachine) initFutureHandler(element slot.SlotElementHelper) (interface{}, statemachine.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    payload, state, err := s.cleanStateMachine.initFutureHandler(aInput, element.GetPayload())
    return payload, state, err
}
func (s *RawTestStateMachine) initPastHandler(element slot.SlotElementHelper) (interface{}, statemachine.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    payload, state, err := s.cleanStateMachine.initPastHandler(aInput, element.GetPayload())
    return payload, state, err
}
func (s *RawTestStateMachine) errorPresentInit(element slot.SlotElementHelper, err error) (interface{}, statemachine.ElementState) {
    payload, state := s.cleanStateMachine.errorPresentInit(element.GetInputEvent(), element.GetPayload(), err)
    return payload, state
}
func (s *RawTestStateMachine) errorFutureInit(element slot.SlotElementHelper, err error) (interface{}, statemachine.ElementState) {
    payload, state := s.cleanStateMachine.errorFutureInit(element.GetInputEvent(), element.GetPayload(), err)
    return payload, state
}
func (s *RawTestStateMachine) errorPastInit(element slot.SlotElementHelper, err error) (interface{}, statemachine.ElementState) {
    payload, state := s.cleanStateMachine.errorPastInit(element.GetInputEvent(), element.GetPayload(), err)
    return payload, state
}


func (s *RawTestStateMachine) transitPresentFirst(element slot.SlotElementHelper) (interface{}, statemachine.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    payload, state, err := s.cleanStateMachine.transitPresentFirst(aInput, aPayload)
    return payload, state, err
}
func (s *RawTestStateMachine) transitFutureFirst(element slot.SlotElementHelper) (interface{}, statemachine.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    payload, state, err := s.cleanStateMachine.transitFutureFirst(aInput, aPayload)
    return payload, state, err
}
func (s *RawTestStateMachine) transitPastFirst(element slot.SlotElementHelper) (interface{}, statemachine.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    payload, state, err := s.cleanStateMachine.transitPastFirst(aInput, aPayload)
    return payload, state, err
}
func (s *RawTestStateMachine) migrateFromPresentFirst(element slot.SlotElementHelper) (interface{}, statemachine.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    payload, state, err := s.cleanStateMachine.migrateFromPresentFirst(aInput, aPayload)
    return payload, state, err
}
func (s *RawTestStateMachine) migrateFromFutureFirst(element slot.SlotElementHelper) (interface{}, statemachine.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    payload, state, err := s.cleanStateMachine.migrateFromFutureFirst(aInput, aPayload)
    return payload, state, err
}
func (s *RawTestStateMachine) errorPresentFirst(element slot.SlotElementHelper, err error) (interface{}, statemachine.ElementState) {
    payload, state := s.cleanStateMachine.errorPresentFirst(element.GetInputEvent(), element.GetPayload(), err)
    return payload, state
}
func (s *RawTestStateMachine) errorFutureFirst(element slot.SlotElementHelper, err error) (interface{}, statemachine.ElementState) {
    payload, state := s.cleanStateMachine.errorFutureFirst(element.GetInputEvent(), element.GetPayload(), err)
    return payload, state
}
func (s *RawTestStateMachine) errorPastFirst(element slot.SlotElementHelper, err error) (interface{}, statemachine.ElementState) {
    payload, state := s.cleanStateMachine.errorPastFirst(element.GetInputEvent(), element.GetPayload(), err)
    return payload, state
}
func (s *RawTestStateMachine) responsePresentFirst(element slot.SlotElementHelper, ar adapter.IAdapterResponse) (interface{}, statemachine.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    aResponse, ok := ar.GetRespPayload().(TAR)
    if !ok { return nil, 0, errors.New("wrong response type") }
    payload, state, err := s.cleanStateMachine.responsePresentFirst(aInput, aPayload, aResponse)
    return payload, state, err
}
func (s *RawTestStateMachine) responseFutureFirst(element slot.SlotElementHelper, ar adapter.IAdapterResponse) (interface{}, statemachine.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    aResponse, ok := ar.GetRespPayload().(TAR)
    if !ok { return nil, 0, errors.New("wrong response type") }
    payload, state, err := s.cleanStateMachine.responseFutureFirst(aInput, aPayload, aResponse)
    return payload, state, err
}
func (s *RawTestStateMachine) responsePastFirst(element slot.SlotElementHelper, ar adapter.IAdapterResponse) (interface{}, statemachine.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    aResponse, ok := ar.GetRespPayload().(TAR)
    if !ok { return nil, 0, errors.New("wrong response type") }
    payload, state, err := s.cleanStateMachine.responsePastFirst(aInput, aPayload, aResponse)
    return payload, state, err
}
func (s *RawTestStateMachine) errorResponsePresentFirst(element slot.SlotElementHelper, ar adapter.IAdapterResponse, err error) (interface{}, statemachine.ElementState) {
    payload, state := s.cleanStateMachine.errorResponsePresentFirst(element.GetInputEvent(), element.GetPayload(), ar, err)
    return payload, state
}
func (s *RawTestStateMachine) errorResponseFutureFirst(element slot.SlotElementHelper, ar adapter.IAdapterResponse, err error) (interface{}, statemachine.ElementState) {
    payload, state := s.cleanStateMachine.errorResponseFutureFirst(element.GetInputEvent(), element.GetPayload(), ar, err)
    return payload, state
}
func (s *RawTestStateMachine) errorResponsePastFirst(element slot.SlotElementHelper, ar adapter.IAdapterResponse, err error) (interface{}, statemachine.ElementState) {
    payload, state := s.cleanStateMachine.errorResponsePastFirst(element.GetInputEvent(), element.GetPayload(), ar, err)
    return payload, state
}

func (s *RawTestStateMachine) transitPresentFirstState(element slot.SlotElementHelper) (interface{}, statemachine.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    payload, state, err := s.cleanStateMachine.transitPresentFirstState(aInput, aPayload)
    return payload, state, err
}
func (s *RawTestStateMachine) transitFutureFirstState(element slot.SlotElementHelper) (interface{}, statemachine.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    payload, state, err := s.cleanStateMachine.transitFutureFirstState(aInput, aPayload)
    return payload, state, err
}
func (s *RawTestStateMachine) transitPastFirstState(element slot.SlotElementHelper) (interface{}, statemachine.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    payload, state, err := s.cleanStateMachine.transitPastFirstState(aInput, aPayload)
    return payload, state, err
}
func (s *RawTestStateMachine) migrateFromPresentSecond(element slot.SlotElementHelper) (interface{}, statemachine.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    payload, state, err := s.cleanStateMachine.migrateFromPresentSecond(aInput, aPayload)
    return payload, state, err
}
func (s *RawTestStateMachine) migrateFromFutureSecond(element slot.SlotElementHelper) (interface{}, statemachine.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    payload, state, err := s.cleanStateMachine.migrateFromFutureSecond(aInput, aPayload)
    return payload, state, err
}
func (s *RawTestStateMachine) errorPresentSecond(element slot.SlotElementHelper, err error) (interface{}, statemachine.ElementState) {
    payload, state := s.cleanStateMachine.errorPresentSecond(element.GetInputEvent(), element.GetPayload(), err)
    return payload, state
}
func (s *RawTestStateMachine) errorFutureSecond(element slot.SlotElementHelper, err error) (interface{}, statemachine.ElementState) {
    payload, state := s.cleanStateMachine.errorFutureSecond(element.GetInputEvent(), element.GetPayload(), err)
    return payload, state
}
func (s *RawTestStateMachine) errorPastSecond(element slot.SlotElementHelper, err error) (interface{}, statemachine.ElementState) {
    payload, state := s.cleanStateMachine.errorPastSecond(element.GetInputEvent(), element.GetPayload(), err)
    return payload, state
}
func (s *RawTestStateMachine) responsePresentSecond(element slot.SlotElementHelper, ar adapter.IAdapterResponse) (interface{}, statemachine.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    aResponse, ok := ar.GetRespPayload().(TAR)
    if !ok { return nil, 0, errors.New("wrong response type") }
    payload, state, err := s.cleanStateMachine.responsePresentSecond(aInput, aPayload, aResponse)
    return payload, state, err
}
func (s *RawTestStateMachine) responseFutureSecond(element slot.SlotElementHelper, ar adapter.IAdapterResponse) (interface{}, statemachine.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    aResponse, ok := ar.GetRespPayload().(TAR)
    if !ok { return nil, 0, errors.New("wrong response type") }
    payload, state, err := s.cleanStateMachine.responseFutureSecond(aInput, aPayload, aResponse)
    return payload, state, err
}
func (s *RawTestStateMachine) responsePastSecond(element slot.SlotElementHelper, ar adapter.IAdapterResponse) (interface{}, statemachine.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    aResponse, ok := ar.GetRespPayload().(TAR)
    if !ok { return nil, 0, errors.New("wrong response type") }
    payload, state, err := s.cleanStateMachine.responsePastSecond(aInput, aPayload, aResponse)
    return payload, state, err
}
func (s *RawTestStateMachine) errorResponsePresentSecond(element slot.SlotElementHelper, ar adapter.IAdapterResponse, err error) (interface{}, statemachine.ElementState) {
    payload, state := s.cleanStateMachine.errorResponsePresentSecond(element.GetInputEvent(), element.GetPayload(), ar, err)
    return payload, state
}
func (s *RawTestStateMachine) errorResponseFutureSecond(element slot.SlotElementHelper, ar adapter.IAdapterResponse, err error) (interface{}, statemachine.ElementState) {
    payload, state := s.cleanStateMachine.errorResponseFutureSecond(element.GetInputEvent(), element.GetPayload(), ar, err)
    return payload, state
}
func (s *RawTestStateMachine) errorResponsePastSecond(element slot.SlotElementHelper, ar adapter.IAdapterResponse, err error) (interface{}, statemachine.ElementState) {
    payload, state := s.cleanStateMachine.errorResponsePastSecond(element.GetInputEvent(), element.GetPayload(), ar, err)
    return payload, state
}









