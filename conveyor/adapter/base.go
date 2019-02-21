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
	"github.com/insolar/insolar/core"
)

type idType = uint32

// PulseConveyorAdapterTaskSink is iface which helps to slot to push task to adapter
type PulseConveyorAdapterTaskSink interface {
	PushTask(respSink PulseConveyorSlotResponseSink, elementId idType, handlerId idType, taskPayload interface{}) error
	CancelElementTasks(pulseNumber idType, elementId idType)
	CancelPulseTasks(pulseNumber idType)
	FlushPulseTasks(pulseNumber uint32)
	FlushNodeTasks(nodeId idType)
}

// PulseConveyorSlotResponseSink is iface which helps to adapter to access to slot
type PulseConveyorSlotResponseSink interface {
	PushResponse(adapterId idType, elementId idType, handlerId idType, respPayload interface{})
	PushNestedEvent(adapterId idType, parentElementId idType, handlerId idType, eventPayload interface{})
	GetPulseNumber() uint32
	GetNodeId() uint32
	GetSlotDetails() SlotDetails
}

// SlotDetails holds info about slot
type SlotDetails struct {
}

// GetPulseNumber returns pulse number
func (sd *SlotDetails) GetPulseNumber() uint32 {
	return 0
}

// GetNodeId returns consensus's node id
func (sd *SlotDetails) GetNodeId() uint32 {
	return 32
}

// GetPulseData returns pulse data
func (sd *SlotDetails) GetPulseData() *core.Pulse {
	return core.GenesisPulse
}

// AdapterTask contains info for launch adapter task
type AdapterTask struct {
	respSink    PulseConveyorSlotResponseSink
	elementID   idType
	handlerID   idType
	taskPayload interface{}
}

// AdapterResponse contains info with adapter response
type AdapterResponse struct {
	adapterID   idType
	elementID   idType
	handlerID   idType
	respPayload interface{}
}

// AdapterNestedEvent contains info with adapter nested event
type AdapterNestedEvent struct {
	adapterID       idType
	parentElementID idType
	handlerID       idType
	eventPayload    interface{}
}
