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
	"errors"

	"github.com/insolar/insolar/conveyor/fsm"
	"github.com/insolar/insolar/conveyor/statemachine"
)

type BaseTestStateMachine struct {}

func (*BaseTestStateMachine) GetTypeID() fsm.ID {
    return 1
}

func (*BaseTestStateMachine) stateFirst() fsm.StateID {
    return 1
}
func (*BaseTestStateMachine) stateSecond() fsm.StateID {
    return 2
}

type RawTestStateMachine struct {
    cleanStateMachine TestStateMachine
}

func RawTestStateMachineFactory() [3]*statemachine.StateMachine {
    m := RawTestStateMachine{
        cleanStateMachine: &CleanTestStateMachine{},
    }

    var x = [3][]statemachine.State{}
    // future state machine
    x[0] = append(x[0], statemachine.State{
        Transition: m.initFutureHandler,
        ErrorState: m.errorFutureInit,
    },
    statemachine.State{
        Migration: m.migrateFromFutureFirst,
        Transition: m.transitFutureFirst,
        AdapterResponse: m.responseFutureFirst,
        ErrorState: m.errorFutureFirst,
        AdapterResponseError: m.errorResponseFutureFirst,
    },
    statemachine.State{
        Migration: m.migrateFromFutureSecond,
        Transition: m.transitFutureSecond,
        AdapterResponse: m.responseFutureSecond,
        ErrorState: m.errorFutureSecond,
        AdapterResponseError: m.errorResponseFutureSecond,
    },)

    // present state machine
    x[1] = append(x[1], statemachine.State{
        Transition: m.initPresentHandler,
        ErrorState: m.errorPresentInit,
    },
    statemachine.State{
        Migration: m.migrateFromPresentFirst,
        Transition: m.transitPresentFirst,
        AdapterResponse: m.responsePresentFirst,
        ErrorState: m.errorPresentFirst,
        AdapterResponseError: m.errorResponsePresentFirst,
    },
    statemachine.State{
        Migration: m.migrateFromPresentSecond,
        Transition: m.transitPresentSecond,
        AdapterResponse: m.responsePresentSecond,
        ErrorState: m.errorPresentSecond,
        AdapterResponseError: m.errorResponsePresentSecond,
    },)

    // past state machine
    x[2] = append(x[2], statemachine.State{
        Transition: m.initPastHandler,
        ErrorState: m.errorPastInit,
    },
    statemachine.State{
        Transition: m.transitPastFirst,
        AdapterResponse: m.responsePastFirst,
        ErrorState: m.errorPastFirst,
        AdapterResponseError: m.errorResponsePastFirst,
    },
    statemachine.State{
        Transition: m.transitPastSecond,
        AdapterResponse: m.responsePastSecond,
        ErrorState: m.errorPastSecond,
        AdapterResponseError: m.errorResponsePastSecond,
    },)


    smFuture := &statemachine.StateMachine{
        ID:     m.cleanStateMachine.(TestStateMachine).GetTypeID(),
        States: x[0],
    }

    smPresent := &statemachine.StateMachine{
        ID:     m.cleanStateMachine.(TestStateMachine).GetTypeID(),
        States: x[1],
    }

    smPast := &statemachine.StateMachine{
        ID:     m.cleanStateMachine.(TestStateMachine).GetTypeID(),
        States: x[2],
    }

    return [3]*statemachine.StateMachine{
        smFuture, smPresent, smPast,
    }
}

func (s *RawTestStateMachine) initPresentHandler(element fsm.SlotElementHelper) (interface{}, fsm.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    payload, state, err := s.cleanStateMachine.initPresentHandler(aInput, element.GetPayload())
    return payload, state, err
}
func (s *RawTestStateMachine) initFutureHandler(element fsm.SlotElementHelper) (interface{}, fsm.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    payload, state, err := s.cleanStateMachine.initFutureHandler(aInput, element.GetPayload())
    return payload, state, err
}
func (s *RawTestStateMachine) initPastHandler(element fsm.SlotElementHelper) (interface{}, fsm.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    payload, state, err := s.cleanStateMachine.initPastHandler(aInput, element.GetPayload())
    return payload, state, err
}
func (s *RawTestStateMachine) errorPresentInit(element fsm.SlotElementHelper, err error) (interface{}, fsm.ElementState) {
    payload, state := s.cleanStateMachine.errorPresentInit(element.GetInputEvent(), element.GetPayload(), err)
    return payload, state
}
func (s *RawTestStateMachine) errorFutureInit(element fsm.SlotElementHelper, err error) (interface{}, fsm.ElementState) {
    payload, state := s.cleanStateMachine.errorFutureInit(element.GetInputEvent(), element.GetPayload(), err)
    return payload, state
}
func (s *RawTestStateMachine) errorPastInit(element fsm.SlotElementHelper, err error) (interface{}, fsm.ElementState) {
    payload, state := s.cleanStateMachine.errorPastInit(element.GetInputEvent(), element.GetPayload(), err)
    return payload, state
}


func (s *RawTestStateMachine) transitPresentFirst(element fsm.SlotElementHelper) (interface{}, fsm.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    payload, state, err := s.cleanStateMachine.transitPresentFirst(aInput, aPayload)
    return payload, state, err
}
func (s *RawTestStateMachine) transitFutureFirst(element fsm.SlotElementHelper) (interface{}, fsm.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    payload, state, err := s.cleanStateMachine.transitFutureFirst(aInput, aPayload)
    return payload, state, err
}
func (s *RawTestStateMachine) transitPastFirst(element fsm.SlotElementHelper) (interface{}, fsm.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    payload, state, err := s.cleanStateMachine.transitPastFirst(aInput, aPayload)
    return payload, state, err
}
func (s *RawTestStateMachine) migrateFromPresentFirst(element fsm.SlotElementHelper) (interface{}, fsm.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    payload, state, err := s.cleanStateMachine.migrateFromPresentFirst(aInput, aPayload)
    return payload, state, err
}
func (s *RawTestStateMachine) migrateFromFutureFirst(element fsm.SlotElementHelper) (interface{}, fsm.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    payload, state, err := s.cleanStateMachine.migrateFromFutureFirst(aInput, aPayload)
    return payload, state, err
}
func (s *RawTestStateMachine) errorPresentFirst(element fsm.SlotElementHelper, err error) (interface{}, fsm.ElementState) {
    payload, state := s.cleanStateMachine.errorPresentFirst(element.GetInputEvent(), element.GetPayload(), err)
    return payload, state
}
func (s *RawTestStateMachine) errorFutureFirst(element fsm.SlotElementHelper, err error) (interface{}, fsm.ElementState) {
    payload, state := s.cleanStateMachine.errorFutureFirst(element.GetInputEvent(), element.GetPayload(), err)
    return payload, state
}
func (s *RawTestStateMachine) errorPastFirst(element fsm.SlotElementHelper, err error) (interface{}, fsm.ElementState) {
    payload, state := s.cleanStateMachine.errorPastFirst(element.GetInputEvent(), element.GetPayload(), err)
    return payload, state
}
func (s *RawTestStateMachine) responsePresentFirst(element fsm.SlotElementHelper, ar interface{}) (interface{}, fsm.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    aResponse, ok := ar.(TAR)
    if !ok { return nil, 0, errors.New("wrong response type") }
    payload, state, err := s.cleanStateMachine.responsePresentFirst(aInput, aPayload, aResponse)
    return payload, state, err
}
func (s *RawTestStateMachine) responseFutureFirst(element fsm.SlotElementHelper, ar interface{}) (interface{}, fsm.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    aResponse, ok := ar.(TAR)
    if !ok { return nil, 0, errors.New("wrong response type") }
    payload, state, err := s.cleanStateMachine.responseFutureFirst(aInput, aPayload, aResponse)
    return payload, state, err
}
func (s *RawTestStateMachine) responsePastFirst(element fsm.SlotElementHelper, ar interface{}) (interface{}, fsm.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    aResponse, ok := ar.(TAR)
    if !ok { return nil, 0, errors.New("wrong response type") }
    payload, state, err := s.cleanStateMachine.responsePastFirst(aInput, aPayload, aResponse)
    return payload, state, err
}
func (s *RawTestStateMachine) errorResponsePresentFirst(element fsm.SlotElementHelper, ar interface{}, err error) (interface{}, fsm.ElementState) {
    payload, state := s.cleanStateMachine.errorResponsePresentFirst(element.GetInputEvent(), element.GetPayload(), ar, err)
    return payload, state
}
func (s *RawTestStateMachine) errorResponseFutureFirst(element fsm.SlotElementHelper, ar interface{}, err error) (interface{}, fsm.ElementState) {
    payload, state := s.cleanStateMachine.errorResponseFutureFirst(element.GetInputEvent(), element.GetPayload(), ar, err)
    return payload, state
}
func (s *RawTestStateMachine) errorResponsePastFirst(element fsm.SlotElementHelper, ar interface{}, err error) (interface{}, fsm.ElementState) {
    payload, state := s.cleanStateMachine.errorResponsePastFirst(element.GetInputEvent(), element.GetPayload(), ar, err)
    return payload, state
}

func (s *RawTestStateMachine) transitPresentSecond(element fsm.SlotElementHelper) (interface{}, fsm.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    payload, state, err := s.cleanStateMachine.transitPresentSecond(aInput, aPayload)
    return payload, state, err
}
func (s *RawTestStateMachine) transitFutureSecond(element fsm.SlotElementHelper) (interface{}, fsm.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    payload, state, err := s.cleanStateMachine.transitFutureSecond(aInput, aPayload)
    return payload, state, err
}
func (s *RawTestStateMachine) transitPastSecond(element fsm.SlotElementHelper) (interface{}, fsm.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    payload, state, err := s.cleanStateMachine.transitPastSecond(aInput, aPayload)
    return payload, state, err
}
func (s *RawTestStateMachine) migrateFromPresentSecond(element fsm.SlotElementHelper) (interface{}, fsm.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    payload, state, err := s.cleanStateMachine.migrateFromPresentSecond(aInput, aPayload)
    return payload, state, err
}
func (s *RawTestStateMachine) migrateFromFutureSecond(element fsm.SlotElementHelper) (interface{}, fsm.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    payload, state, err := s.cleanStateMachine.migrateFromFutureSecond(aInput, aPayload)
    return payload, state, err
}
func (s *RawTestStateMachine) errorPresentSecond(element fsm.SlotElementHelper, err error) (interface{}, fsm.ElementState) {
    payload, state := s.cleanStateMachine.errorPresentSecond(element.GetInputEvent(), element.GetPayload(), err)
    return payload, state
}
func (s *RawTestStateMachine) errorFutureSecond(element fsm.SlotElementHelper, err error) (interface{}, fsm.ElementState) {
    payload, state := s.cleanStateMachine.errorFutureSecond(element.GetInputEvent(), element.GetPayload(), err)
    return payload, state
}
func (s *RawTestStateMachine) errorPastSecond(element fsm.SlotElementHelper, err error) (interface{}, fsm.ElementState) {
    payload, state := s.cleanStateMachine.errorPastSecond(element.GetInputEvent(), element.GetPayload(), err)
    return payload, state
}
func (s *RawTestStateMachine) responsePresentSecond(element fsm.SlotElementHelper, ar interface{}) (interface{}, fsm.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    aResponse, ok := ar.(TAR)
    if !ok { return nil, 0, errors.New("wrong response type") }
    payload, state, err := s.cleanStateMachine.responsePresentSecond(aInput, aPayload, aResponse)
    return payload, state, err
}
func (s *RawTestStateMachine) responseFutureSecond(element fsm.SlotElementHelper, ar interface{}) (interface{}, fsm.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    aResponse, ok := ar.(TAR)
    if !ok { return nil, 0, errors.New("wrong response type") }
    payload, state, err := s.cleanStateMachine.responseFutureSecond(aInput, aPayload, aResponse)
    return payload, state, err
}
func (s *RawTestStateMachine) responsePastSecond(element fsm.SlotElementHelper, ar interface{}) (interface{}, fsm.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    aResponse, ok := ar.(TAR)
    if !ok { return nil, 0, errors.New("wrong response type") }
    payload, state, err := s.cleanStateMachine.responsePastSecond(aInput, aPayload, aResponse)
    return payload, state, err
}
func (s *RawTestStateMachine) errorResponsePresentSecond(element fsm.SlotElementHelper, ar interface{}, err error) (interface{}, fsm.ElementState) {
    payload, state := s.cleanStateMachine.errorResponsePresentSecond(element.GetInputEvent(), element.GetPayload(), ar, err)
    return payload, state
}
func (s *RawTestStateMachine) errorResponseFutureSecond(element fsm.SlotElementHelper, ar interface{}, err error) (interface{}, fsm.ElementState) {
    payload, state := s.cleanStateMachine.errorResponseFutureSecond(element.GetInputEvent(), element.GetPayload(), ar, err)
    return payload, state
}
func (s *RawTestStateMachine) errorResponsePastSecond(element fsm.SlotElementHelper, ar interface{}, err error) (interface{}, fsm.ElementState) {
    payload, state := s.cleanStateMachine.errorResponsePastSecond(element.GetInputEvent(), element.GetPayload(), ar, err)
    return payload, state
}









