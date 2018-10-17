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
	"github.com/insolar/insolar/core"
)

// ClassState is common class state record.
type ClassState interface {
	// IsDeactivation determines if current state is deactivation.
	IsDeactivation() bool
	// GetCode returns state code.
	GetCode() *Reference
	// GetMachineType returns state code machine type.
	GetMachineType() core.MachineType
}

// ObjectState is common object state record.
type ObjectState interface {
	// IsDeactivation determines if current state is deactivation.
	IsDeactivation() bool
	// GetMemory returns state memory.
	GetMemory() []byte
}

// ResultRecord is a record which is created in response to a request.
type ResultRecord struct {
	Domain  Reference
	Request Reference
}

// CodeRecord is a code storage record.
type CodeRecord struct {
	ResultRecord

	Code        []byte
	MachineType core.MachineType
}

// TypeRecord is a code interface declaration.
type TypeRecord struct {
	ResultRecord

	TypeDeclaration []byte
}

// ClassStateRecord is a record containing data for a class state.
type ClassStateRecord struct {
	Code        Reference
	MachineType core.MachineType
}

// GetMachineType returns state code machine type.
func (r *ClassStateRecord) GetMachineType() core.MachineType {
	return r.MachineType
}

// GetCode returns state code.
func (r *ClassStateRecord) GetCode() *Reference {
	return &r.Code
}

// IsDeactivation determines if current state is deactivation.
func (r *ClassStateRecord) IsDeactivation() bool {
	return false
}

// ClassActivateRecord is produced when we "activate" new contract class.
type ClassActivateRecord struct {
	ResultRecord
	ClassStateRecord
}

// ClassAmendRecord is an amendment record for classes.
type ClassAmendRecord struct {
	ResultRecord
	ClassStateRecord

	PrevState  ID
	Migrations []Reference
}

// ObjectStateRecord is a record containing data for an object state.
type ObjectStateRecord struct {
	Memory Memory
}

// IsDeactivation determines if current state is deactivation.
func (r *ObjectStateRecord) IsDeactivation() bool {
	return false
}

// GetMemory returns state memory.
func (r *ObjectStateRecord) GetMemory() []byte {
	return r.Memory
}

// ObjectActivateRecord is produced when we instantiate new object from an available class.
type ObjectActivateRecord struct {
	ResultRecord
	ObjectStateRecord

	Class    Reference
	Parent   Reference
	Delegate bool
}

// ObjectAmendRecord is an amendment record for objects.
type ObjectAmendRecord struct {
	ResultRecord
	ObjectStateRecord

	PrevState ID
}

// DeactivationRecord marks targeted object as disabled.
type DeactivationRecord struct {
	ResultRecord
	PrevState ID
}

// GetMachineType returns state code machine type.
func (*DeactivationRecord) GetMachineType() core.MachineType {
	return core.MachineTypeNotExist
}

// IsDeactivation determines if current state is deactivation.
func (*DeactivationRecord) IsDeactivation() bool {
	return true
}

// IsAmend determines if current state is amend.
func (*DeactivationRecord) IsAmend() bool {
	return false
}

// GetMemory returns state memory.
func (*DeactivationRecord) GetMemory() []byte {
	return nil
}

// GetCode returns state code.
func (*DeactivationRecord) GetCode() *Reference {
	return nil
}
