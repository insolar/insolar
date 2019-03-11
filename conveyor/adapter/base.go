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

package adapter

import (
	"github.com/insolar/insolar/conveyor/interfaces/islot"
)

type idType = uint32

// PulseConveyorAdapterTaskSink is iface which helps to islot to push task to iadapter
type PulseConveyorAdapterTaskSink interface {
	PushTask(respSink AdaptorToSlotResponseSink, elementID idType, handlerID idType, taskPayload interface{}) error
	CancelElementTasks(pulseNumber idType, elementID idType)
	CancelPulseTasks(pulseNumber uint32)
	FlushPulseTasks(pulseNumber uint32)
	FlushNodeTasks(nodeID idType)
}

// AdaptorToSlotResponseSink is iface which helps to iadapter to access to islot
type AdaptorToSlotResponseSink interface {
	PushResponse(adapterID idType, elementID idType, handlerID idType, respPayload interface{})
	PushNestedEvent(adapterID idType, parentElementID idType, handlerID idType, eventPayload interface{})
	GetPulseNumber() uint32
	GetNodeID() uint32
	GetSlotDetails() islot.SlotDetails
}

// AdapterTask contains info for launch iadapter task
type AdapterTask struct {
	respSink    AdaptorToSlotResponseSink
	elementID   idType
	handlerID   idType
	taskPayload interface{}
}

// AdapterResponse contains info with iadapter response
type AdapterResponse struct {
	adapterID   idType
	elementID   idType
	handlerID   idType
	respPayload interface{}
}

func (ar *AdapterResponse) SetElementID(id idType) {
	ar.elementID = id
}

func (ar *AdapterResponse) GetAdapterID() uint32 {
	return ar.adapterID
}

func (ar *AdapterResponse) GetElementID() uint32 {
	return ar.elementID
}

func (ar *AdapterResponse) GetHandlerID() uint32 {
	return ar.handlerID
}

func (ar *AdapterResponse) GetRespPayload() interface{} {
	return ar.respPayload
}

// AdapterNestedEvent contains info with iadapter nested event
type AdapterNestedEvent struct {
	AdapterID       idType
	ParentElementID idType
	HandlerID       idType
	EventPayload    interface{}
}
