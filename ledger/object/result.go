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

package object

import (
	"io"

	"github.com/insolar/insolar/insolar"
)

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

// ResultRecord represents result of a VM method.
type ResultRecord struct {
	Object  insolar.ID
	Request insolar.Reference
	Payload []byte
}

// WriteHashData writes record data to provided writer. This data is used to calculate record's hash.
func (r *ResultRecord) WriteHashData(w io.Writer) (int, error) {
	return w.Write(EncodeVirtual(r))
}

// SideEffectRecord is a record which is created in response to a request.
type SideEffectRecord struct {
	Domain  insolar.Reference
	Request insolar.Reference
}

// TypeRecord is a code interface declaration.
type TypeRecord struct {
	SideEffectRecord

	TypeDeclaration []byte
}

// WriteHashData writes record data to provided writer. This data is used to calculate record's hash.
func (r *TypeRecord) WriteHashData(w io.Writer) (int, error) {
	return w.Write(EncodeVirtual(r))
}

// CodeRecord is a code storage record.
type CodeRecord struct {
	SideEffectRecord

	Code        *insolar.ID
	MachineType insolar.MachineType
}

// WriteHashData writes record data to provided writer. This data is used to calculate record's hash.
func (r *CodeRecord) WriteHashData(w io.Writer) (int, error) {
	return w.Write(EncodeVirtual(r))
}

// StateRecord is a record containing data for an object state.
type StateRecord struct {
	Memory      *insolar.ID
	Image       insolar.Reference // If code or prototype object reference.
	IsPrototype bool              // If true, Image should point to a prototype object. Otherwise to a code.
}

// GetMemory returns state indexStorage.
func (r *StateRecord) GetMemory() *insolar.ID {
	return r.Memory
}

// GetImage returns state code.
func (r *StateRecord) GetImage() *insolar.Reference {
	return &r.Image
}

// GetIsPrototype returns state code.
func (r *StateRecord) GetIsPrototype() bool {
	return r.IsPrototype
}

// ActivateRecord is produced when we instantiate new object from an available prototype.
type ActivateRecord struct {
	SideEffectRecord
	StateRecord

	Parent     insolar.Reference
	IsDelegate bool
}

// PrevStateID returns previous state id.
func (r *ActivateRecord) PrevStateID() *insolar.ID {
	return nil
}

// StateID returns state id.
func (r *ActivateRecord) ID() StateID {
	return StateActivation
}

// WriteHashData writes record data to provided writer. This data is used to calculate record's hash.
func (r *ActivateRecord) WriteHashData(w io.Writer) (int, error) {
	return w.Write(EncodeVirtual(r))
}

// AmendRecord is an amendment record for objects.
type AmendRecord struct {
	SideEffectRecord
	StateRecord

	PrevState insolar.ID
}

// PrevStateID returns previous state id.
func (r *AmendRecord) PrevStateID() *insolar.ID {
	return &r.PrevState
}

// StateID returns state id.
func (r *AmendRecord) ID() StateID {
	return StateAmend
}

// WriteHashData writes record data to provided writer. This data is used to calculate record's hash.
func (r *AmendRecord) WriteHashData(w io.Writer) (int, error) {
	return w.Write(EncodeVirtual(r))
}

// DeactivationRecord marks targeted object as disabled.
type DeactivationRecord struct {
	SideEffectRecord
	PrevState insolar.ID
}

// PrevStateID returns previous state id.
func (r *DeactivationRecord) PrevStateID() *insolar.ID {
	return &r.PrevState
}

// StateID returns state id.
func (r *DeactivationRecord) ID() StateID {
	return StateDeactivation
}

// WriteHashData writes record data to provided writer. This data is used to calculate record's hash.
func (r *DeactivationRecord) WriteHashData(w io.Writer) (int, error) {
	return w.Write(EncodeVirtual(r))
}

// GetMachineType returns state code machine type.
func (*DeactivationRecord) GetMachineType() insolar.MachineType {
	return insolar.MachineTypeNotExist
}

// GetMemory returns state indexStorage.
func (*DeactivationRecord) GetMemory() *insolar.ID {
	return nil
}

// GetImage returns state code.
func (r *DeactivationRecord) GetImage() *insolar.Reference {
	return nil
}

// GetIsPrototype returns state code.
func (r *DeactivationRecord) GetIsPrototype() bool {
	return false
}
