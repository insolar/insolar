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

// ReactivateMode represents reason of reactivating of slot element
type ReactivateMode int

//go:generate stringer -type=ReactivateMode
const (
	Empty = ReactivateMode(iota)
	Response
	Tick
	SeqHead
)

// SlotElementHelper gives access to slot element
//go:generate minimock -i github.com/insolar/insolar/conveyor/interfaces/slot.SlotElementHelper -o ./ -s _mock.go
type SlotElementHelper interface {
	SlotElementRestrictedHelper
	InformParent(payload interface{}) bool
	DeactivateTill(reactivateOn ReactivateMode)
	SendTask(adapterID uint32, taskPayload interface{}, respHandlerID uint32) error
	// JoinSequence( sequenceKey map-key,sequenceOrder uint64 )
	// IsSequenceHead() bool
}

// SlotElementRestrictedHelper is restricted part of SlotElementHelper
//go:generate minimock -i github.com/insolar/insolar/conveyor/interfaces/slot.SlotElementRestrictedHelper -o ./ -s _mock.go
type SlotElementRestrictedHelper interface {
	SlotElementReadOnly

	GetParentElementID() uint32
	GetInputEvent() interface{}
	GetPayload() interface{}

	Reactivate()
	LeaveSequence()
}

// SlotElementReadOnly gives read-only access to slot element
//go:generate minimock -i github.com/insolar/insolar/conveyor/interfaces/slot.SlotElementReadOnly -o ./ -s _mock.go
type SlotElementReadOnly interface {
	GetElementID() uint32
	GetNodeID() uint32
	GetType() int
	GetState() uint16
}
