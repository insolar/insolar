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

type reactivateMode interface{}

//go:generate minimock -i github.com/insolar/insolar/conveyor/interfaces/slot.SlotElementHelper -o ./ -s _mock.go
type SlotElementHelper interface {
	SlotElementRestrictedHelper
	InformParent(payload interface{}) bool
	DeactivateTill(reactivateOn reactivateMode)
	SendTask(adapterID uint32, taskPayload interface{}, respHandlerID uint32) error
	// joinSequence( sequenceKey map-key,sequenceOrder uint64 )
	// isSequenceHead() bool
}

//go:generate minimock -i github.com/insolar/insolar/conveyor/interfaces/slot.SlotElementRestrictedHelper -o ./ -s _mock.go
type SlotElementRestrictedHelper interface {
	SlotElementReadOnly

	GetParentElementID() uint32
	GetInputEvent() interface{}
	GetPayload() interface{}

	Reactivate()
	LeaveSequence()
}

//go:generate minimock -i github.com/insolar/insolar/conveyor/interfaces/slot.SlotElementReadOnly -o ./ -s _mock.go
type SlotElementReadOnly interface {
	GetElementID() uint32
	GetNodeID() uint32
	GetType() int
	GetState() int
}
