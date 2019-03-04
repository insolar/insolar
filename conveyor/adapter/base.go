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
	"github.com/insolar/insolar/conveyor/interfaces/slot"
)

type idType = uint32

// PulseConveyorAdapterTaskSink is iface which helps to slot to push task to adapter
type PulseConveyorAdapterTaskSink interface {
	PushTask(respSink AdaptorToSlotResponseSink, elementID idType, handlerID idType, taskPayload interface{}) error
	CancelElementTasks(pulseNumber idType, elementID idType)
	CancelPulseTasks(pulseNumber uint32)
	FlushPulseTasks(pulseNumber uint32)
	FlushNodeTasks(nodeID idType)
}

// AdaptorToSlotResponseSink is iface which helps to adapter to access to slot
type AdaptorToSlotResponseSink interface {
	PushResponse(adapterID idType, elementID idType, handlerID idType, respPayload interface{})
	PushNestedEvent(adapterID idType, parentElementID idType, handlerID idType, eventPayload interface{})
	GetPulseNumber() uint32
	GetNodeID() uint32
	GetSlotDetails() slot.SlotDetails
}

// AdapterTask contains info for launch adapter task
type AdapterTask struct {
	respSink    AdaptorToSlotResponseSink
	elementID   idType
	handlerID   idType
	taskPayload interface{}
}

// AdapterResponse contains info with adapter response
type AdapterResponse struct {
	AdapterID   idType
	ElementID   idType
	HandlerID   idType
	RespPayload interface{}
}

// AdapterNestedEvent contains info with adapter nested event
type AdapterNestedEvent struct {
	AdapterID       idType
	ParentElementID idType
	HandlerID       idType
	EventPayload    interface{}
}
