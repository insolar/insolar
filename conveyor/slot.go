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

package conveyor

import (
	"github.com/insolar/insolar/conveyor/queue"
	"github.com/insolar/insolar/core"
)

// PulseState is the states of pulse inside slot
type PulseState int

//go:generate stringer -type=PulseState
const (
	Unallocated = PulseState(iota)
	Future
	Present
	Past
	Antique
)

// SlotState shows slot working mode
type SlotState int

//go:generate stringer -type=SlotState
const (
	Initializing = SlotState(iota)
	Working
	Suspending
)

const slotSize = 10000
const slotElementDelta = 1000000

type SlotDetails interface {
	getPulseNumber() core.PulseNumber
	getNodeId() uint32
	getPulseData() *core.Pulse
	getNodeData() interface{}
}

// HandlersConfiguration contains configuration of handlers for specific pulse state
// TODO: logic will be provided after pulse change mechanism
type HandlersConfiguration struct {
	state SlotState
}

// TODO: logic will be provided after pulse change mechanism
func (s *HandlersConfiguration) getMachineConfiguration(smType int) StateMachineType {
	return nil
}

// ActivationStatus represents status of work for slot element
type ActivationStatus int

//go:generate stringer -type=ActivationStatus
const (
	EmptyElement = ActivationStatus(iota)
	ActiveElement
	NotActiveElement
)

// ActivationStatusList is a list of slotElements by ActivationStatus with pointers to head and tail
type ActivationStatusList struct {
	head *slotElement
	tail *slotElement
}

// Slot holds info about specific pulse and events for it
type Slot struct {
	handlersConfiguration HandlersConfiguration
	inputQueue            queue.IQueue
	pulseState            PulseState
	slotState             SlotState
	pulse                 *core.Pulse
	pulseNumber           core.PulseNumber
	nodeId                uint32
	nodeData              interface{}
	elements              []slotElement
	elementListMap        map[ActivationStatus]ActivationStatusList
}

// SlotStateMachine represents state machine of slot itself
var SlotStateMachine = slotElement{
	id:               0,
	state:            0,
	stateMachineType: 0,
}

func initElementsBuf() []slotElement {
	elements := make([]slotElement, slotSize)
	elements[0] = SlotStateMachine
	var nextElement *slotElement
	for i := slotSize - 1; i > 0; i-- {
		elements[i] = *NewSlotElement(EmptyElement)
		elements[i].id = uint32(i)
		elements[i].listNext = nextElement
		nextElement = &elements[i]
	}
	return elements
}

// NewSlot creates new instance of Slot
func NewSlot(pulseState PulseState, pulseNumber core.PulseNumber) *Slot {
	slotState := Initializing
	if pulseState == Antique {
		slotState = Working
	}

	elements := initElementsBuf()

	empty := ActivationStatusList{
		// elements[0] contains SlotStateMachine, so first empty element is elements[1]
		head: &elements[1],
		tail: &elements[slotSize-1],
	}

	elementListMap := map[ActivationStatus]ActivationStatusList{
		EmptyElement:     empty,
		ActiveElement:    ActivationStatusList{},
		NotActiveElement: ActivationStatusList{},
	}
	return &Slot{
		pulseState:     pulseState,
		inputQueue:     queue.NewMutexQueue(),
		pulseNumber:    pulseNumber,
		slotState:      slotState,
		elements:       elements,
		elementListMap: elementListMap,
	}
}

func (s *Slot) getPulseNumber() core.PulseNumber {
	return s.pulseNumber
}

func (s *Slot) getPulseData() *core.Pulse {
	return s.pulse
}

func (s *Slot) getNodeId() uint32 {
	return s.nodeId
}

func (s *Slot) getNodeData() interface{} {
	return s.nodeData
}

// createElement creates new active element from empty element
func (s *Slot) createElement(stateMachineType StateMachineType, state uint16, event queue.OutputElement) *slotElement {
	element := s.getElement(EmptyElement)
	element.id = element.id + slotElementDelta
	element.stateMachineType = stateMachineType
	element.state = state
	element.activationStatus = ActiveElement
	// Set other fields to element, like:
	// element.payload = event.GetPayload()

	s.addElement(ActiveElement, element)
	return element
}

// getElement gets element of provided status from correspondent linked list (and remove it from that list)
func (s *Slot) getElement(status ActivationStatus) *slotElement {
	list := s.elementListMap[status]
	result := list.head
	if result == nil {
		return nil
	}
	list.head = list.head.listNext
	return result
}

// addElement adds element of provided status to correspondent linked list
func (s *Slot) addElement(status ActivationStatus, element *slotElement) {
	list := s.elementListMap[status]
	if list.head == nil {
		list.head = element
		list.tail = element
	} else {
		list.tail.listNext = element
		list.tail = element
	}
}

type StateMachineType interface{}

type slotElement struct {
	id               uint32
	payload          interface{}
	state            uint16
	stateMachineType StateMachineType

	listNext         *slotElement
	activationStatus ActivationStatus
}

func NewSlotElement(activationStatus ActivationStatus) *slotElement {
	return &slotElement{activationStatus: activationStatus}
}
