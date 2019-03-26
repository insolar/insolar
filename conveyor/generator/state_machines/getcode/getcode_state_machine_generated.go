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
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/ledger/artifactmanager"

	"errors"
)

type BaseGetCodeStateMachine struct{}

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
			Migration: m.migrateFromPresentFirst,

			ErrorState: m.errorPresentFirst,
		},
		common.State{
			Migration: m.migrateFromPresentSecond,

			ErrorState: m.errorPresentSecond,
		})

	// present state machine
	x[1] = append(x[1], common.State{
		Transition: m.initPresentHandler,
		ErrorState: m.errorPresentInit,
	},
		common.State{
			Migration:            m.migrateFromPresentFirst,
			Transition:           m.transitPresentFirst,
			AdapterResponse:      m.responsePresentFirst,
			ErrorState:           m.errorPresentFirst,
			AdapterResponseError: m.errorResponsePresentFirst,
		},
		common.State{
			Migration:            m.migrateFromPresentSecond,
			Transition:           m.transitPresentSecond,
			AdapterResponse:      m.responsePresentSecond,
			ErrorState:           m.errorPresentSecond,
			AdapterResponseError: m.errorResponsePresentSecond,
		})

	// past state machine
	x[2] = append(x[2], common.State{
		Transition: m.initPresentHandler,
		ErrorState: m.errorPresentInit,
	},
		common.State{
			Transition: m.transitPresentFirst,

			ErrorState: m.errorPresentFirst,
		},
		common.State{
			Transition: m.transitPresentSecond,

			ErrorState: m.errorPresentSecond,
		})

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
	if !ok {
		return nil, 0, errors.New("wrong input event type")
	}
	payload, state, err := s.cleanStateMachine.initPresentHandler(aInput, element.GetPayload())
	return payload, state, err
}
func (s *RawGetCodeStateMachine) initFutureHandler(element slot.SlotElementHelper) (interface{}, fsm.ElementState, error) {
	aInput, ok := element.GetInputEvent().(Event)
	if !ok {
		return nil, 0, errors.New("wrong input event type")
	}
	payload, state, err := s.cleanStateMachine.initFutureHandler(aInput, element.GetPayload(), element)
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

func (s *RawGetCodeStateMachine) transitPresentFirst(element slot.SlotElementHelper) (interface{}, fsm.ElementState, error) {
	aInput, ok := element.GetInputEvent().(Event)
	if !ok {
		return nil, 0, errors.New("wrong input event type")
	}
	aPayload, ok := element.GetPayload().(insolar.ConveyorPendingMessage)
	if !ok {
		return nil, 0, errors.New("wrong payload type")
	}
	payload, state, err := s.cleanStateMachine.transitPresentFirst(aInput, aPayload, element)
	return payload, state, err
}

func (s *RawGetCodeStateMachine) migrateFromPresentFirst(element slot.SlotElementHelper) (interface{}, fsm.ElementState, error) {
	aInput, ok := element.GetInputEvent().(Event)
	if !ok {
		return nil, 0, errors.New("wrong input event type")
	}
	aPayload, ok := element.GetPayload().(*Payload)
	if !ok {
		return nil, 0, errors.New("wrong payload type")
	}
	payload, state, err := s.cleanStateMachine.migrateFromPresentFirst(aInput, aPayload)
	return payload, state, err
}

func (s *RawGetCodeStateMachine) errorPresentFirst(element slot.SlotElementHelper, err error) (interface{}, fsm.ElementState) {
	payload, state := s.cleanStateMachine.errorPresentFirst(element.GetInputEvent(), element.GetPayload(), err)
	return payload, state
}

func (s *RawGetCodeStateMachine) responsePresentFirst(element slot.SlotElementHelper, ar iadapter.Response) (interface{}, fsm.ElementState, error) {
	aInput, ok := element.GetInputEvent().(Event)
	if !ok {
		return nil, 0, errors.New("wrong input event type")
	}
	aResponse, ok := ar.GetRespPayload().(artifactmanager.GetCodeResp)
	if !ok {
		return nil, 0, errors.New("wrong response type")
	}
	payload, state, err := s.cleanStateMachine.responsePresentFirst(aInput, aResponse, element)
	return payload, state, err
}

func (s *RawGetCodeStateMachine) errorResponsePresentFirst(element slot.SlotElementHelper, ar iadapter.Response, err error) (interface{}, fsm.ElementState) {
	payload, state := s.cleanStateMachine.errorResponsePresentFirst(element.GetInputEvent(), element.GetPayload(), ar, err)
	return payload, state
}

func (s *RawGetCodeStateMachine) transitPresentSecond(element slot.SlotElementHelper) (interface{}, fsm.ElementState, error) {
	aInput, ok := element.GetInputEvent().(Event)
	if !ok {
		return nil, 0, errors.New("wrong input event type")
	}
	aPayload, ok := element.GetPayload().(*Payload)
	if !ok {
		return nil, 0, errors.New("wrong payload type")
	}
	payload, state, err := s.cleanStateMachine.transitPresentSecond(aInput, aPayload)
	return payload, state, err
}

func (s *RawGetCodeStateMachine) migrateFromPresentSecond(element slot.SlotElementHelper) (interface{}, fsm.ElementState, error) {
	aInput, ok := element.GetInputEvent().(Event)
	if !ok {
		return nil, 0, errors.New("wrong input event type")
	}
	aPayload, ok := element.GetPayload().(*Payload)
	if !ok {
		return nil, 0, errors.New("wrong payload type")
	}
	payload, state, err := s.cleanStateMachine.migrateFromPresentSecond(aInput, aPayload)
	return payload, state, err
}

func (s *RawGetCodeStateMachine) errorPresentSecond(element slot.SlotElementHelper, err error) (interface{}, fsm.ElementState) {
	payload, state := s.cleanStateMachine.errorPresentSecond(element.GetInputEvent(), element.GetPayload(), err)
	return payload, state
}

func (s *RawGetCodeStateMachine) responsePresentSecond(element slot.SlotElementHelper, ar iadapter.Response) (interface{}, fsm.ElementState, error) {
	aInput, ok := element.GetInputEvent().(Event)
	if !ok {
		return nil, 0, errors.New("wrong input event type")
	}
	aResponse := ar.GetRespPayload()
	payload, state, err := s.cleanStateMachine.responsePresentSecond(aInput, aResponse, element)
	return payload, state, err
}

func (s *RawGetCodeStateMachine) errorResponsePresentSecond(element slot.SlotElementHelper, ar iadapter.Response, err error) (interface{}, fsm.ElementState) {
	payload, state := s.cleanStateMachine.errorResponsePresentSecond(element.GetInputEvent(), element.GetPayload(), ar, err)
	return payload, state
}
