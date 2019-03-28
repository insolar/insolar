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
	"fmt"

	"github.com/insolar/insolar/conveyor/adapter"
	"github.com/insolar/insolar/conveyor/fsm"
	"github.com/insolar/insolar/conveyor/interfaces/iadapter"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/ledger/artifactmanager"
)

// custom types
type Event struct{}
type Payload struct {
	reply insolar.Reply
	err   error
}
type TA1 string
type TAR string

// conveyor: state_machine
type GetCodeStateMachine interface {
	GetTypeID() fsm.ID

	initPresentHandler(input Event, payload interface{}) (*Payload, fsm.ElementState, error)
	initFutureHandler(input Event, payload interface{}, element fsm.SlotElementHelper) (*Payload, fsm.ElementState, error)

	errorPresentInit(input interface{}, payload interface{}, err error) (*Payload, fsm.ElementState)
	errorFutureInit(input interface{}, payload interface{}, err error) (*Payload, fsm.ElementState)

	// State Declaration
	stateFirst() fsm.StateID

	// Migration
	migrateFromPresentFirst(input Event, payload *Payload) (*Payload, error)

	// Transition
	transitPresentFirst(input Event, payload insolar.ConveyorPendingMessage, element fsm.SlotElementHelper) (*Payload, fsm.ElementState, error)

	// Adapter Response
	responsePresentFirst(input Event, payload artifactmanager.GetCodeResp, element fsm.SlotElementHelper) (*Payload, fsm.ElementState, error)

	// State Error
	errorPresentFirst(input interface{}, payload interface{}, err error) (*Payload, fsm.ElementState)

	// Adapter Response Error
	errorResponsePresentFirst(input interface{}, payload interface{}, ar interface{}, err error) (*Payload, fsm.ElementState)

	// State Declaration
	stateSecond() fsm.StateID

	// Migration
	migrateFromPresentSecond(input Event, payload *Payload) (*Payload, error)

	// Transition
	transitPresentSecond(input Event, payload *Payload /* todo: , adapterHelper TA1*/) (*Payload, fsm.ElementState, error)

	// Adapter Response
	responsePresentSecond(input Event, payload *Payload, respPayload artifactmanager.GetCodeResp) (*Payload, fsm.ElementState, error)

	// State Error
	errorPresentSecond(input interface{}, payload interface{}, err error) (*Payload, fsm.ElementState)

	// Adapter Response Error
	errorResponsePresentSecond(input interface{}, payload interface{}, ar interface{}, err error) (*Payload, fsm.ElementState)

	// State Declaration
	stateThird() fsm.StateID

	// Migration
	migrateFromPresentThird(input Event, payload *Payload) (*Payload, error)

	// Transition
	transitPresentThird(input Event, payload *Payload, element fsm.SlotElementHelper) (*Payload, fsm.ElementState, error)

	// Adapter Response
	responsePresentThird(input Event, payload interface{}) (*Payload, fsm.ElementState, error)

	// State Error
	errorPresentThird(input interface{}, payload interface{}, err error) (*Payload, fsm.ElementState)

	// Adapter Response Error
	errorResponsePresentThird(input interface{}, payload interface{}, ar interface{}, err error) (*Payload, fsm.ElementState)

	// State Declaration
	stateFourth() fsm.StateID

	// Migration
	migrateFromPresentFourth(input Event, payload *Payload) (*Payload, error)

	// Transition
	transitPresentFourth(input Event, payload *Payload /* todo: , adapterHelper TA1*/) (*Payload, fsm.ElementState, error)

	// Adapter Response
	responsePresentFourth(input Event, payload interface{}, element fsm.SlotElementHelper) (*Payload, fsm.ElementState, error)

	// State Error
	errorPresentFourth(input interface{}, payload interface{}, err error) (*Payload, fsm.ElementState)

	// Adapter Response Error
	errorResponsePresentFourth(input interface{}, payload interface{}, ar interface{}, err error) (*Payload, fsm.ElementState)
}

type CleanGetCodeStateMachine struct {
	BaseGetCodeStateMachine
}

func (sm *CleanGetCodeStateMachine) initPresentHandler(input Event, payload interface{}) (*Payload, fsm.ElementState, error) {
	p := &Payload{}
	return p, fsm.NewElementState(sm.GetTypeID(), sm.stateFirst()), nil
}

func (sm *CleanGetCodeStateMachine) initFutureHandler(input Event, payload interface{}, element fsm.SlotElementHelper) (*Payload, fsm.ElementState, error) {
	element.DeactivateTill(fsm.Response)
	p := &Payload{}
	return p, fsm.NewElementState(sm.GetTypeID(), sm.stateFirst()), nil
}

func (sm *CleanGetCodeStateMachine) errorPresentInit(input interface{}, payload interface{}, err error) (*Payload, fsm.ElementState) {
	panic("implement me")
}

func (sm *CleanGetCodeStateMachine) errorFutureInit(input interface{}, payload interface{}, err error) (*Payload, fsm.ElementState) {
	panic("implement me")
}

func (sm *CleanGetCodeStateMachine) migrateFromPresentFirst(input Event, payload *Payload) (*Payload, error) {
	return payload, nil
}

func (sm *CleanGetCodeStateMachine) transitPresentFirst(input Event, payload insolar.ConveyorPendingMessage, element fsm.SlotElementHelper) (*Payload, fsm.ElementState, error) {
	parcel := payload.Msg
	err := adapter.CurrentCatalog.GetCode.GetCode(element, parcel, 2)
	if err != nil {
		return nil, 0, nil
	}
	return nil, fsm.NewElementState(sm.GetTypeID(), sm.stateSecond()), nil
}

func (sm *CleanGetCodeStateMachine) responsePresentFirst(input Event, payload artifactmanager.GetCodeResp, element fsm.SlotElementHelper) (*Payload, fsm.ElementState, error) {
	panic("implement me")
}

func (sm *CleanGetCodeStateMachine) errorPresentFirst(input interface{}, payload interface{}, err error) (*Payload, fsm.ElementState) {
	panic("implement me")
}

func (sm *CleanGetCodeStateMachine) errorResponsePresentFirst(input interface{}, payload interface{}, ar iadapter.Response, err error) (*Payload, fsm.ElementState) {
	panic("implement me")
}

func (sm *CleanGetCodeStateMachine) migrateFromPresentSecond(input Event, payload *Payload) (*Payload, error) {
	return payload, nil
}

func (sm *CleanGetCodeStateMachine) transitPresentSecond(input Event, payload *Payload /* todo: , adapterHelper TA1*/) (*Payload, fsm.ElementState, error) {
	panic("implement me")
}

func (sm *CleanGetCodeStateMachine) responsePresentSecond(input Event, payload *Payload, respPayload artifactmanager.GetCodeResp) (*Payload, fsm.ElementState, error) {
	payload.err = respPayload.Err
	payload.reply = respPayload.Reply

	return payload, fsm.NewElementState(sm.GetTypeID(), sm.stateThird()), nil
}

func (sm *CleanGetCodeStateMachine) errorPresentSecond(input interface{}, payload interface{}, err error) (*Payload, fsm.ElementState) {
	panic("implement me")
}

func (sm *CleanGetCodeStateMachine) errorResponsePresentSecond(input interface{}, payload interface{}, ar iadapter.Response, err error) (*Payload, fsm.ElementState) {
	panic("implement me")
}

func (sm *CleanGetCodeStateMachine) migrateFromPresentThird(input Event, payload *Payload) (*Payload, error) {
	return payload, nil
}

func (sm *CleanGetCodeStateMachine) transitPresentThird(input Event, payload *Payload, element fsm.SlotElementHelper) (*Payload, fsm.ElementState, error) {
	var err error
	if payload.err != nil {
		// TODO: return error to future
		err = adapter.CurrentCatalog.SendResponse.SendResponse(element, nil, 4)
	} else {
		err = adapter.CurrentCatalog.SendResponse.SendResponse(element, payload.reply, 4)
	}
	if err != nil {
		return nil, 0, nil
	}
	return nil, fsm.NewElementState(sm.GetTypeID(), sm.stateFourth()), nil
}

func (sm *CleanGetCodeStateMachine) responsePresentThird(input Event, payload interface{}) (*Payload, fsm.ElementState, error) {
	panic("implement me")
}

func (sm *CleanGetCodeStateMachine) errorPresentThird(input interface{}, payload interface{}, err error) (*Payload, fsm.ElementState) {
	panic("implement me")
}

func (sm *CleanGetCodeStateMachine) errorResponsePresentThird(input interface{}, payload interface{}, ar iadapter.Response, err error) (*Payload, fsm.ElementState) {
	panic("implement me")
}

func (sm *CleanGetCodeStateMachine) migrateFromPresentFourth(input Event, payload *Payload) (*Payload, error) {
	return payload, nil
}

func (sm *CleanGetCodeStateMachine) transitPresentFourth(input Event, payload *Payload /* todo: , adapterHelper TA1*/) (*Payload, fsm.ElementState, error) {
	panic("implement me")
}

func (sm *CleanGetCodeStateMachine) responsePresentFourth(input Event, payload interface{}, element fsm.SlotElementHelper) (*Payload, fsm.ElementState, error) {
	switch res := payload.(type) {
	case string, error:
		return nil, 0, nil
	default:
		return nil, 0, fmt.Errorf("GetCode: unexpected reply: %T", res)
	}
}

func (sm *CleanGetCodeStateMachine) errorPresentFourth(input interface{}, payload interface{}, err error) (*Payload, fsm.ElementState) {
	panic("implement me")
}

func (sm *CleanGetCodeStateMachine) errorResponsePresentFourth(input interface{}, payload interface{}, ar iadapter.Response, err error) (*Payload, fsm.ElementState) {
	panic("implement me")
}
