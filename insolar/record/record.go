//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package record

import (
	"github.com/insolar/insolar/insolar"
)

type Record interface{}

// type VirtualRecord interface {
// 	// WriteHashData writes record data to provided writer. This data is used to calculate record's hash.
// 	WriteHashData(w io.Writer) (int, error)
// }

// type MaterialRecord struct {
// 	Record VirtualRecord
//
// 	JetID insolar.JetID
// }

// StateID is a state of lifeline records.
type StateID int

const (
	// StateUndefined is used for special cases.
	StateUndefined = StateID(iota)
	// StateActivation means it's an activation record.
	StateActivation
	// StateAmend means it's an amend record.
	StateAmend
	// StateDeactivation means it's a deactivation record.
	StateDeactivation
)

// State is common object state record.
type State interface {
	// StateID returns state id.
	// TODO: rename to StateID()
	ID() StateID
	// GetImage returns state code.
	GetImage() *insolar.Reference
	// GetIsPrototype returns state code.
	GetIsPrototype() bool
	// GetMemory returns state indexStorage.
	GetMemory() *insolar.ID
	// PrevStateID returns previous state id.
	PrevStateID() *insolar.ID
}

// TODO it's a hack for object.State interface compatibility
func (Activate) ID() StateID {
	return StateActivation
}

func (p Activate) GetImage() *insolar.Reference {
	return &p.Image
}

func (p Activate) GetIsPrototype() bool {
	return p.IsPrototype
}

func (p Activate) GetMemory() *insolar.ID {
	return &p.Memory
}

func (Activate) PrevStateID() *insolar.ID {
	return nil
}

// TODO it's a hack for object.State interface compatibility
func (Amend) ID() StateID {
	return StateAmend
}

func (p Amend) GetImage() *insolar.Reference {
	return &p.Image
}

func (p Amend) GetIsPrototype() bool {
	return p.IsPrototype
}

func (p Amend) GetMemory() *insolar.ID {
	return &p.Memory
}

func (p Amend) PrevStateID() *insolar.ID {
	return &p.PrevState
}

// TODO it's a hack for object.State interface compatibility
func (Deactivate) ID() StateID {
	return StateDeactivation
}

func (p Deactivate) GetImage() *insolar.Reference {
	return nil
}

func (p Deactivate) GetIsPrototype() bool {
	return false
}

func (p Deactivate) GetMemory() *insolar.ID {
	return nil
}

func (p Deactivate) PrevStateID() *insolar.ID {
	return &p.PrevState
}

// TODO it's a hack for object.State interface compatibility
func (Genesis) PrevStateID() *insolar.ID {
	return nil
}

func (Genesis) ID() StateID {
	return StateActivation
}

func (Genesis) GetMemory() *insolar.ID {
	return nil
}

func (Genesis) GetImage() *insolar.Reference {
	return nil
}

func (Genesis) GetIsPrototype() bool {
	return false
}
