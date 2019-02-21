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

// SlotDetails provides information about slot
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

// ElementList is a list of slotElements with pointers to head and tail
type ElementList struct {
	head *slotElement
	tail *slotElement
}

// popElement gets element from linked list (and remove it from list)
func (l *ElementList) popElement() *slotElement {
	result := l.head
	if result == nil {
		return nil
	}
	l.head = l.head.listNext
	return result
}

// pushElement adds element to linked list
func (l *ElementList) pushElement(element *slotElement) {
	if l.head == nil {
		l.head = element
		l.tail = element
	} else {
		l.tail.listNext = element
		l.tail = element
	}
}

func (l *ElementList) len() int {
	i := 0
	for element := l.head; element != nil; element = element.listNext {
		i++
	}
	return i
}

// Slot holds info about specific pulse and events for it
type Slot struct {
	handlersConfiguration HandlersConfiguration
	inputQueue            queue.IQueue
	pulseState            PulseState
	slotState             SlotState
	pulse                 core.Pulse
	pulseNumber           core.PulseNumber
	nodeID                uint32
	nodeData              interface{}
	elements              []slotElement
	elementListMap        map[ActivationStatus]*ElementList
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
		elements[i] = *newSlotElement(EmptyElement)
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

	empty := &ElementList{
		// elements[0] contains SlotStateMachine, so first empty element is elements[1]
		head: &elements[1],
		tail: &elements[slotSize-1],
	}

	elementListMap := map[ActivationStatus]*ElementList{
		EmptyElement:     empty,
		ActiveElement:    &ElementList{},
		NotActiveElement: &ElementList{},
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

func (s *Slot) getPulseData() core.Pulse {
	return s.pulse
}

func (s *Slot) getNodeID() uint32 {
	return s.nodeID
}

func (s *Slot) getNodeData() interface{} {
	return s.nodeData
}

// createElement creates new active element from empty element
func (s *Slot) createElement(stateMachineType StateMachineType, state uint16, event queue.OutputElement) *slotElement {
	element := s.popElement(EmptyElement)
	element.stateMachineType = stateMachineType
	element.state = state
	element.activationStatus = ActiveElement
	element.listNext = nil
	// Set other fields to element, like:
	// element.payload = event.GetPayload()

	s.pushElement(ActiveElement, element)
	return element
}

// popElement gets element of provided status from correspondent linked list (and remove it from that list)
func (s *Slot) popElement(status ActivationStatus) *slotElement {
	list := s.elementListMap[status]
	return list.popElement()
}

// pushElement adds element of provided status to correspondent linked list
func (s *Slot) pushElement(status ActivationStatus, element *slotElement) {
	element.activationStatus = status
	if status == EmptyElement {
		element.id = element.id + slotElementDelta
	}
	list := s.elementListMap[status]
	list.pushElement(element)
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

// newSlotElement creates new slot element with provided activation status
func newSlotElement(activationStatus ActivationStatus) *slotElement {
	return &slotElement{activationStatus: activationStatus}
}
