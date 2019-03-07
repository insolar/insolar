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
	"fmt"

	"github.com/insolar/insolar/conveyor/interfaces/constant"
	"github.com/insolar/insolar/conveyor/interfaces/statemachine"
	"github.com/insolar/insolar/conveyor/queue"
	"github.com/insolar/insolar/core"
	"github.com/pkg/errors"
)

// SlotState shows slot working mode
type SlotState uint32

//go:generate stringer -type=SlotState
const (
	Initializing = SlotState(iota)
	Working
	Suspending
)

const slotSize = 10000
const slotElementDelta = 1000000 // nolint: unused

// HandlersConfiguration contains configuration of handlers for specific pulse state
// TODO: logic will be provided after pulse change mechanism
type HandlersConfiguration struct {
	state SlotState // nolint: unused
}

// TODO: logic will be provided after pulse change mechanism
func (s *HandlersConfiguration) getMachineConfiguration(smType int) statemachine.StateMachineType { // nolint: unused
	return nil
}

// ElementList is a list of slotElements with pointers to head and tail
type ElementList struct {
	head *slotElement
	tail *slotElement
}

func (l *ElementList) isEmpty() bool {
	return l.head == nil
}

// popElement gets element from linked list (and remove it from list)
func (l *ElementList) popElement() *slotElement {
	result := l.head
	if result == nil {
		return nil
	}
	l.head = l.head.nextElement
	return result
}

// pushElement adds element to linked list
func (l *ElementList) pushElement(element *slotElement) { // nolint: unused
	if l.head == nil {
		l.head = element
	} else {
		l.tail.nextElement = element
	}
	l.tail = element
}

// Slot holds info about specific pulse and events for it
type Slot struct {
	handlersConfiguration HandlersConfiguration // nolint
	inputQueue            queue.IQueue
	responseQueue         queue.IQueue
	pulseState            constant.PulseState
	slotState             SlotState
	stateMachine          slotElement
	pulse                 core.Pulse
	pulseNumber           core.PulseNumber
	nodeID                uint32
	nodeData              interface{}
	elements              []slotElement
	// we can use slice or just several fields of ElementList, it will be faster but not pretty
	elementListMap map[ActivationStatus]*ElementList
}

// SlotStateMachine represents state machine of slot itself
var SlotStateMachine = slotElement{
	id:               0,
	state:            0,
	stateMachineType: nil, // TODO: add smth correct
}

func initElementsBuf() []slotElement {
	elements := make([]slotElement, slotSize)
	var nextElement *slotElement
	for i := slotSize - 1; i >= 0; i-- {
		elements[i] = *newSlotElement(EmptyElement)
		elements[i].id = uint32(i)
		elements[i].nextElement = nextElement
		nextElement = &elements[i]
	}
	return elements
}

// NewSlot creates new instance of Slot
func NewSlot(pulseState constant.PulseState, pulseNumber core.PulseNumber) *Slot {
	slotState := Initializing
	if pulseState == constant.Antique {
		slotState = Working
	}

	elements := initElementsBuf()

	elementListMap := map[ActivationStatus]*ElementList{
		EmptyElement: {
			head: &elements[0],
			tail: &elements[slotSize-1],
		},
		ActiveElement:    {},
		NotActiveElement: {},
	}
	return &Slot{
		pulseState:     pulseState,
		inputQueue:     queue.NewMutexQueue(),
		responseQueue:  queue.NewMutexQueue(),
		pulseNumber:    pulseNumber,
		slotState:      slotState,
		stateMachine:   SlotStateMachine,
		elements:       elements,
		elementListMap: elementListMap,
	}
}

// GetPulseNumber implements iface SlotDetails
func (s *Slot) GetPulseNumber() core.PulseNumber { // nolint: unused
	return s.pulseNumber
}

// GetPulseData implements iface SlotDetails
func (s *Slot) GetPulseData() core.Pulse { // nolint: unused
	return s.pulse
}

// GetNodeID implements iface SlotDetails
func (s *Slot) GetNodeID() uint32 { // nolint: unused
	return s.nodeID
}

// GetNodeData implements iface SlotDetails
func (s *Slot) GetNodeData() interface{} { // nolint: unused
	return s.nodeData
}

// createElement creates new active element from empty element
func (s *Slot) createElement(stateMachineType statemachine.StateMachineType, state uint32, event queue.OutputElement) (*slotElement, error) { // nolint: unused
	element := s.popElement(EmptyElement)
	element.stateMachineType = stateMachineType
	element.state = state
	element.activationStatus = ActiveElement
	element.nextElement = nil
	// TODO:  Set other fields to element, like:
	element.payload = event.GetData()

	err := s.pushElement(element)
	if err != nil {
		emptyList := s.elementListMap[EmptyElement]
		emptyList.pushElement(element)
		return nil, errors.Wrap(err, "[ createElement ]")
	}
	return element, nil
}

func (s *Slot) hasExpired() bool {
	// TODO: This is used to delete past slot, which doesn't have elements and not active for some configure time
	return false
}

func (s *Slot) hasElements(status ActivationStatus) bool {
	list, ok := s.elementListMap[status]
	if !ok {
		return false
	}
	return !list.isEmpty()
}

func (s *Slot) isSuspending() bool {
	return s.slotState == Suspending
}

func (s *Slot) isWorking() bool {
	return s.slotState == Working
}

func (s *Slot) isInitializing() bool {
	return s.slotState == Initializing
}

// popElement gets element of provided status from correspondent linked list (and remove it from that list)
func (s *Slot) popElement(status ActivationStatus) *slotElement { // nolint: unused
	list, ok := s.elementListMap[status]
	if !ok {
		return nil
	}
	return list.popElement()
}

func (s *Slot) getSlotElementByID(id uint32) *slotElement {
	return &s.elements[id%slotSize]
}

// pushElement adds element of provided status to correspondent linked list
func (s *Slot) pushElement(element *slotElement) error { // nolint: unused
	status := element.activationStatus
	list, ok := s.elementListMap[status]
	if !ok {
		return fmt.Errorf("[ pushElement ] can't push element: list for status %s doesn't exist", status)
	}
	if status == EmptyElement {
		oldID := element.id
		*element = *newSlotElement(EmptyElement)
		element.id = oldID + slotElementDelta
	}
	list.pushElement(element)
	return nil
}
