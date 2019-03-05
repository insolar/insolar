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
)

// custom types
type Event struct{}
type Payload string
type TA1 string
type TAR string

// conveyor: state_machine
type TestStateMachine interface {
	TID() common.ElType

	i_Init(input Event, payload interface{}) (*Payload, common.ElUpdate, error)
	if_Init(input Event, payload interface{}) (*Payload, common.ElUpdate, error)
	ip_Init(input Event, payload interface{}) (*Payload, common.ElUpdate, error)

	es_Init(input interface{}, payload interface{}, err error) (*Payload, common.ElUpdate)
	esf_Init(input interface{}, payload interface{}, err error) (*Payload, common.ElUpdate)
	esp_Init(input interface{}, payload interface{}, err error) (*Payload, common.ElUpdate)

	// State Declaration
	s_First() common.ElState

	// Migration
	m_FirstSecond(input Event, payload *Payload) (*Payload, common.ElUpdate, error)
	mfp_FirstSecond(input Event, payload *Payload) (*Payload, common.ElUpdate, error)

	// Transition
	t_First(input Event, payload *Payload/* todo: , adapterHelper TA1*/) (*Payload, common.ElUpdate, error)
	tf_First(input Event, payload *Payload/* todo: , adapterHelper TA1*/) (*Payload, common.ElUpdate, error)
	tp_First(input Event, payload *Payload) (*Payload, common.ElUpdate, error)

	// todo: Finalization
	// f_First(input Event, payload *Payload)
	// ff_First(input Event, payload *Payload)
	// fp_First(input Event, payload *Payload)

	// Adapter Response
	a_First(input Event, payload *Payload, respPayload TAR) (*Payload, common.ElUpdate, error)
	af_First(input Event, payload *Payload, respPayload TAR) (*Payload, common.ElUpdate, error)
	ap_First(input Event, payload *Payload, respPayload TAR) (*Payload, common.ElUpdate, error)

	// State Error
	es_First(input interface{}, payload interface{}, err error) (*Payload, common.ElUpdate)
	esf_First(input interface{}, payload interface{}, err error) (*Payload, common.ElUpdate)
	esp_First(input interface{}, payload interface{}, err error) (*Payload, common.ElUpdate)

	// Adapter Response Error
	ea_First(input interface{}, payload interface{}, ar adapter.IAdapterResponse, err error) (*Payload, common.ElUpdate)
	eaf_First(input interface{}, payload interface{}, ar adapter.IAdapterResponse, err error) (*Payload, common.ElUpdate)
	eap_First(input interface{}, payload interface{}, ar adapter.IAdapterResponse, err error) (*Payload, common.ElUpdate)

	// State Declaration
	s_Second() common.ElState

	// Migration
	m_SecondThird(input Event, payload *Payload) (*Payload, common.ElUpdate, error)
	mfp_SecondThird(input Event, payload *Payload) (*Payload, common.ElUpdate, error)

	// Transition
	t_Second(input Event, payload *Payload/* todo: , adapterHelper1 TA1*/) (*Payload, common.ElUpdate, error)
	tf_Second(input Event, payload *Payload/* todo: , adapterHelper1 TA1*/) (*Payload, common.ElUpdate, error)
	tp_Second(input Event, payload *Payload/* todo: , adapterHelper1 TA1*/) (*Payload, common.ElUpdate, error)

	// todo: Finalization
	// f_Second(input Event, payload *Payload)
	// ff_Second(input Event, payload *Payload)
	// fp_Second(input Event, payload *Payload)

	// Adapter Response
	a_Second(input Event, payload *Payload, respPayload TAR) (*Payload, common.ElUpdate, error)
	af_Second(input Event, payload *Payload, respPayload TAR) (*Payload, common.ElUpdate, error)
	ap_Second(input Event, payload *Payload, respPayload TAR) (*Payload, common.ElUpdate, error)

	// State Error
	es_Second(input interface{}, payload interface{}, err error) (*Payload, common.ElUpdate)
	esf_Second(input interface{}, payload interface{}, err error) (*Payload, common.ElUpdate)
	esp_Second(input interface{}, payload interface{}, err error) (*Payload, common.ElUpdate)

	// Adapter Response Error
	ea_Second(input interface{}, payload interface{}, ar adapter.IAdapterResponse, err error) (*Payload, common.ElUpdate)
	eaf_Second(input interface{}, payload interface{}, ar adapter.IAdapterResponse, err error) (*Payload, common.ElUpdate)
	eap_Second(input interface{}, payload interface{}, ar adapter.IAdapterResponse, err error) (*Payload, common.ElUpdate)
}

type TestStateMachineImplementation struct {
	SMFIDTestStateMachine
}

func (t *TestStateMachineImplementation) i_Init(input Event, payload interface{}) (*Payload, common.ElUpdate, error) {
	return nil, common.ElUpdate(t.s_First()), nil
}
func (t *TestStateMachineImplementation) if_Init(input Event, payload interface{}) (*Payload, common.ElUpdate, error) {
	return nil, common.ElUpdate(t.s_First()), nil
}
func (t *TestStateMachineImplementation) ip_Init(input Event, payload interface{}) (*Payload, common.ElUpdate, error) {
	return nil, common.ElUpdate(t.s_First()), nil
}

func (t *TestStateMachineImplementation) es_Init(input interface{}, payload interface{}, err error) (*Payload, common.ElUpdate) {
	return nil, common.ElUpdate(t.s_First())
}
func (t *TestStateMachineImplementation) esf_Init(input interface{}, payload interface{}, err error) (*Payload, common.ElUpdate) {
	return nil, common.ElUpdate(t.s_First())
}
func (t *TestStateMachineImplementation) esp_Init(input interface{}, payload interface{}, err error) (*Payload, common.ElUpdate) {
	return nil, common.ElUpdate(t.s_First())
}

// Migration
func (t *TestStateMachineImplementation) m_FirstSecond(input Event, payload *Payload) (*Payload, common.ElUpdate, error) {
	return nil, 0, nil
}
func (t *TestStateMachineImplementation) mfp_FirstSecond(input Event, payload *Payload) (*Payload, common.ElUpdate, error) {
	return nil, 0, nil
}

// Transition
func (t *TestStateMachineImplementation) t_First(input Event, payload *Payload/* todo: , adapterHelper TA1*/) (*Payload, common.ElUpdate, error) {
	return nil, common.ElUpdate(t.s_Second()), nil
}
func (t *TestStateMachineImplementation) tf_First(input Event, payload *Payload/* todo: , adapterHelper TA1*/) (*Payload, common.ElUpdate, error) {
	return nil, common.ElUpdate(t.s_Second()), nil
}
func (t *TestStateMachineImplementation) tp_First(input Event, payload *Payload) (*Payload, common.ElUpdate, error) {
	return nil, common.ElUpdate(t.s_Second()), nil
}

// Adapter Response
func (t *TestStateMachineImplementation) a_First(input Event, payload *Payload, respPayload TAR) (*Payload, common.ElUpdate, error) {
	return nil, 0, nil
}
func (t *TestStateMachineImplementation) af_First(input Event, payload *Payload, respPayload TAR) (*Payload, common.ElUpdate, error) {
	return nil, 0, nil
}
func (t *TestStateMachineImplementation) ap_First(input Event, payload *Payload, respPayload TAR) (*Payload, common.ElUpdate, error) {
	return nil, 0, nil
}

// State Error
func (t *TestStateMachineImplementation) es_First(input interface{}, payload interface{}, err error) (*Payload, common.ElUpdate) {
	return nil, common.ElUpdate(t.s_Second())
}
func (t *TestStateMachineImplementation) esf_First(input interface{}, payload interface{}, err error) (*Payload, common.ElUpdate) {
	return nil, common.ElUpdate(t.s_Second())
}
func (t *TestStateMachineImplementation) esp_First(input interface{}, payload interface{}, err error) (*Payload, common.ElUpdate) {
	return nil, common.ElUpdate(t.s_Second())
}

// Adapter Response Error
func (t *TestStateMachineImplementation) ea_First(input interface{}, payload interface{}, ar adapter.IAdapterResponse, err error) (*Payload, common.ElUpdate) {
	return nil, 0
}
func (t *TestStateMachineImplementation) eaf_First(input interface{}, payload interface{}, ar adapter.IAdapterResponse, err error) (*Payload, common.ElUpdate) {
	return nil, 0
}
func (t *TestStateMachineImplementation) eap_First(input interface{}, payload interface{}, ar adapter.IAdapterResponse, err error) (*Payload, common.ElUpdate) {
	return nil, 0
}

// Migration
func (t *TestStateMachineImplementation) m_SecondThird(input Event, payload *Payload) (*Payload, common.ElUpdate, error) {
	return nil, 0, nil
}
func (t *TestStateMachineImplementation) mfp_SecondThird(input Event, payload *Payload) (*Payload, common.ElUpdate, error) {
	return nil, 0, nil
}

// Transition
func (t *TestStateMachineImplementation) t_Second(input Event, payload *Payload/* todo: , adapterHelper1 TA1*/) (*Payload, common.ElUpdate, error) {
	return nil, 0, nil
}
func (t *TestStateMachineImplementation) tf_Second(input Event, payload *Payload/* todo: , adapterHelper1 TA1*/) (*Payload, common.ElUpdate, error) {
	return nil, 0, nil
}
func (t *TestStateMachineImplementation) tp_Second(input Event, payload *Payload/* todo: , adapterHelper1 TA1*/) (*Payload, common.ElUpdate, error) {
	return nil, 0, nil
}

// Adapter Response
func (t *TestStateMachineImplementation) a_Second(input Event, payload *Payload, respPayload TAR) (*Payload, common.ElUpdate, error) {
	return nil, 0, nil
}
func (t *TestStateMachineImplementation) af_Second(input Event, payload *Payload, respPayload TAR) (*Payload, common.ElUpdate, error) {
	return nil, 0, nil
}
func (t *TestStateMachineImplementation) ap_Second(input Event, payload *Payload, respPayload TAR) (*Payload, common.ElUpdate, error) {
	return nil, 0, nil
}

// State Error
func (t *TestStateMachineImplementation) es_Second(input interface{}, payload interface{}, err error) (*Payload, common.ElUpdate) {
	return nil, 0
}
func (t *TestStateMachineImplementation) esf_Second(input interface{}, payload interface{}, err error) (*Payload, common.ElUpdate) {
	return nil, 0
}
func (t *TestStateMachineImplementation) esp_Second(input interface{}, payload interface{}, err error) (*Payload, common.ElUpdate) {
	return nil, 0
}

// Adapter Response Error
func (t *TestStateMachineImplementation) ea_Second(input interface{}, payload interface{}, ar adapter.IAdapterResponse, err error) (*Payload, common.ElUpdate) {
	return nil, 0
}
func (t *TestStateMachineImplementation) eaf_Second(input interface{}, payload interface{}, ar adapter.IAdapterResponse, err error) (*Payload, common.ElUpdate) {
	return nil, 0
}
func (t *TestStateMachineImplementation) eap_Second(input interface{}, payload interface{}, ar adapter.IAdapterResponse, err error) (*Payload, common.ElUpdate) {
	return nil, 0
}

