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

package slot

import (
	"fmt"

	"github.com/insolar/insolar/conveyor/adapter"
	"github.com/insolar/insolar/conveyor/adapter/adapterid"
	"github.com/insolar/insolar/conveyor/fsm"
	"github.com/insolar/insolar/conveyor/generator/matrix"
	"github.com/insolar/insolar/conveyor/queue"
	"github.com/insolar/insolar/insolar"

	"github.com/pkg/errors"
)

// HandlerStorage gives access to handlers
// TODO: make HandlerStorage not available for export. For now we use it in conveyor_integration_test only
var HandlerStorage matrix.StateMachineHolder

func init() {
	HandlerStorage = matrix.NewMatrix()
}

// State shows slot working mode
type State uint32

//go:generate stringer -type=State
const (
	Initializing = State(iota)
	Working
	Suspending
	Canceling
)

const slotSize = 10000
const slotElementDelta = slotSize // nolint: unused

// RemoveSlotCallback allows to remove slot by pulse number
type RemoveSlotCallback func(number insolar.PulseNumber)

// HandlersConfiguration contains configuration of handlers for specific pulse state
// TODO: logic will be provided after pulse change mechanism
type HandlersConfiguration struct {
	pulseStateMachines matrix.SetAccessor
	initStateMachine   matrix.StateMachine
}

// TODO: logic will be provided after pulse change mechanism
func (h *HandlersConfiguration) getMachineConfiguration(smType int) matrix.StateMachine { // nolint: unused
	return nil
}

// elementList is a list of slotElements with pointers to head and tail
type elementList struct {
	head   *slotElement
	tail   *slotElement
	length int
}

func (l *elementList) isEmpty() bool {
	return l.head == nil
}

// popElement gets element from linked list (and remove it from list)
func (l *elementList) popElement() *slotElement {
	result := l.head
	if result == nil {
		return nil
	}
	l.removeElement(result)
	return result
}

// removeElement removes element from linked list
func (l *elementList) removeElement(element *slotElement) { // nolint: unused
	if element == nil {
		return
	}
	next := element.nextElement
	prev := element.prevElement
	if prev != nil {
		prev.nextElement = next
	} else {
		l.head = next
	}
	if next != nil {
		next.prevElement = prev
	} else {
		l.tail = prev
	}
	element.prevElement = nil
	element.nextElement = nil
	l.length--
}

// pushElement adds element to linked list
func (l *elementList) pushElement(element *slotElement) { // nolint: unused
	if l.head == nil {
		l.head = element
	} else {
		l.tail.nextElement = element
		element.prevElement = l.tail
	}
	element.nextElement = nil
	l.tail = element
	l.length++
}

func (l *elementList) len() int { // nolint: unused
	return l.length
}

// slot holds info about specific pulse and events for it
type slot struct {
	handlersConfiguration HandlersConfiguration // nolint
	inputQueue            queue.Queue
	responseQueue         queue.Queue
	pulseState            PulseState
	slotState             State
	stateMachine          slotElement
	pulse                 insolar.Pulse
	pulseNumber           insolar.PulseNumber
	nodeID                uint32
	nodeData              interface{}
	elements              []slotElement
	// we can use slice or just several fields of elementList, it will be faster but not pretty
	elementListMap     map[ActivationStatus]*elementList
	removeSlotCallback RemoveSlotCallback
}

func (s *slot) SinkPush(data interface{}) error {
	return s.inputQueue.SinkPush(data)
}

// TODO: fix me sweetie, this must not exist
func (s *slot) SetPulse(pulse insolar.Pulse) {
	s.pulse = pulse
}

func (s *slot) SinkPushAll(data []interface{}) error {
	return s.inputQueue.SinkPushAll(data)
}

func (s *slot) PushSignal(signalType uint32, callback queue.SyncDone) error {
	return s.inputQueue.PushSignal(signalType, callback)
}

// slotStateMachine represents state machine of slot itself
var slotStateMachine = slotElement{
	id:           0,
	state:        0,
	stateMachine: nil, // TODO: add smth correct
}

func initElementsBuf() ([]slotElement, *elementList) {
	elements := make([]slotElement, slotSize)
	emptyList := &elementList{}
	for i := 0; i < slotSize; i++ {
		// we don't have *slot here yet. Set it later
		elements[i] = *newSlotElement(EmptyElement, nil)
		elements[i].id = uint32(i)
		emptyList.pushElement(&elements[i])
	}
	return elements, emptyList
}

// NewWorkingSlot creates new instance of slot by TaskPusher interface
func NewWorkingSlot(pulseState PulseState, pulseNumber insolar.PulseNumber, removeSlotCallback RemoveSlotCallback) TaskPusher {

	slot := newSlot(pulseState, pulseNumber, removeSlotCallback)
	slot.runWorker()

	return slot
}

func newSlot(pulseState PulseState, pulseNumber insolar.PulseNumber, removeSlotCallback RemoveSlotCallback) *slot {
	slotState := Initializing
	if pulseState == Antique {
		slotState = Working
	}

	elements, emptyList := initElementsBuf()

	elementListMap := map[ActivationStatus]*elementList{
		EmptyElement:     emptyList,
		ActiveElement:    {},
		NotActiveElement: {},
	}

	slot := &slot{
		pulseState:         pulseState,
		inputQueue:         queue.NewMutexQueue(),
		responseQueue:      queue.NewMutexQueue(),
		pulseNumber:        pulseNumber,
		slotState:          slotState,
		stateMachine:       slotStateMachine,
		elements:           elements,
		elementListMap:     elementListMap,
		removeSlotCallback: removeSlotCallback,
		handlersConfiguration: HandlersConfiguration{
			initStateMachine: HandlerStorage.GetInitialStateMachine(),
		},
	}

	for i := range slot.elements {
		slot.elements[i].slot = slot
	}

	return slot
}

func (s *slot) runWorker() {
	worker := newWorker(s)
	go worker.run()
}

func (s *slot) PushResponse(adapterID adapterid.ID, elementID uint32, handlerID uint32, respPayload interface{}) {
	response := adapter.NewResponse(adapterID, elementID, handlerID, respPayload)
	err := s.responseQueue.SinkPush(response)
	if err != nil {
		panic("[ PushResponse ] Can't SinkPush: " + err.Error())
	}
}

func (s *slot) PushNestedEvent(adapterID adapterid.ID, parentElementID uint32, handlerID uint32, eventPayload interface{}) {
	event := adapter.NewNestedEvent(adapterID, parentElementID, handlerID, eventPayload)
	err := s.responseQueue.SinkPush(event)
	if err != nil {
		panic("[ PushNestedEvent ] Can't SinkPush: " + err.Error())
	}
}

func (s *slot) GetSlotDetails() adapter.SlotDetails {
	return s
}

// GetPulseNumber implements iface SlotDetails
func (s *slot) GetPulseNumber() insolar.PulseNumber { // nolint: unused
	return s.pulseNumber
}

// GetPulseData implements iface SlotDetails
func (s *slot) GetPulseData() insolar.Pulse { // nolint: unused
	return s.pulse
}

// GetNodeID implements iface SlotDetails
func (s *slot) GetNodeID() uint32 { // nolint: unused
	return s.nodeID
}

// GetNodeData implements iface SlotDetails
func (s *slot) GetNodeData() interface{} { // nolint: unused
	return s.nodeData
}

// createElement creates new active element from empty element
func (s *slot) createElement(stateMachine matrix.StateMachine, state fsm.StateID, event queue.OutputElement) (*slotElement, error) { // nolint: unused
	element := s.popElement(EmptyElement)
	element.stateMachine = stateMachine
	element.state = state
	element.activationStatus = ActiveElement
	element.nextElement = nil
	// TODO:  Set other fields to element, like:
	conveyorMsg, ok := event.GetData().(insolar.ConveyorPendingMessage)
	if !ok {
		return nil, errors.Errorf("[ createElement ]Input event must be 'insolar.ConveyorPendingMessage'. Actual: %+v", event.GetData())
	}
	element.payload = nil
	element.inputEvent = conveyorMsg.Msg
	element.responseFuture = conveyorMsg.Future

	err := s.pushElement(element)
	if err != nil {
		emptyList := s.elementListMap[EmptyElement]
		emptyList.pushElement(element)
		return nil, errors.Wrap(err, "[ createElement ]")
	}
	return element, nil
}

func (s *slot) hasExpired() bool {
	// TODO: This is used to delete past slot, which doesn't have elements and not active for some configure time
	return s.len(ActiveElement) == 0 && s.len(NotActiveElement) == 0
}

func (s *slot) hasElements(status ActivationStatus) bool {
	list, ok := s.elementListMap[status]
	if !ok {
		return false
	}
	return !list.isEmpty()
}

func (s *slot) isSuspending() bool {
	return s.slotState == Suspending
}

func (s *slot) isWorking() bool {
	return s.slotState == Working
}

// nolint: unused
func (s *slot) isInitializing() bool {
	return s.slotState == Initializing
}

// popElement gets element of provided status from correspondent linked list (and remove it from that list)
func (s *slot) popElement(status ActivationStatus) *slotElement { // nolint: unused
	list, ok := s.elementListMap[status]
	if !ok {
		return nil
	}
	return list.popElement()
}

func (s *slot) len(status ActivationStatus) int { // nolint: unused
	list, ok := s.elementListMap[status]
	if !ok {
		return 0
	}
	return list.len()
}

func (s *slot) extractSlotElementByID(id uint32) *slotElement { // nolint: unused
	element := &s.elements[id%slotSize]
	if element.id != id {
		return nil
	}

	list, ok := s.elementListMap[element.activationStatus]
	if ok {
		list.removeElement(element)
	}
	return element
}

// pushElement adds element of provided status to correspondent linked list
func (s *slot) pushElement(element *slotElement) error { // nolint: unused
	status := element.activationStatus
	list, ok := s.elementListMap[status]
	if !ok {
		return fmt.Errorf("[ pushElement ] can't push element: list for status %s doesn't exist", status)
	}
	if status == EmptyElement {
		oldID := element.id
		*element = *newSlotElement(EmptyElement, s)
		element.id = oldID + slotElementDelta
	}
	list.pushElement(element)
	return nil
}
