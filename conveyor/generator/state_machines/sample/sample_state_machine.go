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
	"github.com/insolar/insolar/conveyor/interfaces/adapter"
	"github.com/insolar/insolar/conveyor/interfaces/statemachine"
)

// custom types
type Event struct{}
type Payload struct{}
type TA1 string
type TAR string

// conveyor: state_machine
type TestStateMachine interface {
	TID() statemachine.ID

	i_Init(input Event, payload interface{}) (*Payload, statemachine.ElementState, error)
	if_Init(input Event, payload interface{}) (*Payload, statemachine.ElementState, error)
	ip_Init(input Event, payload interface{}) (*Payload, statemachine.ElementState, error)

	es_Init(input interface{}, payload interface{}, err error) (*Payload, statemachine.ElementState)
	esf_Init(input interface{}, payload interface{}, err error) (*Payload, statemachine.ElementState)
	esp_Init(input interface{}, payload interface{}, err error) (*Payload, statemachine.ElementState)

	// State Declaration
	s_First() statemachine.StateID

	// Migration
	m_FirstSecond(input Event, payload *Payload) (*Payload, statemachine.ElementState, error)
	mfp_FirstSecond(input Event, payload *Payload) (*Payload, statemachine.ElementState, error)

	// Transition
	t_First(input Event, payload *Payload /* todo: , adapterHelper TA1*/) (*Payload, statemachine.ElementState, error)
	tf_First(input Event, payload *Payload /* todo: , adapterHelper TA1*/) (*Payload, statemachine.ElementState, error)
	tp_First(input Event, payload *Payload) (*Payload, statemachine.ElementState, error)

	// TODO: Finalization
	// f_First(input Event, payload *Payload)
	// ff_First(input Event, payload *Payload)
	// fp_First(input Event, payload *Payload)

	// Adapter Response
	a_First(input Event, payload *Payload, respPayload TAR) (*Payload, statemachine.ElementState, error)
	af_First(input Event, payload *Payload, respPayload TAR) (*Payload, statemachine.ElementState, error)
	ap_First(input Event, payload *Payload, respPayload TAR) (*Payload, statemachine.ElementState, error)

	// State Error
	es_First(input interface{}, payload interface{}, err error) (*Payload, statemachine.ElementState)
	esf_First(input interface{}, payload interface{}, err error) (*Payload, statemachine.ElementState)
	esp_First(input interface{}, payload interface{}, err error) (*Payload, statemachine.ElementState)

	// Adapter Response Error
	ea_First(input interface{}, payload interface{}, ar adapter.IAdapterResponse, err error) (*Payload, statemachine.ElementState)
	eaf_First(input interface{}, payload interface{}, ar adapter.IAdapterResponse, err error) (*Payload, statemachine.ElementState)
	eap_First(input interface{}, payload interface{}, ar adapter.IAdapterResponse, err error) (*Payload, statemachine.ElementState)

	// State Declaration
	s_Second() statemachine.StateID

	// Migration
	m_SecondThird(input Event, payload *Payload) (*Payload, statemachine.ElementState, error)
	mfp_SecondThird(input Event, payload *Payload) (*Payload, statemachine.ElementState, error)

	// Transition
	t_Second(input Event, payload *Payload /* todo: , adapterHelper1 TA1*/) (*Payload, statemachine.ElementState, error)
	tf_Second(input Event, payload *Payload /* todo: , adapterHelper1 TA1*/) (*Payload, statemachine.ElementState, error)
	tp_Second(input Event, payload *Payload /* todo: , adapterHelper1 TA1*/) (*Payload, statemachine.ElementState, error)

	// TODO: Finalization
	// f_Second(input Event, payload *Payload)
	// ff_Second(input Event, payload *Payload)
	// fp_Second(input Event, payload *Payload)

	// Adapter Response
	a_Second(input Event, payload *Payload, respPayload TAR) (*Payload, statemachine.ElementState, error)
	af_Second(input Event, payload *Payload, respPayload TAR) (*Payload, statemachine.ElementState, error)
	ap_Second(input Event, payload *Payload, respPayload TAR) (*Payload, statemachine.ElementState, error)

	// State Error
	es_Second(input interface{}, payload interface{}, err error) (*Payload, statemachine.ElementState)
	esf_Second(input interface{}, payload interface{}, err error) (*Payload, statemachine.ElementState)
	esp_Second(input interface{}, payload interface{}, err error) (*Payload, statemachine.ElementState)

	// Adapter Response Error
	ea_Second(input interface{}, payload interface{}, ar adapter.IAdapterResponse, err error) (*Payload, statemachine.ElementState)
	eaf_Second(input interface{}, payload interface{}, ar adapter.IAdapterResponse, err error) (*Payload, statemachine.ElementState)
	eap_Second(input interface{}, payload interface{}, ar adapter.IAdapterResponse, err error) (*Payload, statemachine.ElementState)
}

type TestStateMachineImplementation struct {
	SMFIDTestStateMachine
}

func (t *TestStateMachineImplementation) i_Init(input Event, payload interface{}) (*Payload, statemachine.ElementState, error) {
	return nil, statemachine.ElementState(t.s_First()), nil
}
func (t *TestStateMachineImplementation) if_Init(input Event, payload interface{}) (*Payload, statemachine.ElementState, error) {
	return nil, statemachine.ElementState(t.s_First()), nil
}
func (t *TestStateMachineImplementation) ip_Init(input Event, payload interface{}) (*Payload, statemachine.ElementState, error) {
	return nil, statemachine.ElementState(t.s_First()), nil
}

func (t *TestStateMachineImplementation) es_Init(input interface{}, payload interface{}, err error) (*Payload, statemachine.ElementState) {
	return nil, statemachine.ElementState(t.s_First())
}
func (t *TestStateMachineImplementation) esf_Init(input interface{}, payload interface{}, err error) (*Payload, statemachine.ElementState) {
	return nil, statemachine.ElementState(t.s_First())
}
func (t *TestStateMachineImplementation) esp_Init(input interface{}, payload interface{}, err error) (*Payload, statemachine.ElementState) {
	return nil, statemachine.ElementState(t.s_First())
}

// Migration
func (t *TestStateMachineImplementation) m_FirstSecond(input Event, payload *Payload) (*Payload, statemachine.ElementState, error) {
	return nil, 0, nil
}
func (t *TestStateMachineImplementation) mfp_FirstSecond(input Event, payload *Payload) (*Payload, statemachine.ElementState, error) {
	return nil, 0, nil
}

// Transition
func (t *TestStateMachineImplementation) t_First(input Event, payload *Payload /* todo: , adapterHelper TA1*/) (*Payload, statemachine.ElementState, error) {
	return nil, statemachine.ElementState(t.s_Second()), nil
}
func (t *TestStateMachineImplementation) tf_First(input Event, payload *Payload /* todo: , adapterHelper TA1*/) (*Payload, statemachine.ElementState, error) {
	return nil, statemachine.ElementState(t.s_Second()), nil
}
func (t *TestStateMachineImplementation) tp_First(input Event, payload *Payload) (*Payload, statemachine.ElementState, error) {
	return nil, statemachine.ElementState(t.s_Second()), nil
}

// Adapter Response
func (t *TestStateMachineImplementation) a_First(input Event, payload *Payload, respPayload TAR) (*Payload, statemachine.ElementState, error) {
	return nil, 0, nil
}
func (t *TestStateMachineImplementation) af_First(input Event, payload *Payload, respPayload TAR) (*Payload, statemachine.ElementState, error) {
	return nil, 0, nil
}
func (t *TestStateMachineImplementation) ap_First(input Event, payload *Payload, respPayload TAR) (*Payload, statemachine.ElementState, error) {
	return nil, 0, nil
}

// State Error
func (t *TestStateMachineImplementation) es_First(input interface{}, payload interface{}, err error) (*Payload, statemachine.ElementState) {
	return nil, statemachine.ElementState(t.s_Second())
}
func (t *TestStateMachineImplementation) esf_First(input interface{}, payload interface{}, err error) (*Payload, statemachine.ElementState) {
	return nil, statemachine.ElementState(t.s_Second())
}
func (t *TestStateMachineImplementation) esp_First(input interface{}, payload interface{}, err error) (*Payload, statemachine.ElementState) {
	return nil, statemachine.ElementState(t.s_Second())
}

// Adapter Response Error
func (t *TestStateMachineImplementation) ea_First(input interface{}, payload interface{}, ar adapter.IAdapterResponse, err error) (*Payload, statemachine.ElementState) {
	return nil, 0
}
func (t *TestStateMachineImplementation) eaf_First(input interface{}, payload interface{}, ar adapter.IAdapterResponse, err error) (*Payload, statemachine.ElementState) {
	return nil, 0
}
func (t *TestStateMachineImplementation) eap_First(input interface{}, payload interface{}, ar adapter.IAdapterResponse, err error) (*Payload, statemachine.ElementState) {
	return nil, 0
}

// Migration
func (t *TestStateMachineImplementation) m_SecondThird(input Event, payload *Payload) (*Payload, statemachine.ElementState, error) {
	return nil, 0, nil
}
func (t *TestStateMachineImplementation) mfp_SecondThird(input Event, payload *Payload) (*Payload, statemachine.ElementState, error) {
	return nil, 0, nil
}

// Transition
func (t *TestStateMachineImplementation) t_Second(input Event, payload *Payload /* todo: , adapterHelper1 TA1*/) (*Payload, statemachine.ElementState, error) {
	return nil, 0, nil
}
func (t *TestStateMachineImplementation) tf_Second(input Event, payload *Payload /* todo: , adapterHelper1 TA1*/) (*Payload, statemachine.ElementState, error) {
	return nil, 0, nil
}
func (t *TestStateMachineImplementation) tp_Second(input Event, payload *Payload /* todo: , adapterHelper1 TA1*/) (*Payload, statemachine.ElementState, error) {
	return nil, 0, nil
}

// Adapter Response
func (t *TestStateMachineImplementation) a_Second(input Event, payload *Payload, respPayload TAR) (*Payload, statemachine.ElementState, error) {
	return nil, 0, nil
}
func (t *TestStateMachineImplementation) af_Second(input Event, payload *Payload, respPayload TAR) (*Payload, statemachine.ElementState, error) {
	return nil, 0, nil
}
func (t *TestStateMachineImplementation) ap_Second(input Event, payload *Payload, respPayload TAR) (*Payload, statemachine.ElementState, error) {
	return nil, 0, nil
}

// State Error
func (t *TestStateMachineImplementation) es_Second(input interface{}, payload interface{}, err error) (*Payload, statemachine.ElementState) {
	return nil, 0
}
func (t *TestStateMachineImplementation) esf_Second(input interface{}, payload interface{}, err error) (*Payload, statemachine.ElementState) {
	return nil, 0
}
func (t *TestStateMachineImplementation) esp_Second(input interface{}, payload interface{}, err error) (*Payload, statemachine.ElementState) {
	return nil, 0
}

// Adapter Response Error
func (t *TestStateMachineImplementation) ea_Second(input interface{}, payload interface{}, ar adapter.IAdapterResponse, err error) (*Payload, statemachine.ElementState) {
	return nil, 0
}
func (t *TestStateMachineImplementation) eaf_Second(input interface{}, payload interface{}, ar adapter.IAdapterResponse, err error) (*Payload, statemachine.ElementState) {
	return nil, 0
}
func (t *TestStateMachineImplementation) eap_Second(input interface{}, payload interface{}, ar adapter.IAdapterResponse, err error) (*Payload, statemachine.ElementState) {
	return nil, 0
}
