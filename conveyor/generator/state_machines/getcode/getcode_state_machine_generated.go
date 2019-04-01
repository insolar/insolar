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
	"github.com/insolar/insolar/conveyor/fsm"
	"github.com/insolar/insolar/conveyor/statemachine"
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
func (*BaseGetCodeStateMachine) stateThird() fsm.StateID {
	return 3
}
func (*BaseGetCodeStateMachine) stateFourth() fsm.StateID {
	return 4
}

type RawGetCodeStateMachine struct {
	cleanStateMachine GetCodeStateMachine
}

func RawGetCodeStateMachineFactory() [3]*statemachine.StateMachine {
	m := RawGetCodeStateMachine{
		cleanStateMachine: &CleanGetCodeStateMachine{},
	}

	var x = [3][]statemachine.State{}
	// future state machine
	x[0] = append(x[0], statemachine.State{
		Transition: m.initFutureHandler,
		ErrorState: m.errorFutureInit,
	},
		statemachine.State{
			Migration: m.migrateFromPresentFirst,

			ErrorState: m.errorPresentFirst,
		},
		statemachine.State{
			Migration: m.migrateFromPresentSecond,

			ErrorState: m.errorPresentSecond,
		},
		statemachine.State{
			Migration: m.migrateFromPresentThird,

			ErrorState: m.errorPresentThird,
		},
		statemachine.State{
			Migration: m.migrateFromPresentFourth,

			ErrorState: m.errorPresentFourth,
		})

	// present state machine
	x[1] = append(x[1], statemachine.State{
		Transition: m.initPresentHandler,
		ErrorState: m.errorPresentInit,
	},
		statemachine.State{
			Migration:            m.migrateFromPresentFirst,
			Transition:           m.transitPresentFirst,
			AdapterResponse:      m.responsePresentFirst,
			ErrorState:           m.errorPresentFirst,
			AdapterResponseError: m.errorResponsePresentFirst,
		},
		statemachine.State{
			Migration:            m.migrateFromPresentSecond,
			Transition:           m.transitPresentSecond,
			AdapterResponse:      m.responsePresentSecond,
			ErrorState:           m.errorPresentSecond,
			AdapterResponseError: m.errorResponsePresentSecond,
		},
		statemachine.State{
			Migration:            m.migrateFromPresentThird,
			Transition:           m.transitPresentThird,
			AdapterResponse:      m.responsePresentThird,
			ErrorState:           m.errorPresentThird,
			AdapterResponseError: m.errorResponsePresentThird,
		},
		statemachine.State{
			Migration:            m.migrateFromPresentFourth,
			Transition:           m.transitPresentFourth,
			AdapterResponse:      m.responsePresentFourth,
			ErrorState:           m.errorPresentFourth,
			AdapterResponseError: m.errorResponsePresentFourth,
		})

	// past state machine
	x[2] = append(x[2], statemachine.State{
		Transition: m.initPresentHandler,
		ErrorState: m.errorPresentInit,
	},
		statemachine.State{
			Transition:      m.transitPresentFirst,
			AdapterResponse: m.responsePresentFirst,

			ErrorState: m.errorPresentFirst,
		},
		statemachine.State{
			Transition:      m.transitPresentSecond,
			AdapterResponse: m.responsePresentSecond,

			ErrorState: m.errorPresentSecond,
		},
		statemachine.State{
			Transition:      m.transitPresentThird,
			AdapterResponse: m.responsePresentThird,

			ErrorState: m.errorPresentThird,
		},
		statemachine.State{
			Transition:      m.transitPresentFourth,
			AdapterResponse: m.responsePresentFourth,

			ErrorState: m.errorPresentFourth,
		})

	smFuture := statemachine.StateMachine{
		ID:     m.cleanStateMachine.(GetCodeStateMachine).GetTypeID(),
		States: x[0],
	}

	smPresent := statemachine.StateMachine{
		ID:     m.cleanStateMachine.(GetCodeStateMachine).GetTypeID(),
		States: x[1],
	}

	smPast := statemachine.StateMachine{
		ID:     m.cleanStateMachine.(GetCodeStateMachine).GetTypeID(),
		States: x[2],
	}

	return [3]*statemachine.StateMachine{
		&smFuture, &smPresent, &smPast,
	}
}

func (s *RawGetCodeStateMachine) initPresentHandler(element fsm.SlotElementHelper) (interface{}, fsm.ElementState, error) {
	aInput, ok := element.GetInputEvent().(Event)
	if !ok {
		return nil, 0, errors.New("wrong input event type")
	}
	payload, state, err := s.cleanStateMachine.initPresentHandler(aInput, element.GetPayload())
	return payload, state, err
}
func (s *RawGetCodeStateMachine) initFutureHandler(element fsm.SlotElementHelper) (interface{}, fsm.ElementState, error) {
	aInput, ok := element.GetInputEvent().(Event)
	if !ok {
		return nil, 0, errors.New("wrong input event type")
	}
	payload, state, err := s.cleanStateMachine.initFutureHandler(aInput, element.GetPayload(), element)
	return payload, state, err
}

func (s *RawGetCodeStateMachine) errorPresentInit(element fsm.SlotElementHelper, err error) (interface{}, fsm.ElementState) {
	payload, state := s.cleanStateMachine.errorPresentInit(element.GetInputEvent(), element.GetPayload(), err)
	return payload, state
}
func (s *RawGetCodeStateMachine) errorFutureInit(element fsm.SlotElementHelper, err error) (interface{}, fsm.ElementState) {
	payload, state := s.cleanStateMachine.errorFutureInit(element.GetInputEvent(), element.GetPayload(), err)
	return payload, state
}

func (s *RawGetCodeStateMachine) transitPresentFirst(element fsm.SlotElementHelper) (interface{}, fsm.ElementState, error) {
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

func (s *RawGetCodeStateMachine) migrateFromPresentFirst(element fsm.SlotElementHelper) (interface{}, fsm.ElementState, error) {
	aInput, ok := element.GetInputEvent().(Event)
	if !ok {
		return nil, 0, errors.New("wrong input event type")
	}
	aPayload, ok := element.GetPayload().(*Payload)
	if !ok {
		return nil, 0, errors.New("wrong payload type")
	}
	payload, err := s.cleanStateMachine.migrateFromPresentFirst(aInput, aPayload)
	return payload, fsm.NewElementState(element.GetType(), element.GetState()), err
}

func (s *RawGetCodeStateMachine) errorPresentFirst(element fsm.SlotElementHelper, err error) (interface{}, fsm.ElementState) {
	payload, state := s.cleanStateMachine.errorPresentFirst(element.GetInputEvent(), element.GetPayload(), err)
	return payload, state
}

func (s *RawGetCodeStateMachine) responsePresentFirst(element fsm.SlotElementHelper, ar interface{}) (interface{}, fsm.ElementState, error) {
	aInput, ok := element.GetInputEvent().(Event)
	if !ok {
		return nil, 0, errors.New("wrong input event type")
	}
	aResponse, ok := ar.(artifactmanager.GetCodeResp)
	if !ok {
		return nil, 0, errors.New("wrong response type")
	}
	payload, state, err := s.cleanStateMachine.responsePresentFirst(aInput, aResponse, element)
	return payload, state, err
}

func (s *RawGetCodeStateMachine) errorResponsePresentFirst(element fsm.SlotElementHelper, ar interface{}, err error) (interface{}, fsm.ElementState) {
	payload, state := s.cleanStateMachine.errorResponsePresentFirst(element.GetInputEvent(), element.GetPayload(), ar, err)
	return payload, state
}

func (s *RawGetCodeStateMachine) transitPresentSecond(element fsm.SlotElementHelper) (interface{}, fsm.ElementState, error) {
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

func (s *RawGetCodeStateMachine) migrateFromPresentSecond(element fsm.SlotElementHelper) (interface{}, fsm.ElementState, error) {
	aInput, ok := element.GetInputEvent().(Event)
	if !ok {
		return nil, 0, errors.New("wrong input event type")
	}
	aPayload, ok := element.GetPayload().(*Payload)
	if !ok {
		return nil, 0, errors.New("wrong payload type")
	}
	payload, err := s.cleanStateMachine.migrateFromPresentSecond(aInput, aPayload)
	return payload, fsm.NewElementState(element.GetType(), element.GetState()), err
}

func (s *RawGetCodeStateMachine) errorPresentSecond(element fsm.SlotElementHelper, err error) (interface{}, fsm.ElementState) {
	payload, state := s.cleanStateMachine.errorPresentSecond(element.GetInputEvent(), element.GetPayload(), err)
	return payload, state
}

func (s *RawGetCodeStateMachine) responsePresentSecond(element fsm.SlotElementHelper, ar interface{}) (interface{}, fsm.ElementState, error) {
	aInput, ok := element.GetInputEvent().(Event)
	if !ok {
		return nil, 0, errors.New("wrong input event type")
	}
	aPayload, ok := element.GetPayload().(*Payload)
	if !ok {
		return nil, 0, errors.New("wrong payload type")
	}
	aResponse, ok := ar.(artifactmanager.GetCodeResp)
	if !ok {
		return nil, 0, errors.New("wrong response type")
	}
	payload, state, err := s.cleanStateMachine.responsePresentSecond(aInput, aPayload, aResponse)
	return payload, state, err
}

func (s *RawGetCodeStateMachine) errorResponsePresentSecond(element fsm.SlotElementHelper, ar interface{}, err error) (interface{}, fsm.ElementState) {
	payload, state := s.cleanStateMachine.errorResponsePresentSecond(element.GetInputEvent(), element.GetPayload(), ar, err)
	return payload, state
}

func (s *RawGetCodeStateMachine) transitPresentThird(element fsm.SlotElementHelper) (interface{}, fsm.ElementState, error) {
	aInput, ok := element.GetInputEvent().(Event)
	if !ok {
		return nil, 0, errors.New("wrong input event type")
	}
	aPayload, ok := element.GetPayload().(*Payload)
	if !ok {
		return nil, 0, errors.New("wrong payload type")
	}
	payload, state, err := s.cleanStateMachine.transitPresentThird(aInput, aPayload, element)
	return payload, state, err
}

func (s *RawGetCodeStateMachine) migrateFromPresentThird(element fsm.SlotElementHelper) (interface{}, fsm.ElementState, error) {
	aInput, ok := element.GetInputEvent().(Event)
	if !ok {
		return nil, 0, errors.New("wrong input event type")
	}
	aPayload, ok := element.GetPayload().(*Payload)
	if !ok {
		return nil, 0, errors.New("wrong payload type")
	}
	payload, err := s.cleanStateMachine.migrateFromPresentThird(aInput, aPayload)
	return payload, fsm.NewElementState(element.GetType(), element.GetState()), err
}

func (s *RawGetCodeStateMachine) errorPresentThird(element fsm.SlotElementHelper, err error) (interface{}, fsm.ElementState) {
	payload, state := s.cleanStateMachine.errorPresentThird(element.GetInputEvent(), element.GetPayload(), err)
	return payload, state
}

func (s *RawGetCodeStateMachine) responsePresentThird(element fsm.SlotElementHelper, ar interface{}) (interface{}, fsm.ElementState, error) {
	aInput, ok := element.GetInputEvent().(Event)
	if !ok {
		return nil, 0, errors.New("wrong input event type")
	}
	payload, state, err := s.cleanStateMachine.responsePresentThird(aInput, ar)
	return payload, state, err
}

func (s *RawGetCodeStateMachine) errorResponsePresentThird(element fsm.SlotElementHelper, ar interface{}, err error) (interface{}, fsm.ElementState) {
	payload, state := s.cleanStateMachine.errorResponsePresentThird(element.GetInputEvent(), element.GetPayload(), ar, err)
	return payload, state
}

func (s *RawGetCodeStateMachine) transitPresentFourth(element fsm.SlotElementHelper) (interface{}, fsm.ElementState, error) {
	aInput, ok := element.GetInputEvent().(Event)
	if !ok {
		return nil, 0, errors.New("wrong input event type")
	}
	aPayload, ok := element.GetPayload().(*Payload)
	if !ok {
		return nil, 0, errors.New("wrong payload type")
	}
	payload, state, err := s.cleanStateMachine.transitPresentFourth(aInput, aPayload)
	return payload, state, err
}

func (s *RawGetCodeStateMachine) migrateFromPresentFourth(element fsm.SlotElementHelper) (interface{}, fsm.ElementState, error) {
	aInput, ok := element.GetInputEvent().(Event)
	if !ok {
		return nil, 0, errors.New("wrong input event type")
	}
	aPayload, ok := element.GetPayload().(*Payload)
	if !ok {
		return nil, 0, errors.New("wrong payload type")
	}
	payload, err := s.cleanStateMachine.migrateFromPresentFourth(aInput, aPayload)
	return payload, fsm.NewElementState(element.GetType(), element.GetState()), err
}

func (s *RawGetCodeStateMachine) errorPresentFourth(element fsm.SlotElementHelper, err error) (interface{}, fsm.ElementState) {
	payload, state := s.cleanStateMachine.errorPresentFourth(element.GetInputEvent(), element.GetPayload(), err)
	return payload, state
}

func (s *RawGetCodeStateMachine) responsePresentFourth(element fsm.SlotElementHelper, ar interface{}) (interface{}, fsm.ElementState, error) {
	aInput, ok := element.GetInputEvent().(Event)
	if !ok {
		return nil, 0, errors.New("wrong input event type")
	}
	aResponse, ok := ar.(artifactmanager.GetCodeResp)
	if !ok {
		return nil, 0, errors.New("wrong response type")
	}
	payload, state, err := s.cleanStateMachine.responsePresentFourth(aInput, aResponse, element)
	return payload, state, err
}

func (s *RawGetCodeStateMachine) errorResponsePresentFourth(element fsm.SlotElementHelper, ar interface{}, err error) (interface{}, fsm.ElementState) {
	payload, state := s.cleanStateMachine.errorResponsePresentFourth(element.GetInputEvent(), element.GetPayload(), ar, err)
	return payload, state
}
