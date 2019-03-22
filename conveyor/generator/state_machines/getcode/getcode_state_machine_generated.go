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
	"github.com/insolar/insolar/conveyor/generator/common"
    "github.com/insolar/insolar/conveyor/interfaces/fsm"
	"github.com/insolar/insolar/conveyor/interfaces/iadapter"
	"github.com/insolar/insolar/conveyor/interfaces/slot"
	"github.com/insolar/insolar/conveyor/interfaces/statemachine"

	"errors"
)

type BaseGetCodeStateMachine struct {}

func (*BaseGetCodeStateMachine) GetTypeID() fsm.ID {
    return 1
}

func (*BaseGetCodeStateMachine) stateFirst() fsm.StateID {
    return 1
}
func (*BaseGetCodeStateMachine) stateSecond() fsm.StateID {
    return 2
}

type RawGetCodeStateMachine struct {
    cleanStateMachine GetCodeStateMachine
}

func RawGetCodeStateMachineFactory() [3]statemachine.StateMachine {
    m := RawGetCodeStateMachine{
        cleanStateMachine: &CleanGetCodeStateMachine{},
    }

    var x = [3][]common.State{}
    // future state machine
    x[0] = append(x[0], common.State{
        Transition: m.initFutureHandler,
        ErrorState: m.errorFutureInit,
    },
    common.State{
        Migration: m.migrateFromFutureFirst,
        Transition: m.transitFutureFirst,
        AdapterResponse: m.responseFutureFirst,
        ErrorState: m.errorFutureFirst,
        AdapterResponseError: m.errorResponseFutureFirst,
    },
    common.State{
        Migration: m.migrateFromFutureSecond,
        Transition: m.transitFutureSecond,
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
        Migration: m.migrateFromPresentFirst,
        Transition: m.transitPresentFirst,
        AdapterResponse: m.responsePresentFirst,
        ErrorState: m.errorPresentFirst,
        AdapterResponseError: m.errorResponsePresentFirst,
    },
    common.State{
        Migration: m.migrateFromPresentSecond,
        Transition: m.transitPresentSecond,
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
        Transition: m.transitPastSecond,
        AdapterResponse: m.responsePastSecond,
        ErrorState: m.errorPastSecond,
        AdapterResponseError: m.errorResponsePastSecond,
    },)


    smFuture := common.StateMachine{
        ID:     m.cleanStateMachine.(GetCodeStateMachine).GetTypeID(),
        States: x[0],
    }

    smPresent := common.StateMachine{
        ID:     m.cleanStateMachine.(GetCodeStateMachine).GetTypeID(),
        States: x[1],
    }

    smPast := common.StateMachine{
        ID:     m.cleanStateMachine.(GetCodeStateMachine).GetTypeID(),
        States: x[2],
    }

    return [3]statemachine.StateMachine{
        &smFuture, &smPresent, &smPast,
    }
}

func (s *RawGetCodeStateMachine) initPresentHandler(element slot.SlotElementHelper) (interface{}, fsm.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    payload, state, err := s.cleanStateMachine.initPresentHandler(aInput, element.GetPayload())
    return payload, state, err
}
func (s *RawGetCodeStateMachine) initFutureHandler(element slot.SlotElementHelper) (interface{}, fsm.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    payload, state, err := s.cleanStateMachine.initFutureHandler(aInput, element.GetPayload())
    return payload, state, err
}
func (s *RawGetCodeStateMachine) initPastHandler(element slot.SlotElementHelper) (interface{}, fsm.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    payload, state, err := s.cleanStateMachine.initPastHandler(aInput, element.GetPayload())
    return payload, state, err
}
func (s *RawGetCodeStateMachine) errorPresentInit(element slot.SlotElementHelper, err error) (interface{}, fsm.ElementState) {
    payload, state := s.cleanStateMachine.errorPresentInit(element.GetInputEvent(), element.GetPayload(), err)
    return payload, state
}
func (s *RawGetCodeStateMachine) errorFutureInit(element slot.SlotElementHelper, err error) (interface{}, fsm.ElementState) {
    payload, state := s.cleanStateMachine.errorFutureInit(element.GetInputEvent(), element.GetPayload(), err)
    return payload, state
}
func (s *RawGetCodeStateMachine) errorPastInit(element slot.SlotElementHelper, err error) (interface{}, fsm.ElementState) {
    payload, state := s.cleanStateMachine.errorPastInit(element.GetInputEvent(), element.GetPayload(), err)
    return payload, state
}


func (s *RawGetCodeStateMachine) transitPresentFirst(element slot.SlotElementHelper) (interface{}, fsm.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    payload, state, err := s.cleanStateMachine.transitPresentFirst(aInput, aPayload)
    return payload, state, err
}
func (s *RawGetCodeStateMachine) transitFutureFirst(element slot.SlotElementHelper) (interface{}, fsm.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    payload, state, err := s.cleanStateMachine.transitFutureFirst(aInput, aPayload)
    return payload, state, err
}
func (s *RawGetCodeStateMachine) transitPastFirst(element slot.SlotElementHelper) (interface{}, fsm.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    payload, state, err := s.cleanStateMachine.transitPastFirst(aInput, aPayload)
    return payload, state, err
}
func (s *RawGetCodeStateMachine) migrateFromPresentFirst(element slot.SlotElementHelper) (interface{}, fsm.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    payload, state, err := s.cleanStateMachine.migrateFromPresentFirst(aInput, aPayload)
    return payload, state, err
}
func (s *RawGetCodeStateMachine) migrateFromFutureFirst(element slot.SlotElementHelper) (interface{}, fsm.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    payload, state, err := s.cleanStateMachine.migrateFromFutureFirst(aInput, aPayload)
    return payload, state, err
}
func (s *RawGetCodeStateMachine) errorPresentFirst(element slot.SlotElementHelper, err error) (interface{}, fsm.ElementState) {
    payload, state := s.cleanStateMachine.errorPresentFirst(element.GetInputEvent(), element.GetPayload(), err)
    return payload, state
}
func (s *RawGetCodeStateMachine) errorFutureFirst(element slot.SlotElementHelper, err error) (interface{}, fsm.ElementState) {
    payload, state := s.cleanStateMachine.errorFutureFirst(element.GetInputEvent(), element.GetPayload(), err)
    return payload, state
}
func (s *RawGetCodeStateMachine) errorPastFirst(element slot.SlotElementHelper, err error) (interface{}, fsm.ElementState) {
    payload, state := s.cleanStateMachine.errorPastFirst(element.GetInputEvent(), element.GetPayload(), err)
    return payload, state
}
func (s *RawGetCodeStateMachine) responsePresentFirst(element slot.SlotElementHelper, ar iadapter.Response) (interface{}, fsm.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    aResponse, ok := ar.GetRespPayload().(TAR)
    if !ok { return nil, 0, errors.New("wrong response type") }
    payload, state, err := s.cleanStateMachine.responsePresentFirst(aInput, aPayload, aResponse)
    return payload, state, err
}
func (s *RawGetCodeStateMachine) responseFutureFirst(element slot.SlotElementHelper, ar iadapter.Response) (interface{}, fsm.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    aResponse, ok := ar.GetRespPayload().(TAR)
    if !ok { return nil, 0, errors.New("wrong response type") }
    payload, state, err := s.cleanStateMachine.responseFutureFirst(aInput, aPayload, aResponse)
    return payload, state, err
}
func (s *RawGetCodeStateMachine) responsePastFirst(element slot.SlotElementHelper, ar iadapter.Response) (interface{}, fsm.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    aResponse, ok := ar.GetRespPayload().(TAR)
    if !ok { return nil, 0, errors.New("wrong response type") }
    payload, state, err := s.cleanStateMachine.responsePastFirst(aInput, aPayload, aResponse)
    return payload, state, err
}
func (s *RawGetCodeStateMachine) errorResponsePresentFirst(element slot.SlotElementHelper, ar iadapter.Response, err error) (interface{}, fsm.ElementState) {
    payload, state := s.cleanStateMachine.errorResponsePresentFirst(element.GetInputEvent(), element.GetPayload(), ar, err)
    return payload, state
}
func (s *RawGetCodeStateMachine) errorResponseFutureFirst(element slot.SlotElementHelper, ar iadapter.Response, err error) (interface{}, fsm.ElementState) {
    payload, state := s.cleanStateMachine.errorResponseFutureFirst(element.GetInputEvent(), element.GetPayload(), ar, err)
    return payload, state
}
func (s *RawGetCodeStateMachine) errorResponsePastFirst(element slot.SlotElementHelper, ar iadapter.Response, err error) (interface{}, fsm.ElementState) {
    payload, state := s.cleanStateMachine.errorResponsePastFirst(element.GetInputEvent(), element.GetPayload(), ar, err)
    return payload, state
}

func (s *RawGetCodeStateMachine) transitPresentSecond(element slot.SlotElementHelper) (interface{}, fsm.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    payload, state, err := s.cleanStateMachine.transitPresentSecond(aInput, aPayload)
    return payload, state, err
}
func (s *RawGetCodeStateMachine) transitFutureSecond(element slot.SlotElementHelper) (interface{}, fsm.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    payload, state, err := s.cleanStateMachine.transitFutureSecond(aInput, aPayload)
    return payload, state, err
}
func (s *RawGetCodeStateMachine) transitPastSecond(element slot.SlotElementHelper) (interface{}, fsm.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    payload, state, err := s.cleanStateMachine.transitPastSecond(aInput, aPayload)
    return payload, state, err
}
func (s *RawGetCodeStateMachine) migrateFromPresentSecond(element slot.SlotElementHelper) (interface{}, fsm.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    payload, state, err := s.cleanStateMachine.migrateFromPresentSecond(aInput, aPayload)
    return payload, state, err
}
func (s *RawGetCodeStateMachine) migrateFromFutureSecond(element slot.SlotElementHelper) (interface{}, fsm.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    payload, state, err := s.cleanStateMachine.migrateFromFutureSecond(aInput, aPayload)
    return payload, state, err
}
func (s *RawGetCodeStateMachine) errorPresentSecond(element slot.SlotElementHelper, err error) (interface{}, fsm.ElementState) {
    payload, state := s.cleanStateMachine.errorPresentSecond(element.GetInputEvent(), element.GetPayload(), err)
    return payload, state
}
func (s *RawGetCodeStateMachine) errorFutureSecond(element slot.SlotElementHelper, err error) (interface{}, fsm.ElementState) {
    payload, state := s.cleanStateMachine.errorFutureSecond(element.GetInputEvent(), element.GetPayload(), err)
    return payload, state
}
func (s *RawGetCodeStateMachine) errorPastSecond(element slot.SlotElementHelper, err error) (interface{}, fsm.ElementState) {
    payload, state := s.cleanStateMachine.errorPastSecond(element.GetInputEvent(), element.GetPayload(), err)
    return payload, state
}
func (s *RawGetCodeStateMachine) responsePresentSecond(element slot.SlotElementHelper, ar iadapter.Response) (interface{}, fsm.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    aResponse, ok := ar.GetRespPayload().(TAR)
    if !ok { return nil, 0, errors.New("wrong response type") }
    payload, state, err := s.cleanStateMachine.responsePresentSecond(aInput, aPayload, aResponse)
    return payload, state, err
}
func (s *RawGetCodeStateMachine) responseFutureSecond(element slot.SlotElementHelper, ar iadapter.Response) (interface{}, fsm.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    aResponse, ok := ar.GetRespPayload().(TAR)
    if !ok { return nil, 0, errors.New("wrong response type") }
    payload, state, err := s.cleanStateMachine.responseFutureSecond(aInput, aPayload, aResponse)
    return payload, state, err
}
func (s *RawGetCodeStateMachine) responsePastSecond(element slot.SlotElementHelper, ar iadapter.Response) (interface{}, fsm.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    aResponse, ok := ar.GetRespPayload().(TAR)
    if !ok { return nil, 0, errors.New("wrong response type") }
    payload, state, err := s.cleanStateMachine.responsePastSecond(aInput, aPayload, aResponse)
    return payload, state, err
}
func (s *RawGetCodeStateMachine) errorResponsePresentSecond(element slot.SlotElementHelper, ar iadapter.Response, err error) (interface{}, fsm.ElementState) {
    payload, state := s.cleanStateMachine.errorResponsePresentSecond(element.GetInputEvent(), element.GetPayload(), ar, err)
    return payload, state
}
func (s *RawGetCodeStateMachine) errorResponseFutureSecond(element slot.SlotElementHelper, ar iadapter.Response, err error) (interface{}, fsm.ElementState) {
    payload, state := s.cleanStateMachine.errorResponseFutureSecond(element.GetInputEvent(), element.GetPayload(), ar, err)
    return payload, state
}
func (s *RawGetCodeStateMachine) errorResponsePastSecond(element slot.SlotElementHelper, ar iadapter.Response, err error) (interface{}, fsm.ElementState) {
    payload, state := s.cleanStateMachine.errorResponsePastSecond(element.GetInputEvent(), element.GetPayload(), ar, err)
    return payload, state
}









