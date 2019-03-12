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

func (*BaseTestStateMachine) s_First() statemachine.StateID {
    return 1
}
func (*BaseTestStateMachine) s_Second() statemachine.StateID {
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
        Transition: m.if_Init,
        ErrorState: m.esf_Init,
    },
    common.State{
        Transition: m.tf_First,
        AdapterResponse: m.af_First,
        ErrorState: m.esf_First,
        AdapterResponseError: m.eaf_First,
    },
    common.State{
        Transition: m.tf_Second,
        AdapterResponse: m.af_Second,
        ErrorState: m.esf_Second,
        AdapterResponseError: m.eaf_Second,
    },)

    // present state machine
    x[1] = append(x[1], common.State{
        Transition: m.i_Init,
        ErrorState: m.es_Init,
    },
    common.State{
        Migration: m.mfp_FirstSecond,
        Transition: m.t_First,
        AdapterResponse: m.a_First,
        ErrorState: m.es_First,
        AdapterResponseError: m.ea_First,
    },
    common.State{
        Migration: m.mfp_SecondThird,
        Transition: m.t_Second,
        AdapterResponse: m.a_Second,
        ErrorState: m.es_Second,
        AdapterResponseError: m.ea_Second,
    },)

    // past state machine
    x[2] = append(x[2], common.State{
        Transition: m.ip_Init,
        ErrorState: m.esp_Init,
    },
    common.State{
        Transition: m.tp_First,
        AdapterResponse: m.ap_First,
        ErrorState: m.esp_First,
        AdapterResponseError: m.eap_First,
    },
    common.State{
        Transition: m.tp_Second,
        AdapterResponse: m.ap_Second,
        ErrorState: m.esp_Second,
        AdapterResponseError: m.eap_Second,
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

func (s *RawTestStateMachine) i_Init(element slot.SlotElementHelper) (interface{}, statemachine.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    payload, state, err := s.cleanStateMachine.i_Init(aInput, element.GetPayload())
    return payload, state, err
}
func (s *RawTestStateMachine) if_Init(element slot.SlotElementHelper) (interface{}, statemachine.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    payload, state, err := s.cleanStateMachine.if_Init(aInput, element.GetPayload())
    return payload, state, err
}
func (s *RawTestStateMachine) ip_Init(element slot.SlotElementHelper) (interface{}, statemachine.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    payload, state, err := s.cleanStateMachine.ip_Init(aInput, element.GetPayload())
    return payload, state, err
}
func (s *RawTestStateMachine) es_Init(element slot.SlotElementHelper, err error) (interface{}, statemachine.ElementState) {
    payload, state := s.cleanStateMachine.es_Init(element.GetInputEvent(), element.GetPayload(), err)
    return payload, state
}
func (s *RawTestStateMachine) esf_Init(element slot.SlotElementHelper, err error) (interface{}, statemachine.ElementState) {
    payload, state := s.cleanStateMachine.esf_Init(element.GetInputEvent(), element.GetPayload(), err)
    return payload, state
}
func (s *RawTestStateMachine) esp_Init(element slot.SlotElementHelper, err error) (interface{}, statemachine.ElementState) {
    payload, state := s.cleanStateMachine.esp_Init(element.GetInputEvent(), element.GetPayload(), err)
    return payload, state
}


func (s *RawTestStateMachine) t_First(element slot.SlotElementHelper) (interface{}, statemachine.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    payload, state, err := s.cleanStateMachine.t_First(aInput, aPayload)
    return payload, state, err
}
func (s *RawTestStateMachine) tf_First(element slot.SlotElementHelper) (interface{}, statemachine.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    payload, state, err := s.cleanStateMachine.tf_First(aInput, aPayload)
    return payload, state, err
}
func (s *RawTestStateMachine) tp_First(element slot.SlotElementHelper) (interface{}, statemachine.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    payload, state, err := s.cleanStateMachine.tp_First(aInput, aPayload)
    return payload, state, err
}
func (s *RawTestStateMachine) m_FirstSecond(element slot.SlotElementHelper) (interface{}, statemachine.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    payload, state, err := s.cleanStateMachine.m_FirstSecond(aInput, aPayload)
    return payload, state, err
}
func (s *RawTestStateMachine) mfp_FirstSecond(element slot.SlotElementHelper) (interface{}, statemachine.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    payload, state, err := s.cleanStateMachine.mfp_FirstSecond(aInput, aPayload)
    return payload, state, err
}
func (s *RawTestStateMachine) es_First(element slot.SlotElementHelper, err error) (interface{}, statemachine.ElementState) {
    payload, state := s.cleanStateMachine.es_First(element.GetInputEvent(), element.GetPayload(), err)
    return payload, state
}
func (s *RawTestStateMachine) esf_First(element slot.SlotElementHelper, err error) (interface{}, statemachine.ElementState) {
    payload, state := s.cleanStateMachine.esf_First(element.GetInputEvent(), element.GetPayload(), err)
    return payload, state
}
func (s *RawTestStateMachine) esp_First(element slot.SlotElementHelper, err error) (interface{}, statemachine.ElementState) {
    payload, state := s.cleanStateMachine.esp_First(element.GetInputEvent(), element.GetPayload(), err)
    return payload, state
}
func (s *RawTestStateMachine) a_First(element slot.SlotElementHelper, ar adapter.IAdapterResponse) (interface{}, statemachine.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    aResponse, ok := ar.GetRespPayload().(TAR)
    if !ok { return nil, 0, errors.New("wrong response type") }
    payload, state, err := s.cleanStateMachine.a_First(aInput, aPayload, aResponse)
    return payload, state, err
}
func (s *RawTestStateMachine) af_First(element slot.SlotElementHelper, ar adapter.IAdapterResponse) (interface{}, statemachine.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    aResponse, ok := ar.GetRespPayload().(TAR)
    if !ok { return nil, 0, errors.New("wrong response type") }
    payload, state, err := s.cleanStateMachine.af_First(aInput, aPayload, aResponse)
    return payload, state, err
}
func (s *RawTestStateMachine) ap_First(element slot.SlotElementHelper, ar adapter.IAdapterResponse) (interface{}, statemachine.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    aResponse, ok := ar.GetRespPayload().(TAR)
    if !ok { return nil, 0, errors.New("wrong response type") }
    payload, state, err := s.cleanStateMachine.ap_First(aInput, aPayload, aResponse)
    return payload, state, err
}
func (s *RawTestStateMachine) ea_First(element slot.SlotElementHelper, ar adapter.IAdapterResponse, err error) (interface{}, statemachine.ElementState) {
    payload, state := s.cleanStateMachine.ea_First(element.GetInputEvent(), element.GetPayload(), ar, err)
    return payload, state
}
func (s *RawTestStateMachine) eaf_First(element slot.SlotElementHelper, ar adapter.IAdapterResponse, err error) (interface{}, statemachine.ElementState) {
    payload, state := s.cleanStateMachine.eaf_First(element.GetInputEvent(), element.GetPayload(), ar, err)
    return payload, state
}
func (s *RawTestStateMachine) eap_First(element slot.SlotElementHelper, ar adapter.IAdapterResponse, err error) (interface{}, statemachine.ElementState) {
    payload, state := s.cleanStateMachine.eap_First(element.GetInputEvent(), element.GetPayload(), ar, err)
    return payload, state
}

func (s *RawTestStateMachine) t_Second(element slot.SlotElementHelper) (interface{}, statemachine.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    payload, state, err := s.cleanStateMachine.t_Second(aInput, aPayload)
    return payload, state, err
}
func (s *RawTestStateMachine) tf_Second(element slot.SlotElementHelper) (interface{}, statemachine.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    payload, state, err := s.cleanStateMachine.tf_Second(aInput, aPayload)
    return payload, state, err
}
func (s *RawTestStateMachine) tp_Second(element slot.SlotElementHelper) (interface{}, statemachine.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    payload, state, err := s.cleanStateMachine.tp_Second(aInput, aPayload)
    return payload, state, err
}
func (s *RawTestStateMachine) m_SecondThird(element slot.SlotElementHelper) (interface{}, statemachine.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    payload, state, err := s.cleanStateMachine.m_SecondThird(aInput, aPayload)
    return payload, state, err
}
func (s *RawTestStateMachine) mfp_SecondThird(element slot.SlotElementHelper) (interface{}, statemachine.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    payload, state, err := s.cleanStateMachine.mfp_SecondThird(aInput, aPayload)
    return payload, state, err
}
func (s *RawTestStateMachine) es_Second(element slot.SlotElementHelper, err error) (interface{}, statemachine.ElementState) {
    payload, state := s.cleanStateMachine.es_Second(element.GetInputEvent(), element.GetPayload(), err)
    return payload, state
}
func (s *RawTestStateMachine) esf_Second(element slot.SlotElementHelper, err error) (interface{}, statemachine.ElementState) {
    payload, state := s.cleanStateMachine.esf_Second(element.GetInputEvent(), element.GetPayload(), err)
    return payload, state
}
func (s *RawTestStateMachine) esp_Second(element slot.SlotElementHelper, err error) (interface{}, statemachine.ElementState) {
    payload, state := s.cleanStateMachine.esp_Second(element.GetInputEvent(), element.GetPayload(), err)
    return payload, state
}
func (s *RawTestStateMachine) a_Second(element slot.SlotElementHelper, ar adapter.IAdapterResponse) (interface{}, statemachine.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    aResponse, ok := ar.GetRespPayload().(TAR)
    if !ok { return nil, 0, errors.New("wrong response type") }
    payload, state, err := s.cleanStateMachine.a_Second(aInput, aPayload, aResponse)
    return payload, state, err
}
func (s *RawTestStateMachine) af_Second(element slot.SlotElementHelper, ar adapter.IAdapterResponse) (interface{}, statemachine.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    aResponse, ok := ar.GetRespPayload().(TAR)
    if !ok { return nil, 0, errors.New("wrong response type") }
    payload, state, err := s.cleanStateMachine.af_Second(aInput, aPayload, aResponse)
    return payload, state, err
}
func (s *RawTestStateMachine) ap_Second(element slot.SlotElementHelper, ar adapter.IAdapterResponse) (interface{}, statemachine.ElementState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    aResponse, ok := ar.GetRespPayload().(TAR)
    if !ok { return nil, 0, errors.New("wrong response type") }
    payload, state, err := s.cleanStateMachine.ap_Second(aInput, aPayload, aResponse)
    return payload, state, err
}
func (s *RawTestStateMachine) ea_Second(element slot.SlotElementHelper, ar adapter.IAdapterResponse, err error) (interface{}, statemachine.ElementState) {
    payload, state := s.cleanStateMachine.ea_Second(element.GetInputEvent(), element.GetPayload(), ar, err)
    return payload, state
}
func (s *RawTestStateMachine) eaf_Second(element slot.SlotElementHelper, ar adapter.IAdapterResponse, err error) (interface{}, statemachine.ElementState) {
    payload, state := s.cleanStateMachine.eaf_Second(element.GetInputEvent(), element.GetPayload(), ar, err)
    return payload, state
}
func (s *RawTestStateMachine) eap_Second(element slot.SlotElementHelper, ar adapter.IAdapterResponse, err error) (interface{}, statemachine.ElementState) {
    payload, state := s.cleanStateMachine.eap_Second(element.GetInputEvent(), element.GetPayload(), ar, err)
    return payload, state
}









