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
	"errors"
	"testing"

	"github.com/insolar/insolar/conveyor/fsm"
	"github.com/insolar/insolar/conveyor/generator/matrix"
	"github.com/insolar/insolar/conveyor/handler"
	"github.com/insolar/insolar/conveyor/queue"
	"github.com/insolar/insolar/insolar"

	"github.com/stretchr/testify/require"
)

const testRealPulse = insolar.PulseNumber(1000)

func mockQueue(t *testing.T) *queue.QueueMock {
	qMock := queue.NewQueueMock(t)
	qMock.SinkPushFunc = func(p interface{}) (r error) {
		return nil
	}
	qMock.SinkPushAllFunc = func(p []interface{}) (r error) {
		return nil
	}
	qMock.PushSignalFunc = func(p uint32, p1 queue.SyncDone) (r error) {
		p1.SetResult(333)
		return nil
	}
	qMock.RemoveAllFunc = func() (r []queue.OutputElement) {
		return []queue.OutputElement{}
	}
	qMock.HasSignalFunc = func() (r bool) {
		return false
	}
	return qMock
}

func mockQueueReturnFalse(t *testing.T) *queue.QueueMock {
	qMock := queue.NewQueueMock(t)
	qMock.SinkPushFunc = func(p interface{}) (r error) {
		return errors.New("test error")
	}
	qMock.SinkPushAllFunc = func(p []interface{}) (r error) {
		return errors.New("test error")
	}
	qMock.PushSignalFunc = func(p uint32, p1 queue.SyncDone) (r error) {
		return errors.New("test error")
	}
	return qMock
}

type mockSyncDone struct {
	waiter    chan interface{}
	doneCount int
}

func (s *mockSyncDone) GetResult() []byte {
	result := <-s.waiter
	hash, ok := result.([]byte)
	if !ok {
		return []byte{}
	}
	return hash
}

func (s *mockSyncDone) SetResult(result interface{}) {
	s.waiter <- result
	s.doneCount += 1
}

func mockCallback() queue.SyncDone {
	return &mockSyncDone{
		waiter:    make(chan interface{}, 3),
		doneCount: 0,
	}
}

func testSlot(t *testing.T, isQueueOk bool, pulseNumber insolar.PulseNumber) *slot {
	var q *queue.QueueMock
	if isQueueOk {
		q = mockQueue(t)
	} else {
		q = mockQueueReturnFalse(t)
	}

	return &slot{
		inputQueue:  q,
		pulseNumber: pulseNumber,
		pulse:       insolar.Pulse{PulseNumber: pulseNumber},
	}
}

func len3List() elementList {
	el1 := &slotElement{id: 1}
	el2 := &slotElement{id: 2}
	el3 := &slotElement{id: 3}
	el1.nextElement = el2
	el2.nextElement = el3
	el3.prevElement = el2
	el2.prevElement = el1

	l := elementList{
		head:   el1,
		tail:   el3,
		length: 3,
	}
	return l
}

func TestElementList_removeElement_Nil(t *testing.T) {
	l := elementList{}
	l.removeElement(nil)
}

func TestElementList_removeElement_OnlyOne(t *testing.T) {
	expectedElement := &slotElement{id: 1}
	l := elementList{head: expectedElement, tail: expectedElement, length: 1}
	l.removeElement(expectedElement)
	require.Equal(t, elementList{}, l)
}

func TestElementList_removeElement_Head(t *testing.T) {
	l := len3List()
	headNext := l.head.nextElement
	tail := l.tail
	l.removeElement(l.head)
	require.Equal(t, headNext, l.head)
	require.Equal(t, tail, l.tail)
	require.Equal(t, 2, l.len())
}

func TestElementList_removeElement_Tail(t *testing.T) {
	l := len3List()
	tailPrev := l.tail.prevElement
	head := l.head
	l.removeElement(l.tail)
	require.Equal(t, tailPrev, l.tail)
	require.Equal(t, head, l.head)
	require.Equal(t, 2, l.len())
}

func TestElementList_removeElement_MiddleMultiple(t *testing.T) {
	l := elementList{}
	numElements := 333
	var el *slotElement
	var elements []*slotElement
	for i := 0; i < numElements; i++ {
		el = &slotElement{id: uint32(i)}
		l.pushElement(el)
		elements = append(elements, el)
	}

	for i := 1; i < numElements-1; i++ {
		prevHeadID := l.head.id
		prevTailID := l.tail.id

		l.removeElement(elements[i])

		require.Equal(t, prevHeadID, l.head.id)
		require.Equal(t, prevTailID, l.tail.id)
		require.Equal(t, numElements-i, l.len())
	}
}

func TestElementList_popElement_FromEmptyList(t *testing.T) {
	l := elementList{}
	el := l.popElement()
	require.Nil(t, el)
}

func TestElementList_isEmpty(t *testing.T) {
	list := elementList{}
	require.True(t, list.isEmpty())

	list.pushElement(newSlotElement(ActiveElement, nil))
	require.False(t, list.isEmpty())

	list.popElement()
	require.True(t, list.isEmpty())
}

func TestElementList_popElement_FromLenOneList(t *testing.T) {
	expectedElement := &slotElement{id: 1}
	l := elementList{head: expectedElement, tail: expectedElement, length: 1}

	el := l.popElement()
	require.Equal(t, expectedElement, el)
	require.Equal(t, 0, l.len())
}

func TestElementList_popElement_Multiple(t *testing.T) {
	l := elementList{}
	numElements := 333
	var el *slotElement
	for i := 0; i < numElements; i++ {
		el = &slotElement{id: uint32(i)}
		l.pushElement(el)
	}

	for i := 1; i < numElements; i++ {
		prevHead := *l.head
		prevTail := *l.tail

		el := l.popElement()

		require.Equal(t, prevHead.id, el.id)
		require.Equal(t, prevHead.nextElement, l.head)
		require.Equal(t, prevTail.id, l.tail.id)
		require.Equal(t, numElements-i, l.len())
	}
	// pop last element
	prevHead := *l.head
	el = l.popElement()
	require.Equal(t, prevHead.id, el.id)
	require.Equal(t, elementList{}, l)
}

func TestElementList_pushElement_ToEmptyList(t *testing.T) {
	l := elementList{}
	el := &slotElement{}

	l.pushElement(el)
	require.Equal(t, 1, l.len())
}

func TestElementList_pushElement_ToLenOneList(t *testing.T) {
	expectedElement := &slotElement{id: 1}
	l := elementList{head: expectedElement, tail: expectedElement, length: 1}
	el := &slotElement{}

	l.pushElement(el)
	require.Equal(t, el, expectedElement.nextElement)
	require.Equal(t, expectedElement, el.prevElement)
	require.Equal(t, el, l.tail)
	require.Equal(t, 2, l.len())
}

func TestElementList_pushElement_Multiple(t *testing.T) {
	firstElement := &slotElement{id: 1}
	l := elementList{head: firstElement, tail: firstElement, length: 1}
	numElements := 333

	for i := 2; i < numElements; i++ {
		prevHead := l.head
		prevTail := l.tail
		el := &slotElement{id: uint32(i)}

		l.pushElement(el)

		require.Equal(t, prevHead, l.head)
		require.Equal(t, el, prevTail.nextElement)
		require.Equal(t, el, l.tail)
		require.Equal(t, prevTail, el.prevElement)
		require.Equal(t, i, l.len())
	}
}

func TestElementList_pushElement_popElement(t *testing.T) {
	l := elementList{}
	el := &slotElement{id: 1}

	l.pushElement(el)
	res := l.popElement()
	require.Equal(t, el, res)
	require.Equal(t, 0, l.len())
}

func TestInitElementsBuf(t *testing.T) {
	elements, emptyList := initElementsBuf()
	require.Len(t, elements, slotSize)
	for i := 0; i < slotSize; i++ {
		require.Equal(t, EmptyElement, elements[i].activationStatus)
	}
	require.Equal(t, slotSize, emptyList.len())
}

func TestNewSlot(t *testing.T) {
	s := newSlot(Future, testRealPulse, nil)
	require.NotNil(t, s)
	require.Equal(t, Future, s.pulseState)
	require.Equal(t, testRealPulse, s.pulseNumber)
	require.Equal(t, Initializing, s.slotState)
	require.Empty(t, s.inputQueue.RemoveAll())
	require.Len(t, s.elements, slotSize)
	require.Len(t, s.elementListMap, 3)
	require.Equal(t, slotStateMachine, s.stateMachine)
}

func TestSlot_getPulseNumber(t *testing.T) {
	s := testSlot(t, true, testRealPulse)

	pn := s.GetPulseNumber()
	require.Equal(t, testRealPulse, pn)
}

func TestSlot_getPulseData(t *testing.T) {
	s := testSlot(t, true, testRealPulse)

	pulse := s.GetPulseData()
	require.Equal(t, insolar.Pulse{PulseNumber: testRealPulse}, pulse)
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
	s := newSlot(Future, testRealPulse, nil)
	oldEmptyLen := s.elementListMap[EmptyElement].len()
	event := queue.OutputElement{}

	stateMachineMock := makeMockStateMachine(t)

	element, err := s.createElement(stateMachineMock, 1, event)
	require.NotNil(t, element)
	require.NoError(t, err)
	require.Equal(t, stateMachineMock, element.stateMachine)
	require.Equal(t, fsm.StateID(1), element.state)
	require.Equal(t, uint32(0), element.id)
	require.Equal(t, ActiveElement, element.activationStatus)
	require.Equal(t, 1, s.elementListMap[ActiveElement].len())
	require.Equal(t, oldEmptyLen-1, s.elementListMap[EmptyElement].len())
}

func TestSlot_createElement_Err(t *testing.T) {
	s := newSlot(Future, testRealPulse, nil)
	oldEmptyLen := s.elementListMap[EmptyElement].len()
	delete(s.elementListMap, ActiveElement)
	event := queue.OutputElement{}

	stateMachineMock := makeMockStateMachine(t)

	element, err := s.createElement(stateMachineMock, 1, event)
	require.Nil(t, element)
	require.EqualError(t, err, "[ createElement ]: [ pushElement ] can't push element: list for status ActiveElement doesn't exist")
	require.Equal(t, oldEmptyLen, s.elementListMap[EmptyElement].len())
}

func TestSlot_hasElements_UnexistingState(t *testing.T) {
	s := newSlot(Present, 10, nil)
	badState := ActivationStatus(4444444)
	require.False(t, s.hasElements(badState))
}

func TestSlot_hasElements(t *testing.T) {
	s := newSlot(Present, 10, nil)
	require.False(t, s.hasElements(ActiveElement))
	require.False(t, s.hasElements(NotActiveElement))
	require.True(t, s.hasElements(EmptyElement))

	sm := makeMockStateMachine(t)

	_, err := s.createElement(sm, 20, queue.OutputElement{})
	require.NoError(t, err)

	require.True(t, s.hasElements(ActiveElement))
	require.False(t, s.hasElements(NotActiveElement))
	require.True(t, s.hasElements(EmptyElement))

	el := s.popElement(ActiveElement)
	require.False(t, s.hasElements(ActiveElement))
	require.False(t, s.hasElements(NotActiveElement))
	require.True(t, s.hasElements(EmptyElement))

	el.activationStatus = NotActiveElement
	err = s.pushElement(el)
	require.NoError(t, err)
	require.False(t, s.hasElements(ActiveElement))
	require.True(t, s.hasElements(NotActiveElement))
	require.True(t, s.hasElements(EmptyElement))

}

func TestSlot_popElement(t *testing.T) {
	l := len3List()
	s := slot{
		elementListMap: map[ActivationStatus]*elementList{
			ActiveElement: &l,
		},
	}
	prevHead := *s.elementListMap[ActiveElement].head
	prevTail := s.elementListMap[ActiveElement].tail

	element := s.popElement(ActiveElement)
	require.Equal(t, prevHead.id, element.id)
	require.Equal(t, prevHead.nextElement.id, s.elementListMap[ActiveElement].head.id)
	require.Nil(t, s.elementListMap[ActiveElement].head.prevElement)
	require.Equal(t, prevTail, s.elementListMap[ActiveElement].tail)
	require.Equal(t, 2, s.elementListMap[ActiveElement].len())
}

func TestSlot_popElement_UnknownStatus(t *testing.T) {
	s := slot{}
	unknownStatus := ActivationStatus(6767)

	element := s.popElement(unknownStatus)
	require.Nil(t, element)
}

func TestSlot_pushElement(t *testing.T) {
	l := len3List()
	s := slot{
		elementListMap: map[ActivationStatus]*elementList{
			ActiveElement: &l,
		},
	}
	prevHead := s.elementListMap[ActiveElement].head
	prevTail := s.elementListMap[ActiveElement].tail
	element := &slotElement{id: 777}
	prevID := element.id
	element.activationStatus = ActiveElement

	err := s.pushElement(element)
	require.NoError(t, err)
	require.Equal(t, prevHead, s.elementListMap[ActiveElement].head)
	require.Equal(t, element, prevTail.nextElement)
	require.Equal(t, element, s.elementListMap[ActiveElement].tail)
	require.Equal(t, prevTail, s.elementListMap[ActiveElement].tail.prevElement)
	require.Equal(t, 4, s.elementListMap[ActiveElement].len())

	require.Equal(t, prevID, element.id)
}

func TestSlot_pushElement_UnknownStatus(t *testing.T) {
	s := slot{}
	unknownStatus := ActivationStatus(6767)
	element := &slotElement{id: 777}
	element.activationStatus = unknownStatus

	err := s.pushElement(element)
	require.EqualError(t, err, "[ pushElement ] can't push element: list for status ActivationStatus(6767) doesn't exist")
}

func TestSlot_pushElement_Empty(t *testing.T) {
	l := len3List()
	s := slot{
		elementListMap: map[ActivationStatus]*elementList{
			EmptyElement: &l,
		},
	}
	prevHead := s.elementListMap[EmptyElement].head
	prevTail := s.elementListMap[EmptyElement].tail
	element := &slotElement{id: 777}
	element.activationStatus = EmptyElement
	prevID := element.id

	err := s.pushElement(element)
	require.NoError(t, err)
	require.Equal(t, prevHead, s.elementListMap[EmptyElement].head)
	require.Equal(t, element, prevTail.nextElement)
	require.Equal(t, element, s.elementListMap[EmptyElement].tail)
	require.Equal(t, prevTail, s.elementListMap[EmptyElement].tail.prevElement)
	require.Equal(t, 4, s.elementListMap[EmptyElement].len())

	require.Equal(t, prevID+slotElementDelta, element.id)
}

func TestSlot_extractSlotElementByID(t *testing.T) {
	sm := makeMockStateMachine(t)
	slot := newSlot(Present, 10, nil)

	var elements []*slotElement

	for i := 1; i < 100; i++ {
		el, err := slot.createElement(sm, fsm.StateID(20+i), queue.OutputElement{})
		require.NoError(t, err)
		elements = append(elements, el)
	}

	require.NotEqual(t, 0, len(elements))

	listLen := slot.len(ActiveElement)
	for i := 1; i < len(elements); i++ {
		require.Equal(t, elements[i].state, slot.extractSlotElementByID(elements[i].id).state)
		require.Equal(t, listLen-i, slot.len(ActiveElement))
	}

	for i := 1; i < 100; i++ {
		_, err := slot.createElement(sm, fsm.StateID(2000+i), queue.OutputElement{})
		require.NoError(t, err)

		el := slot.popElement(ActiveElement)
		slot.pushElement(el)
	}

	listLen = slot.len(ActiveElement)
	for i := 1; i < len(elements); i++ {
		require.Equal(t, elements[i].state, slot.extractSlotElementByID(elements[i].id).state)
		require.Equal(t, listLen-i, slot.len(ActiveElement))
	}
}

func makeMockStateMachine(t *testing.T) matrix.StateMachine {
	sm := matrix.NewStateMachineMock(t)

	sm.GetTransitionHandlerFunc = func(p fsm.StateID) (r handler.TransitHandler) {
		return func(element fsm.SlotElementHelper) (interface{}, fsm.ElementState, error) {
			return nil, 0, nil
		}
	}

	sm.GetMigrationHandlerFunc = func(p fsm.StateID) (r handler.MigrationHandler) {
		return func(element fsm.SlotElementHelper) (interface{}, fsm.ElementState, error) {
			return nil, 0, nil
		}
	}

	return sm
}

func TestSlot_PopPushMultiple(t *testing.T) {
	sm := makeMockStateMachine(t)
	slot := newSlot(Present, 10, nil)

	slot.createElement(sm, 33, queue.OutputElement{})

	el := slot.popElement(ActiveElement)
	require.NotNil(t, el)
	slot.pushElement(el)

	el = slot.popElement(ActiveElement)
	require.NotNil(t, el)
	slot.pushElement(el)

	el = slot.popElement(ActiveElement)
	require.NotNil(t, el)
}

func TestSlot_pushElementToEmpty_ExtractByID(t *testing.T) {
	s := newSlot(Future, testRealPulse, nil)

	element := s.popElement(EmptyElement)
	oldID := element.id
	err := s.pushElement(element)
	require.NoError(t, err)

	elementByID := s.extractSlotElementByID(oldID)
	require.Nil(t, elementByID)

	elementByID = s.extractSlotElementByID(element.id)
	require.NotNil(t, elementByID)
}

func TestSlot_extractSlotElementByID_NotExist(t *testing.T) {
	s := newSlot(Present, 10, nil)

	elementByID := s.extractSlotElementByID(slotElementDelta)
	require.Nil(t, elementByID)
}

func TestSlot_extractSlotElementByID_pushElement(t *testing.T) {
	s := newSlot(Present, 10, nil)

	elementByID := s.extractSlotElementByID(0)
	require.NotNil(t, elementByID)

	elementByID.activationStatus = ActiveElement

	err := s.pushElement(elementByID)
	require.NoError(t, err)

	require.Equal(t, 1, s.elementListMap[ActiveElement].len())
	require.Equal(t, slotSize-1, s.elementListMap[EmptyElement].len())
}

func TestNewSlotElement(t *testing.T) {
	s := newSlotElement(ActiveElement, nil)
	require.NotNil(t, s)
	require.Equal(t, ActiveElement, s.activationStatus)
}
