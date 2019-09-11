///
//    Copyright 2019 Insolar Technologies
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
///

package smachine

type ContextMarker = uint64

type StateUpdate struct {
	marker  ContextMarker
	link    *Slot
	param0  uint32
	param1  interface{}
	step    SlotStep
	updType uint16
}

func (u StateUpdate) IsZero() bool {
	return u.marker == 0 && u.updType == 0
}

func (u StateUpdate) getLink() SlotLink {
	if u.link == nil {
		return NoLink()
	}
	return SlotLink{SlotID(u.param0), u.link}
}

func (u StateUpdate) ensureMarker(marker ContextMarker) StateUpdate {
	if u.marker != marker {
		panic("illegal state")
	}
	return u
}

func newStateUpdate(marker ContextMarker, updType uint16, slotStep SlotStep, param interface{}) StateUpdate {
	return StateUpdate{
		marker:  marker,
		param1:  param,
		step:    slotStep,
		updType: updType,
	}
}

func newStateUpdateUint(marker ContextMarker, updType uint16, slotStep SlotStep, param uint32) StateUpdate {
	return StateUpdate{
		marker:  marker,
		param0:  param,
		step:    slotStep,
		updType: updType,
	}
}

func newStateUpdateUintAndParam(marker ContextMarker, updType uint16, slotStep SlotStep, param0 uint32, param1 interface{}) StateUpdate {
	return StateUpdate{
		marker:  marker,
		param0:  param0,
		param1:  param1,
		step:    slotStep,
		updType: updType,
	}
}

func newStateUpdateLink(marker ContextMarker, updType uint16, link SlotLink, slotStep SlotStep, param interface{}) StateUpdate {
	return StateUpdate{
		marker:  marker,
		param1:  param,
		link:    link.s,
		param0:  uint32(link.id),
		step:    slotStep,
		updType: updType,
	}
}
