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

package conveyer

//
// import (
// 	"github.com/insolar/insolar/core"
// )
//
// type PulseState int
//
// //go:generate stringer -type=PulseState
// const (
// 	Unallocated = PulseState(iota)
// 	Future
// 	Present
// 	Past
// 	Antique
// )
//
// type SlotState int
//
// //go:generate stringer -type=SlotState
// const (
// 	Initializing = SlotState(iota)
// 	Working
// 	Suspending
// )
//
// const slotSize = 10000
//
// type slotDetails interface {
// 	getPulseNumber() core.PulseNumber
// 	// getNodeId() uint32
// 	getPulseData() *core.Pulse
// 	// getNodeData() interface{}
// }
//
// type SlotConfiguration struct{}
//
// type Slot struct {
// 	configuration SlotConfiguration
// 	inputQueue    NonBlockingQueue
// 	pulseState    PulseState
// 	slotState     SlotState
// 	pulse         *core.Pulse
// 	pulseNumber   core.PulseNumber
// 	elements      []SlotElement
// 	activeElement *SlotElement
// 	emptyElement  *SlotElement
// }
//
// func NewSlot(pulseState PulseState, pulseNumber core.PulseNumber) *Slot {
// 	slotState := Initializing
// 	if pulseState == Antique {
// 		slotState = Working
// 	}
//
// 	elements := make([]SlotElement, slotSize)
// 	emptySlotElement := NewSlotElement()
// 	for i := 0; i < slotSize; i++ {
// 		elements[i] = *emptySlotElement
// 	}
//
// 	return &Slot{
// 		pulseState: pulseState,
// 		// TODO: use newQueue or smth
// 		inputQueue:    &Queue{},
// 		pulseNumber:   pulseNumber,
// 		slotState:     slotState,
// 		elements:      elements,
// 		activeElement: nil,
// 		emptyElement:  &elements[0],
// 	}
// }
//
// type SlotElement struct {
// 	id       uint32
// 	payload  interface{}
// 	state    uint16
// 	listNext *SlotElement
// 	isActive bool
// }
//
// func NewSlotElement() *SlotElement {
// 	return &SlotElement{}
// }
//
// func (s *Slot) getPulseNumber() core.PulseNumber {
// 	return s.pulseNumber
// }
//
// func (s *Slot) getPulseData() *core.Pulse {
// 	return s.pulse
// }
//
// func (s *Slot) createElement() *SlotElement {
// 	return &SlotElement{}
// }
//
// func (s *Slot) getActiveElement() *SlotElement {
// 	return s.activeElement
// }
