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
	"testing"

	"github.com/insolar/insolar/conveyor/interfaces/statemachine"
	"github.com/insolar/insolar/conveyor/queue"
	"github.com/insolar/insolar/core"
	"github.com/stretchr/testify/require"
)

func testSlot(t *testing.T, isQueueOk bool, pulseNumber core.PulseNumber) *Slot {
	var q *queue.IQueueMock
	if isQueueOk {
		q = mockQueue(t)
	} else {
		q = mockQueueReturnFalse(t)
	}

	return &Slot{
		inputQueue:  q,
		pulseNumber: pulseNumber,
		pulse:       core.Pulse{PulseNumber: pulseNumber},
	}
}

func len3List() ElementList {
	el1 := &slotElement{id: 1}
	el2 := &slotElement{id: 2}
	el3 := &slotElement{id: 3}
	el1.nextElement = el2
	el2.nextElement = el3

	l := ElementList{
		head: el1,
		tail: el3,
	}
	return l
}

func elementListLength(l *ElementList) int {
	i := 0
	for element := l.head; element != nil; element = element.nextElement {
		i++
	}
	return i
}

func TestElementList_popElement_FromEmptyList(t *testing.T) {
	l := ElementList{}
	el := l.popElement()
	require.Nil(t, el)
}

func TestElementList_popElement_FromLenOneList(t *testing.T) {
	expectedElement := &slotElement{id: 1}
	l := ElementList{head: expectedElement, tail: expectedElement}

	el := l.popElement()
	require.Equal(t, expectedElement, el)
	require.Equal(t, 0, elementListLength(&l))
}

func TestElementList_popElement_Multiple(t *testing.T) {
	l := ElementList{}
	numElements := 333
	prevElement := &slotElement{id: uint32(numElements)}
	l.tail = prevElement
	var el *slotElement
	for i := numElements - 1; i > 0; i-- {
		el = &slotElement{id: uint32(i)}
		el.nextElement = prevElement
		prevElement = el
	}
	l.head = el

	for i := 1; i <= numElements; i++ {
		prevHead := l.head
		prevTail := l.tail

		el := l.popElement()

		require.Equal(t, prevHead, el)
		require.Equal(t, prevHead.nextElement, l.head)
		require.Equal(t, prevTail, l.tail)
		require.Equal(t, numElements-i, elementListLength(&l))
	}
}

func TestElementList_pushElement_ToEmptyList(t *testing.T) {
	l := ElementList{}
	el := &slotElement{}

	l.pushElement(el)
	require.Equal(t, 1, elementListLength(&l))
}

func TestElementList_pushElement_ToLenOneList(t *testing.T) {
	expectedElement := &slotElement{id: 1}
	l := ElementList{head: expectedElement, tail: expectedElement}
	el := &slotElement{}

	l.pushElement(el)
	require.Equal(t, el, expectedElement.nextElement)
	require.Equal(t, el, l.tail)
	require.Equal(t, 2, elementListLength(&l))
}

func TestElementList_pushElement_Multiple(t *testing.T) {
	firstElement := &slotElement{id: 1}
	l := ElementList{head: firstElement, tail: firstElement}
	numElements := 333

	for i := 2; i < numElements; i++ {
		prevHead := l.head
		prevTail := l.tail
		el := &slotElement{id: uint32(i)}

		l.pushElement(el)

		require.Equal(t, prevHead, l.head)
		require.Equal(t, el, prevTail.nextElement)
		require.Equal(t, el, l.tail)
		require.Equal(t, i, elementListLength(&l))
	}
}

func TestElementList_pushElement_popElement(t *testing.T) {
	l := ElementList{}
	el := &slotElement{id: 1}

	l.pushElement(el)
	res := l.popElement()
	require.Equal(t, el, res)
	require.Equal(t, 0, elementListLength(&l))
}

func TestInitElementsBuf(t *testing.T) {
	elements := initElementsBuf()
	require.Len(t, elements, slotSize)
	for i := 0; i < slotSize; i++ {
		require.Equal(t, EmptyElement, elements[i].activationStatus)
	}
}

func TestNewSlot(t *testing.T) {
	s := NewSlot(Future, testRealPulse)
	require.NotNil(t, s)
	require.Equal(t, Future, s.pulseState)
	require.Equal(t, testRealPulse, s.pulseNumber)
	require.Equal(t, Initializing, s.slotState)
	require.Empty(t, s.inputQueue.RemoveAll())
	require.Len(t, s.elements, slotSize)
	require.Len(t, s.elementListMap, 3)
	require.Equal(t, SlotStateMachine, s.slotStateMachine)
}

func TestSlot_getPulseNumber(t *testing.T) {
	s := testSlot(t, true, testRealPulse)

	pn := s.GetPulseNumber()
	require.Equal(t, testRealPulse, pn)
}

func TestSlot_getPulseData(t *testing.T) {
	s := testSlot(t, true, testRealPulse)

	pulse := s.GetPulseData()
	require.Equal(t, core.Pulse{PulseNumber: testRealPulse}, pulse)
}

func TestSlot_getNodeId(t *testing.T) {
	s := testSlot(t, true, testRealPulse)
	expectedNodeID := uint32(112233)
	s.nodeID = expectedNodeID

	nodeID := s.GetNodeID()
	require.Equal(t, expectedNodeID, nodeID)
}

func TestSlot_getNodeData(t *testing.T) {
	s := testSlot(t, true, testRealPulse)
	expectedNodeData := "some_test_node_data"
	s.nodeData = expectedNodeData

	nodeData := s.GetNodeData()
	require.Equal(t, expectedNodeData, nodeData)
}

func TestSlot_createElement(t *testing.T) {
	s := NewSlot(Future, testRealPulse)
	oldEmptyLen := elementListLength(s.elementListMap[EmptyElement])
	event := queue.OutputElement{}

	stateMachineMock := statemachine.NewStateMachineTypeMock(t)

	element, err := s.createElement(stateMachineMock, 1, event)
	require.NotNil(t, element)
	require.NoError(t, err)
	require.Equal(t, stateMachineMock, element.stateMachineType)
	require.Equal(t, uint16(1), element.state)
	require.Equal(t, uint32(0), element.id)
	require.Equal(t, ActiveElement, element.activationStatus)
	require.Equal(t, 1, elementListLength(s.elementListMap[ActiveElement]))
	require.Equal(t, oldEmptyLen-1, elementListLength(s.elementListMap[EmptyElement]))
}

func TestSlot_createElement_Err(t *testing.T) {
	s := NewSlot(Future, testRealPulse)
	oldEmptyLen := elementListLength(s.elementListMap[EmptyElement])
	delete(s.elementListMap, ActiveElement)
	event := queue.OutputElement{}

	stateMachineMock := statemachine.NewStateMachineTypeMock(t)

	element, err := s.createElement(stateMachineMock, 1, event)
	require.Nil(t, element)
	require.EqualError(t, err, "[ createElement ]: [ pushElement ] can't push element: list for status ActiveElement doesn't exist")
	require.Equal(t, oldEmptyLen, elementListLength(s.elementListMap[EmptyElement]))
}

func TestSlot_popElement(t *testing.T) {
	l := len3List()
	s := Slot{
		elementListMap: map[ActivationStatus]*ElementList{
			ActiveElement: &l,
		},
	}
	prevHead := s.elementListMap[ActiveElement].head
	prevTail := s.elementListMap[ActiveElement].tail

	element := s.popElement(ActiveElement)
	require.Equal(t, prevHead, element)
	require.Equal(t, prevHead.nextElement, s.elementListMap[ActiveElement].head)
	require.Equal(t, prevTail, s.elementListMap[ActiveElement].tail)
	require.Equal(t, 2, elementListLength(s.elementListMap[ActiveElement]))
}

func TestSlot_popElement_UnknownStatus(t *testing.T) {
	s := Slot{}
	unknownStatus := ActivationStatus(6767)

	element := s.popElement(unknownStatus)
	require.Nil(t, element)
}

func TestSlot_pushElement(t *testing.T) {
	l := len3List()
	s := Slot{
		elementListMap: map[ActivationStatus]*ElementList{
			ActiveElement: &l,
		},
	}
	prevHead := s.elementListMap[ActiveElement].head
	prevTail := s.elementListMap[ActiveElement].tail
	element := &slotElement{id: 777}
	prevID := element.id

	err := s.pushElement(ActiveElement, element)
	require.NoError(t, err)
	require.Equal(t, prevHead, s.elementListMap[ActiveElement].head)
	require.Equal(t, element, prevTail.nextElement)
	require.Equal(t, element, s.elementListMap[ActiveElement].tail)
	require.Equal(t, 4, elementListLength(s.elementListMap[ActiveElement]))

	require.Equal(t, prevID, element.id)
}

func TestSlot_pushElement_UnknownStatus(t *testing.T) {
	s := Slot{}
	unknownStatus := ActivationStatus(6767)
	element := &slotElement{id: 777}

	err := s.pushElement(unknownStatus, element)
	require.EqualError(t, err, "[ pushElement ] can't push element: list for status ActivationStatus(6767) doesn't exist")
}

func TestSlot_pushElement_Empty(t *testing.T) {
	l := len3List()
	s := Slot{
		elementListMap: map[ActivationStatus]*ElementList{
			EmptyElement: &l,
		},
	}
	prevHead := s.elementListMap[EmptyElement].head
	prevTail := s.elementListMap[EmptyElement].tail
	element := &slotElement{id: 777}
	prevID := element.id

	err := s.pushElement(EmptyElement, element)
	require.NoError(t, err)
	require.Equal(t, prevHead, s.elementListMap[EmptyElement].head)
	require.Equal(t, element, prevTail.nextElement)
	require.Equal(t, element, s.elementListMap[EmptyElement].tail)
	require.Equal(t, 4, elementListLength(s.elementListMap[EmptyElement]))

	require.Equal(t, prevID+slotElementDelta, element.id)
}

func TestNewSlotElement(t *testing.T) {
	s := newSlotElement(ActiveElement)
	require.NotNil(t, s)
	require.Equal(t, ActiveElement, s.activationStatus)
}
