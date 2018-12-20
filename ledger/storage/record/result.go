/*
 *    Copyright 2018 Insolar
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

package record

import (
	"io"

	"github.com/insolar/insolar/core"
)

// State is a state of lifeline records.
type State int

const (
	// StateUndefined is used for special cases.
	StateUndefined = State(iota)
	// StateActivation means it's an activation record.
	StateActivation
	// StateAmend means it's an amend record.
	StateAmend
	// StateDeactivation means it's a deactivation record.
	StateDeactivation
)

// ObjectState is common object state record.
type ObjectState interface {
	// State returns state id.
	State() State
	// GetImage returns state code.
	GetImage() *core.RecordRef
	// GetIsPrototype returns state code.
	GetIsPrototype() bool
	// GetMemory returns state memory.
	GetMemory() *core.RecordID
	// PrevStateID returns previous state id.
	PrevStateID() *core.RecordID
}

// ResultRecord represents result of a VM method.
type ResultRecord struct {
	Object  core.RecordID
	Request core.RecordRef
	Payload []byte
}

// Type implementation of Record interface.
func (ResultRecord) Type() TypeID {
	return typeResult
}

// WriteHashData writes record data to provided writer. This data is used to calculate record's hash.
func (r *ResultRecord) WriteHashData(w io.Writer) (int, error) {
	return w.Write(SerializeRecord(r))
}

// SideEffectRecord is a record which is created in response to a request.
type SideEffectRecord struct {
	Domain  core.RecordRef
	Request core.RecordRef
}

// TypeRecord is a code interface declaration.
type TypeRecord struct {
	SideEffectRecord

	TypeDeclaration []byte
}

// Type implementation of Record interface.
func (r *TypeRecord) Type() TypeID { return typeType }

// WriteHashData writes record data to provided writer. This data is used to calculate record's hash.
func (r *TypeRecord) WriteHashData(w io.Writer) (int, error) {
	return w.Write(SerializeRecord(r))
}

// CodeRecord is a code storage record.
type CodeRecord struct {
	SideEffectRecord

	Code        *core.RecordID
	MachineType core.MachineType
}

// Type implementation of Record interface.
func (r *CodeRecord) Type() TypeID { return typeCode }

// WriteHashData writes record data to provided writer. This data is used to calculate record's hash.
func (r *CodeRecord) WriteHashData(w io.Writer) (int, error) {
	return w.Write(SerializeRecord(r))
}

// ObjectStateRecord is a record containing data for an object state.
type ObjectStateRecord struct {
	Memory      *core.RecordID
	Image       core.RecordRef // If code or prototype object reference.
	IsPrototype bool           // If true, Image should point to a prototype object. Otherwise to a code.
}

// GetMemory returns state memory.
func (r *ObjectStateRecord) GetMemory() *core.RecordID {
	return r.Memory
}

// GetImage returns state code.
func (r *ObjectStateRecord) GetImage() *core.RecordRef {
	return &r.Image
}

// GetIsPrototype returns state code.
func (r *ObjectStateRecord) GetIsPrototype() bool {
	return r.IsPrototype
}

// ObjectActivateRecord is produced when we instantiate new object from an available prototype.
type ObjectActivateRecord struct {
	SideEffectRecord
	ObjectStateRecord

	Parent     core.RecordRef
	IsDelegate bool
}

// PrevStateID returns previous state id.
func (r *ObjectActivateRecord) PrevStateID() *core.RecordID {
	return nil
}

// State returns state id.
func (r *ObjectActivateRecord) State() State {
	return StateActivation
}

// Type implementation of Record interface.
func (r *ObjectActivateRecord) Type() TypeID { return typeActivate }

// WriteHashData writes record data to provided writer. This data is used to calculate record's hash.
func (r *ObjectActivateRecord) WriteHashData(w io.Writer) (int, error) {
	return w.Write(SerializeRecord(r))
}

// ObjectAmendRecord is an amendment record for objects.
type ObjectAmendRecord struct {
	SideEffectRecord
	ObjectStateRecord

	PrevState core.RecordID
}

// PrevStateID returns previous state id.
func (r *ObjectAmendRecord) PrevStateID() *core.RecordID {
	return &r.PrevState
}

// State returns state id.
func (r *ObjectAmendRecord) State() State {
	return StateAmend
}

// Type implementation of Record interface.
func (r *ObjectAmendRecord) Type() TypeID { return typeAmend }

// WriteHashData writes record data to provided writer. This data is used to calculate record's hash.
func (r *ObjectAmendRecord) WriteHashData(w io.Writer) (int, error) {
	return w.Write(SerializeRecord(r))
}

// DeactivationRecord marks targeted object as disabled.
type DeactivationRecord struct {
	SideEffectRecord
	PrevState core.RecordID
}

// PrevStateID returns previous state id.
func (r *DeactivationRecord) PrevStateID() *core.RecordID {
	return &r.PrevState
}

// State returns state id.
func (r *DeactivationRecord) State() State {
	return StateDeactivation
}

// Type implementation of Record interface.
func (r *DeactivationRecord) Type() TypeID { return typeDeactivate }

// WriteHashData writes record data to provided writer. This data is used to calculate record's hash.
func (r *DeactivationRecord) WriteHashData(w io.Writer) (int, error) {
	return w.Write(SerializeRecord(r))
}

// GetMachineType returns state code machine type.
func (*DeactivationRecord) GetMachineType() core.MachineType {
	return core.MachineTypeNotExist
}

// GetMemory returns state memory.
func (*DeactivationRecord) GetMemory() *core.RecordID {
	return nil
}

// GetImage returns state code.
func (r *DeactivationRecord) GetImage() *core.RecordRef {
	return nil
}

// GetIsPrototype returns state code.
func (r *DeactivationRecord) GetIsPrototype() bool {
	return false
}
