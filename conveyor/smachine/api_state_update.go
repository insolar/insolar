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

type ContextMarker *struct{}

type StateUpdate struct {
	marker  ContextMarker
	link    SlotLink
	param   interface{}
	step    SlotStep
	updType uint16
}

func (u StateUpdate) IsZero() bool {
	return u.marker == nil && u.updType == 0
}

func NewStateUpdate(marker ContextMarker, updType uint16, slotStep SlotStep, param interface{}) StateUpdate {
	return StateUpdate{
		marker:  marker,
		param:   param,
		step:    slotStep,
		updType: updType,
	}
}

func NewStateUpdateLink(marker ContextMarker, updType uint16, link SlotLink, slotStep SlotStep, param interface{}) StateUpdate {
	return StateUpdate{
		marker:  marker,
		param:   param,
		link:    link,
		step:    slotStep,
		updType: updType,
	}
}

func EnsureUpdateContext(p ContextMarker, u StateUpdate) StateUpdate {
	if u.marker != p {
		panic("illegal value")
	}
	return u
}

func ExtractStateUpdate(u StateUpdate) (updType uint16, slotStep SlotStep, param interface{}) {
	return u.updType, u.step, u.param
}

func ExtractStateUpdateParam(u StateUpdate) (updType uint16, param interface{}) {
	return u.updType, u.param
}
