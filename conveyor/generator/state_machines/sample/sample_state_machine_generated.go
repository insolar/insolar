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

type SMFIDTestStateMachine struct{}

func (*SMFIDTestStateMachine) TID() statemachine.ID {
	return 1
}

func (*SMFIDTestStateMachine) s_First() statemachine.StateID {
	return 1
}
func (*SMFIDTestStateMachine) s_Second() statemachine.StateID {
	return 2
}

type SMRHTestStateMachine struct {
	cleanHandlers TestStateMachine
}

func SMRHTestStateMachineFactory() [3]common.StateMachine {
	m := SMRHTestStateMachine{
		cleanHandlers: &TestStateMachineImplementation{},
	}

	var x = [3][]common.State{}
	// future state machine
	x[0] = append(x[0], common.State{
		Transition: m.if_Init,
		ErrorState: m.esf_Init,
	},
		common.State{
			Transition:           m.tf_First,
			AdapterResponse:      m.af_First,
			ErrorState:           m.esf_First,
			AdapterResponseError: m.eaf_First,
		},
		common.State{
			Transition:           m.tf_Second,
			AdapterResponse:      m.af_Second,
			ErrorState:           m.esf_Second,
			AdapterResponseError: m.eaf_Second,
		})

	// present state machine
	x[1] = append(x[1], common.State{
		Transition: m.i_Init,
		ErrorState: m.es_Init,
	},
		common.State{
			Migration:            m.mfp_FirstSecond,
			Transition:           m.t_First,
			AdapterResponse:      m.a_First,
			ErrorState:           m.es_First,
			AdapterResponseError: m.ea_First,
		},
		common.State{
			Migration:            m.mfp_SecondThird,
			Transition:           m.t_Second,
			AdapterResponse:      m.a_Second,
			ErrorState:           m.es_Second,
			AdapterResponseError: m.ea_Second,
		})

	// past state machine
	x[2] = append(x[2], common.State{
		Transition: m.ip_Init,
		ErrorState: m.esp_Init,
	},
		common.State{
			Transition:           m.tp_First,
			AdapterResponse:      m.ap_First,
			ErrorState:           m.esp_First,
			AdapterResponseError: m.eap_First,
		},
		common.State{
			Transition:           m.tp_Second,
			AdapterResponse:      m.ap_Second,
			ErrorState:           m.esp_Second,
			AdapterResponseError: m.eap_Second,
		})

	return [3]common.StateMachine{
		common.StateMachine{
			ID:     int(m.cleanHandlers.(TestStateMachine).TID()),
			States: x[0],
		},
		common.StateMachine{
			ID:     int(m.cleanHandlers.(TestStateMachine).TID()),
			States: x[0],
		},
		common.StateMachine{
			ID:     int(m.cleanHandlers.(TestStateMachine).TID()),
			States: x[0],
		},
	}
}

func (s *SMRHTestStateMachine) i_Init(element slot.SlotElementHelper) (interface{}, statemachine.ElementState, error) {
	aInput, ok := element.GetInputEvent().(Event)
	if !ok {
		return nil, 0, errors.New("wrong input event type")
	}
	payload, state, err := s.cleanHandlers.i_Init(aInput, element.GetPayload())
	return payload, state, err
}
func (s *SMRHTestStateMachine) if_Init(element slot.SlotElementHelper) (interface{}, statemachine.ElementState, error) {
	aInput, ok := element.GetInputEvent().(Event)
	if !ok {
		return nil, 0, errors.New("wrong input event type")
	}
	payload, state, err := s.cleanHandlers.if_Init(aInput, element.GetPayload())
	return payload, state, err
}
func (s *SMRHTestStateMachine) ip_Init(element slot.SlotElementHelper) (interface{}, statemachine.ElementState, error) {
	aInput, ok := element.GetInputEvent().(Event)
	if !ok {
		return nil, 0, errors.New("wrong input event type")
	}
	payload, state, err := s.cleanHandlers.ip_Init(aInput, element.GetPayload())
	return payload, state, err
}
func (s *SMRHTestStateMachine) es_Init(element slot.SlotElementHelper, err error) (interface{}, statemachine.ElementState) {
	payload, state := s.cleanHandlers.es_Init(element.GetInputEvent(), element.GetPayload(), err)
	return payload, state
}
func (s *SMRHTestStateMachine) esf_Init(element slot.SlotElementHelper, err error) (interface{}, statemachine.ElementState) {
	payload, state := s.cleanHandlers.esf_Init(element.GetInputEvent(), element.GetPayload(), err)
	return payload, state
}
func (s *SMRHTestStateMachine) esp_Init(element slot.SlotElementHelper, err error) (interface{}, statemachine.ElementState) {
	payload, state := s.cleanHandlers.esp_Init(element.GetInputEvent(), element.GetPayload(), err)
	return payload, state
}

func (s *SMRHTestStateMachine) t_First(element slot.SlotElementHelper) (interface{}, statemachine.ElementState, error) {
	aInput, ok := element.GetInputEvent().(Event)
	if !ok {
		return nil, 0, errors.New("wrong input event type")
	}
	aPayload, ok := element.GetPayload().(*Payload)
	if !ok {
		return nil, 0, errors.New("wrong payload type")
	}
	payload, state, err := s.cleanHandlers.t_First(aInput, aPayload)
	return payload, state, err
}
func (s *SMRHTestStateMachine) tf_First(element slot.SlotElementHelper) (interface{}, statemachine.ElementState, error) {
	aInput, ok := element.GetInputEvent().(Event)
	if !ok {
		return nil, 0, errors.New("wrong input event type")
	}
	aPayload, ok := element.GetPayload().(*Payload)
	if !ok {
		return nil, 0, errors.New("wrong payload type")
	}
	payload, state, err := s.cleanHandlers.tf_First(aInput, aPayload)
	return payload, state, err
}
func (s *SMRHTestStateMachine) tp_First(element slot.SlotElementHelper) (interface{}, statemachine.ElementState, error) {
	aInput, ok := element.GetInputEvent().(Event)
	if !ok {
		return nil, 0, errors.New("wrong input event type")
	}
	aPayload, ok := element.GetPayload().(*Payload)
	if !ok {
		return nil, 0, errors.New("wrong payload type")
	}
	payload, state, err := s.cleanHandlers.tp_First(aInput, aPayload)
	return payload, state, err
}
func (s *SMRHTestStateMachine) m_FirstSecond(element slot.SlotElementHelper) (interface{}, statemachine.ElementState, error) {
	aInput, ok := element.GetInputEvent().(Event)
	if !ok {
		return nil, 0, errors.New("wrong input event type")
	}
	aPayload, ok := element.GetPayload().(*Payload)
	if !ok {
		return nil, 0, errors.New("wrong payload type")
	}
	payload, state, err := s.cleanHandlers.m_FirstSecond(aInput, aPayload)
	return payload, state, err
}
func (s *SMRHTestStateMachine) mfp_FirstSecond(element slot.SlotElementHelper) (interface{}, statemachine.ElementState, error) {
	aInput, ok := element.GetInputEvent().(Event)
	if !ok {
		return nil, 0, errors.New("wrong input event type")
	}
	aPayload, ok := element.GetPayload().(*Payload)
	if !ok {
		return nil, 0, errors.New("wrong payload type")
	}
	payload, state, err := s.cleanHandlers.mfp_FirstSecond(aInput, aPayload)
	return payload, state, err
}
func (s *SMRHTestStateMachine) es_First(element slot.SlotElementHelper, err error) (interface{}, statemachine.ElementState) {
	payload, state := s.cleanHandlers.es_First(element.GetInputEvent(), element.GetPayload(), err)
	return payload, state
}
func (s *SMRHTestStateMachine) esf_First(element slot.SlotElementHelper, err error) (interface{}, statemachine.ElementState) {
	payload, state := s.cleanHandlers.esf_First(element.GetInputEvent(), element.GetPayload(), err)
	return payload, state
}
func (s *SMRHTestStateMachine) esp_First(element slot.SlotElementHelper, err error) (interface{}, statemachine.ElementState) {
	payload, state := s.cleanHandlers.esp_First(element.GetInputEvent(), element.GetPayload(), err)
	return payload, state
}
func (s *SMRHTestStateMachine) a_First(element slot.SlotElementHelper, ar adapter.IAdapterResponse) (interface{}, statemachine.ElementState, error) {
	aInput, ok := element.GetInputEvent().(Event)
	if !ok {
		return nil, 0, errors.New("wrong input event type")
	}
	aPayload, ok := element.GetPayload().(*Payload)
	if !ok {
		return nil, 0, errors.New("wrong payload type")
	}
	aResponse, ok := ar.GetRespPayload().(TAR)
	if !ok {
		return nil, 0, errors.New("wrong response type")
	}
	payload, state, err := s.cleanHandlers.a_First(aInput, aPayload, aResponse)
	return payload, state, err
}
func (s *SMRHTestStateMachine) af_First(element slot.SlotElementHelper, ar adapter.IAdapterResponse) (interface{}, statemachine.ElementState, error) {
	aInput, ok := element.GetInputEvent().(Event)
	if !ok {
		return nil, 0, errors.New("wrong input event type")
	}
	aPayload, ok := element.GetPayload().(*Payload)
	if !ok {
		return nil, 0, errors.New("wrong payload type")
	}
	aResponse, ok := ar.GetRespPayload().(TAR)
	if !ok {
		return nil, 0, errors.New("wrong response type")
	}
	payload, state, err := s.cleanHandlers.af_First(aInput, aPayload, aResponse)
	return payload, state, err
}
func (s *SMRHTestStateMachine) ap_First(element slot.SlotElementHelper, ar adapter.IAdapterResponse) (interface{}, statemachine.ElementState, error) {
	aInput, ok := element.GetInputEvent().(Event)
	if !ok {
		return nil, 0, errors.New("wrong input event type")
	}
	aPayload, ok := element.GetPayload().(*Payload)
	if !ok {
		return nil, 0, errors.New("wrong payload type")
	}
	aResponse, ok := ar.GetRespPayload().(TAR)
	if !ok {
		return nil, 0, errors.New("wrong response type")
	}
	payload, state, err := s.cleanHandlers.ap_First(aInput, aPayload, aResponse)
	return payload, state, err
}
func (s *SMRHTestStateMachine) ea_First(element slot.SlotElementHelper, ar adapter.IAdapterResponse, err error) (interface{}, statemachine.ElementState) {
	payload, state := s.cleanHandlers.ea_First(element.GetInputEvent(), element.GetPayload(), ar, err)
	return payload, state
}
func (s *SMRHTestStateMachine) eaf_First(element slot.SlotElementHelper, ar adapter.IAdapterResponse, err error) (interface{}, statemachine.ElementState) {
	payload, state := s.cleanHandlers.eaf_First(element.GetInputEvent(), element.GetPayload(), ar, err)
	return payload, state
}
func (s *SMRHTestStateMachine) eap_First(element slot.SlotElementHelper, ar adapter.IAdapterResponse, err error) (interface{}, statemachine.ElementState) {
	payload, state := s.cleanHandlers.eap_First(element.GetInputEvent(), element.GetPayload(), ar, err)
	return payload, state
}

func (s *SMRHTestStateMachine) t_Second(element slot.SlotElementHelper) (interface{}, statemachine.ElementState, error) {
	aInput, ok := element.GetInputEvent().(Event)
	if !ok {
		return nil, 0, errors.New("wrong input event type")
	}
	aPayload, ok := element.GetPayload().(*Payload)
	if !ok {
		return nil, 0, errors.New("wrong payload type")
	}
	payload, state, err := s.cleanHandlers.t_Second(aInput, aPayload)
	return payload, state, err
}
func (s *SMRHTestStateMachine) tf_Second(element slot.SlotElementHelper) (interface{}, statemachine.ElementState, error) {
	aInput, ok := element.GetInputEvent().(Event)
	if !ok {
		return nil, 0, errors.New("wrong input event type")
	}
	aPayload, ok := element.GetPayload().(*Payload)
	if !ok {
		return nil, 0, errors.New("wrong payload type")
	}
	payload, state, err := s.cleanHandlers.tf_Second(aInput, aPayload)
	return payload, state, err
}
func (s *SMRHTestStateMachine) tp_Second(element slot.SlotElementHelper) (interface{}, statemachine.ElementState, error) {
	aInput, ok := element.GetInputEvent().(Event)
	if !ok {
		return nil, 0, errors.New("wrong input event type")
	}
	aPayload, ok := element.GetPayload().(*Payload)
	if !ok {
		return nil, 0, errors.New("wrong payload type")
	}
	payload, state, err := s.cleanHandlers.tp_Second(aInput, aPayload)
	return payload, state, err
}
func (s *SMRHTestStateMachine) m_SecondThird(element slot.SlotElementHelper) (interface{}, statemachine.ElementState, error) {
	aInput, ok := element.GetInputEvent().(Event)
	if !ok {
		return nil, 0, errors.New("wrong input event type")
	}
	aPayload, ok := element.GetPayload().(*Payload)
	if !ok {
		return nil, 0, errors.New("wrong payload type")
	}
	payload, state, err := s.cleanHandlers.m_SecondThird(aInput, aPayload)
	return payload, state, err
}
func (s *SMRHTestStateMachine) mfp_SecondThird(element slot.SlotElementHelper) (interface{}, statemachine.ElementState, error) {
	aInput, ok := element.GetInputEvent().(Event)
	if !ok {
		return nil, 0, errors.New("wrong input event type")
	}
	aPayload, ok := element.GetPayload().(*Payload)
	if !ok {
		return nil, 0, errors.New("wrong payload type")
	}
	payload, state, err := s.cleanHandlers.mfp_SecondThird(aInput, aPayload)
	return payload, state, err
}
func (s *SMRHTestStateMachine) es_Second(element slot.SlotElementHelper, err error) (interface{}, statemachine.ElementState) {
	payload, state := s.cleanHandlers.es_Second(element.GetInputEvent(), element.GetPayload(), err)
	return payload, state
}
func (s *SMRHTestStateMachine) esf_Second(element slot.SlotElementHelper, err error) (interface{}, statemachine.ElementState) {
	payload, state := s.cleanHandlers.esf_Second(element.GetInputEvent(), element.GetPayload(), err)
	return payload, state
}
func (s *SMRHTestStateMachine) esp_Second(element slot.SlotElementHelper, err error) (interface{}, statemachine.ElementState) {
	payload, state := s.cleanHandlers.esp_Second(element.GetInputEvent(), element.GetPayload(), err)
	return payload, state
}
func (s *SMRHTestStateMachine) a_Second(element slot.SlotElementHelper, ar adapter.IAdapterResponse) (interface{}, statemachine.ElementState, error) {
	aInput, ok := element.GetInputEvent().(Event)
	if !ok {
		return nil, 0, errors.New("wrong input event type")
	}
	aPayload, ok := element.GetPayload().(*Payload)
	if !ok {
		return nil, 0, errors.New("wrong payload type")
	}
	aResponse, ok := ar.GetRespPayload().(TAR)
	if !ok {
		return nil, 0, errors.New("wrong response type")
	}
	payload, state, err := s.cleanHandlers.a_Second(aInput, aPayload, aResponse)
	return payload, state, err
}
func (s *SMRHTestStateMachine) af_Second(element slot.SlotElementHelper, ar adapter.IAdapterResponse) (interface{}, statemachine.ElementState, error) {
	aInput, ok := element.GetInputEvent().(Event)
	if !ok {
		return nil, 0, errors.New("wrong input event type")
	}
	aPayload, ok := element.GetPayload().(*Payload)
	if !ok {
		return nil, 0, errors.New("wrong payload type")
	}
	aResponse, ok := ar.GetRespPayload().(TAR)
	if !ok {
		return nil, 0, errors.New("wrong response type")
	}
	payload, state, err := s.cleanHandlers.af_Second(aInput, aPayload, aResponse)
	return payload, state, err
}
func (s *SMRHTestStateMachine) ap_Second(element slot.SlotElementHelper, ar adapter.IAdapterResponse) (interface{}, statemachine.ElementState, error) {
	aInput, ok := element.GetInputEvent().(Event)
	if !ok {
		return nil, 0, errors.New("wrong input event type")
	}
	aPayload, ok := element.GetPayload().(*Payload)
	if !ok {
		return nil, 0, errors.New("wrong payload type")
	}
	aResponse, ok := ar.GetRespPayload().(TAR)
	if !ok {
		return nil, 0, errors.New("wrong response type")
	}
	payload, state, err := s.cleanHandlers.ap_Second(aInput, aPayload, aResponse)
	return payload, state, err
}
func (s *SMRHTestStateMachine) ea_Second(element slot.SlotElementHelper, ar adapter.IAdapterResponse, err error) (interface{}, statemachine.ElementState) {
	payload, state := s.cleanHandlers.ea_Second(element.GetInputEvent(), element.GetPayload(), ar, err)
	return payload, state
}
func (s *SMRHTestStateMachine) eaf_Second(element slot.SlotElementHelper, ar adapter.IAdapterResponse, err error) (interface{}, statemachine.ElementState) {
	payload, state := s.cleanHandlers.eaf_Second(element.GetInputEvent(), element.GetPayload(), ar, err)
	return payload, state
}
func (s *SMRHTestStateMachine) eap_Second(element slot.SlotElementHelper, ar adapter.IAdapterResponse, err error) (interface{}, statemachine.ElementState) {
	payload, state := s.cleanHandlers.eap_Second(element.GetInputEvent(), element.GetPayload(), ar, err)
	return payload, state
}
