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
	"github.com/insolar/insolar/conveyor/adapter/adapterid"
	"github.com/insolar/insolar/conveyor/queue"
	"github.com/insolar/insolar/insolar"
)

// TaskSink is iface which helps to slot to push task to adapter
//go:generate minimock -i github.com/insolar/insolar/conveyor/adapter.TaskSink -o ./ -s _mock.go
type TaskSink interface {
	PushTask(respSink ResponseSink, elementID uint32, handlerID uint32, taskPayload interface{}) error
	CancelElementTasks(pulseNumber insolar.PulseNumber, elementID uint32)
	CancelPulseTasks(pulseNumber insolar.PulseNumber)
	FlushPulseTasks(pulseNumber insolar.PulseNumber)
	FlushNodeTasks(nodeID uint32)
	GetAdapterID() adapterid.ID
}

// ResponseSink is iface which helps to adapter to access to slot
type ResponseSink interface {
	PushResponse(adapterID adapterid.ID, elementID uint32, handlerID uint32, respPayload interface{})
	PushNestedEvent(adapterID adapterid.ID, parentElementID uint32, handlerID uint32, eventPayload interface{})
	GetPulseNumber() insolar.PulseNumber
	GetNodeID() uint32
	GetSlotDetails() SlotDetails
}

// Task contains info for launch adapter task
type Task struct {
	respSink    ResponseSink
	elementID   uint32
	handlerID   uint32
	TaskPayload interface{}
}

// Response contains info with adapter response
type Response struct {
	adapterID   adapterid.ID
	elementID   uint32
	handlerID   uint32
	respPayload interface{}
}

// NewResponse creates new adapter response
func NewResponse(adapterID adapterid.ID, elementID uint32, handlerID uint32, respPayload interface{}) *Response {
	return &Response{
		adapterID:   adapterID,
		elementID:   elementID,
		handlerID:   handlerID,
		respPayload: respPayload,
	}
}

// GetAdapterID implements Response method
func (ar *Response) GetAdapterID() adapterid.ID {
	return ar.adapterID
}

// GetElementID implements Response method
func (ar *Response) GetElementID() uint32 {
	return ar.elementID
}

// GetHandlerID implements Response method
func (ar *Response) GetHandlerID() uint32 {
	return ar.handlerID
}

// GetRespPayload implements Response method
func (ar *Response) GetRespPayload() interface{} {
	return ar.respPayload
}

// NestedEvent contains info with adapter nested event
type NestedEvent struct {
	adapterID       adapterid.ID
	parentElementID uint32
	handlerID       uint32
	eventPayload    interface{}
}

func NewNestedEvent(adapterID adapterid.ID, parentElementID uint32, handlerID uint32, eventPayload interface{}) *NestedEvent {
	return &NestedEvent{
		adapterID:       adapterID,
		parentElementID: parentElementID,
		handlerID:       handlerID,
		eventPayload:    eventPayload,
	}
}

// GetHandlerID implements NestedEvent method
func (a *NestedEvent) GetAdapterID() adapterid.ID {
	return a.adapterID
}

// GetParentElementID implements NestedEvent method
func (a *NestedEvent) GetParentElementID() uint32 {
	return a.parentElementID
}

// GetHandlerID implements NestedEvent method
func (a *NestedEvent) GetHandlerID() uint32 {
	return a.handlerID
}

// GetEventPayload implements NestedEvent method
func (a *NestedEvent) GetEventPayload() interface{} {
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
	Process(task Task, nestedEventHelper NestedEventHelper, cancelInfo CancelInfo) interface{}
}

// NewAdapterWithQueue creates new instance of Adapter
func NewAdapterWithQueue(processor Processor, id adapterid.ID) TaskSink {
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

// SlotDetails provides information about slot
//go:generate minimock -i github.com/insolar/insolar/conveyor/adapter.SlotDetails -o ./ -s _mock.go
type SlotDetails interface {
	GetPulseNumber() insolar.PulseNumber // nolint: unused
	GetNodeID() uint32                   // nolint: unused
	GetPulseData() insolar.Pulse         // nolint: unused
	GetNodeData() interface{}            // nolint: unused
}
