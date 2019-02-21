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
	el1.listNext = el2
	el2.listNext = el3

	l := ElementList{
		head: el1,
		tail: el3,
	}
	return l
}

func TestElementList_len_Empty(t *testing.T) {
	l := ElementList{}
	listLen := l.len()
	require.Equal(t, 0, listLen)
}

func TestElementList_len(t *testing.T) {
	l := len3List()
	listLen := l.len()
	require.Equal(t, 3, listLen)
}

func TestElementList_popElement_Empty(t *testing.T) {
	l := ElementList{}
	el := l.popElement()
	require.Nil(t, el)
}

func TestElementList_popElement_OnlyOne(t *testing.T) {
	expectedElement := &slotElement{id: 1}
	l := ElementList{head: expectedElement, tail: expectedElement}

	el := l.popElement()
	require.Equal(t, expectedElement, el)
	require.Equal(t, 0, l.len())
}

func TestElementList_popElement(t *testing.T) {
	l := len3List()
	prevHead := l.head
	prevTail := l.tail

	el := l.popElement()
	require.Equal(t, prevHead, el)
	require.Equal(t, prevHead.listNext, l.head)
	require.Equal(t, prevTail, l.tail)
	require.Equal(t, 2, l.len())
}

func TestElementList_pushElement_Empty(t *testing.T) {
	l := ElementList{}
	el := &slotElement{}

	l.pushElement(el)
	require.Equal(t, 1, l.len())
}

func TestElementList_pushElement_OnlyOne(t *testing.T) {
	expectedElement := &slotElement{id: 1}
	l := ElementList{head: expectedElement, tail: expectedElement}
	el := &slotElement{}

	l.pushElement(el)
	require.Equal(t, el, expectedElement.listNext)
	require.Equal(t, el, l.tail)
	require.Equal(t, 2, l.len())
}

func TestElementList_pushElement(t *testing.T) {
	l := len3List()
	prevHead := l.head
	prevTail := l.tail
	el := &slotElement{}

	l.pushElement(el)
	require.Equal(t, prevHead, l.head)
	require.Equal(t, el, prevTail.listNext)
	require.Equal(t, el, l.tail)
	require.Equal(t, 4, l.len())
}

func TestInitElementsBuf(t *testing.T) {
	elements := initElementsBuf()
	require.Len(t, elements, slotSize)
	require.Equal(t, SlotStateMachine, elements[0])
	for i := 1; i < slotSize; i++ {
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
}

func TestSlot_getPulseNumber(t *testing.T) {
	s := testSlot(t, true, testRealPulse)

	pn := s.getPulseNumber()
	require.Equal(t, testRealPulse, pn)
}

func TestSlot_getPulseData(t *testing.T) {
	s := testSlot(t, true, testRealPulse)

	pulse := s.getPulseData()
	require.Equal(t, core.Pulse{PulseNumber: testRealPulse}, pulse)
}

func TestSlot_getNodeId(t *testing.T) {
	s := testSlot(t, true, testRealPulse)
	expectedNodeID := uint32(112233)
	s.nodeID = expectedNodeID

	nodeID := s.getNodeID()
	require.Equal(t, expectedNodeID, nodeID)
}

func TestSlot_getNodeData(t *testing.T) {
	s := testSlot(t, true, testRealPulse)
	expectedNodeData := "some_test_node_data"
	s.nodeData = expectedNodeData

	nodeData := s.getNodeData()
	require.Equal(t, expectedNodeData, nodeData)
}

func TestSlot_createElement(t *testing.T) {
	s := NewSlot(Future, testRealPulse)
	event := queue.OutputElement{}

	element := s.createElement("testStateMachineType", 1, event)
	require.NotNil(t, element)
	require.Equal(t, "testStateMachineType", element.stateMachineType)
	require.Equal(t, uint16(1), element.state)
	require.Equal(t, uint32(1+slotElementDelta), element.id)
	require.Equal(t, ActiveElement, element.activationStatus)
	require.Equal(t, 1, s.elementListMap[ActiveElement].len())
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
	require.Equal(t, prevHead.listNext, s.elementListMap[ActiveElement].head)
	require.Equal(t, prevTail, s.elementListMap[ActiveElement].tail)
	require.Equal(t, 2, s.elementListMap[ActiveElement].len())
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
	element := &slotElement{}

	s.pushElement(ActiveElement, element)
	require.Equal(t, prevHead, s.elementListMap[ActiveElement].head)
	require.Equal(t, element, prevTail.listNext)
	require.Equal(t, element, s.elementListMap[ActiveElement].tail)
	require.Equal(t, 4, s.elementListMap[ActiveElement].len())
}

func TestNewSlotElement(t *testing.T) {
	s := NewSlotElement(ActiveElement)
	require.NotNil(t, s)
	require.Equal(t, ActiveElement, s.activationStatus)
}
