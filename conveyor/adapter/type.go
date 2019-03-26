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
	"github.com/insolar/insolar/conveyor/interfaces/iadapter"
	"github.com/insolar/insolar/conveyor/interfaces/slot"
	"github.com/insolar/insolar/conveyor/queue"
	"github.com/insolar/insolar/insolar"
)

type idType = uint32

// PulseConveyorAdapterTaskSink is iface which helps to slot to push task to adapter
// NestedEvent gives access to nested event of adapter
//go:generate minimock -i github.com/insolar/insolar/conveyor/adapter.PulseConveyorAdapterTaskSink -o ./ -s _mock.go
type PulseConveyorAdapterTaskSink interface {
	PushTask(respSink AdapterToSlotResponseSink, elementID idType, handlerID idType, taskPayload interface{}) error
	CancelElementTasks(pulseNumber insolar.PulseNumber, elementID idType)
	CancelPulseTasks(pulseNumber insolar.PulseNumber)
	FlushPulseTasks(pulseNumber insolar.PulseNumber)
	FlushNodeTasks(nodeID idType)
	GetAdapterID() uint32
}

// AdapterToSlotResponseSink is iface which helps to adapter to access to slot
type AdapterToSlotResponseSink interface {
	PushResponse(adapterID uint32, elementID uint32, handlerID uint32, respPayload interface{})
	PushNestedEvent(adapterID uint32, parentElementID uint32, handlerID uint32, eventPayload interface{})
	GetPulseNumber() insolar.PulseNumber
	GetNodeID() uint32
	GetSlotDetails() slot.SlotDetails
}

// AdapterTask contains info for launch adapter task
type AdapterTask struct {
	respSink    AdapterToSlotResponseSink
	elementID   idType
	handlerID   idType
	TaskPayload interface{}
}

// AdapterResponse contains info with adapter response
type AdapterResponse struct {
	adapterID   idType
	elementID   idType
	handlerID   idType
	respPayload interface{}
}

// NewAdapterResponse creates new adapter response
func NewAdapterResponse(adapterID idType, elementID idType, handlerID idType, respPayload interface{}) iadapter.Response {
	return &AdapterResponse{
		adapterID:   adapterID,
		elementID:   elementID,
		handlerID:   handlerID,
		respPayload: respPayload,
	}
}

// GetAdapterID implements Response method
func (ar *AdapterResponse) GetAdapterID() uint32 {
	return ar.adapterID
}

// GetElementID implements Response method
func (ar *AdapterResponse) GetElementID() uint32 {
	return ar.elementID
}

// GetHandlerID implements Response method
func (ar *AdapterResponse) GetHandlerID() uint32 {
	return ar.handlerID
}

// GetRespPayload implements Response method
func (ar *AdapterResponse) GetRespPayload() interface{} {
	return ar.respPayload
}

// AdapterNestedEvent contains info with adapter nested event
type AdapterNestedEvent struct {
	adapterID       idType
	parentElementID idType
	handlerID       idType
	eventPayload    interface{}
}

func NewAdapterNestedEvent(adapterID uint32, parentElementID uint32, handlerID uint32, eventPayload interface{}) iadapter.NestedEvent {
	return &AdapterNestedEvent{
		adapterID:       adapterID,
		parentElementID: parentElementID,
		handlerID:       handlerID,
		eventPayload:    eventPayload,
	}
}

// GetHandlerID implements NestedEvent method
func (a *AdapterNestedEvent) GetAdapterID() uint32 {
	return a.adapterID
}

// GetParentElementID implements NestedEvent method
func (a *AdapterNestedEvent) GetParentElementID() uint32 {
	return a.parentElementID
}

// GetHandlerID implements NestedEvent method
func (a *AdapterNestedEvent) GetHandlerID() uint32 {
	return a.handlerID
}

// GetEventPayload implements NestedEvent method
func (a *AdapterNestedEvent) GetEventPayload() interface{} {
	return a.eventPayload
}

// CancelInfo provides info about cancellation
type CancelInfo interface {
	Cancel() <-chan bool
	Flush() <-chan bool
	IsCanceled() bool
	IsFlushed() bool
	ID() uint64
}

// NestedEventHelper is helper for sending nested event from Processor
type NestedEventHelper interface {
	Send(eventPayload interface{})
}

// Processor is iface for processing task for adapter
type Processor interface {
	Process(task AdapterTask, nestedEventHelper NestedEventHelper, cancelInfo CancelInfo) interface{}
}

// NewAdapterWithQueue creates new instance of Adapter
func NewAdapterWithQueue(processor Processor, id idType) PulseConveyorAdapterTaskSink {
	adapter := &CancellableQueueAdapter{
		queue:             queue.NewMutexQueue(),
		processingStarted: 0,
		stopProcessing:    0,
		processingStopped: make(chan bool, 1),
		taskHolder:        newTaskHolder(),
		processor:         processor,
		adapterID:         id,
	}
	started := make(chan bool, 1)
	go adapter.StartProcessing(started)
	<-started

	return adapter
}
