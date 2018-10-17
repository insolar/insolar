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
	// IsAmend determines if current state is amend.
	IsAmend() bool
	// GetCode returns state code.
	GetCode() *Reference
	// GetMachineType returns state code machine type.
	GetMachineType() core.MachineType
}

// ObjectState is common object state record.
type ObjectState interface {
	// IsDeactivation determines if current state is deactivation.
	IsDeactivation() bool
	// IsAmend determines if current state is amend.
	IsAmend() bool
	// GetMemory returns state memory.
	GetMemory() []byte
}

// ReasonCode is an error reason code.
type ReasonCode uint32

// ResultRecord is a common type for all results.
type ResultRecord struct {
	DomainRecord  Reference
	RequestRecord Reference
}

// WipeOutRecord is a special record that takes place of another record
// when we need to completely wipe out some information from storage
// (think GDPR).
type WipeOutRecord struct {
	ResultRecord

	Replacement Reference
	WipedHash   [core.RecordHashSize]byte
}

// StatelessResult is a result type that does not need to be stored.
type StatelessResult struct {
	ResultRecord
}

// ReadRecordResult just contains necessary record from storage.
type ReadRecordResult struct {
	StatelessResult

	RecordBody []byte
}

// StatelessCallResult is a contract call result that didn't produce new state.
type StatelessCallResult struct {
	StatelessResult

	ResultMemory Memory
}

// Write allows to write to Request's paramMemory.
func (r *StatelessCallResult) Write(p []byte) (n int, err error) {
	r.ResultMemory = make([]byte, len(p))
	return copy(r.ResultMemory, p), nil
}

// Read allows to read Result's resultMemory.
func (r *StatelessCallResult) Read(p []byte) (n int, err error) {
	return copy(p, r.ResultMemory), nil
}

// StatelessExceptionResult is an exception result that does not need to be stored.
type StatelessExceptionResult struct {
	StatelessCallResult

	ExceptionType Reference
}

// ReadObjectResult contains necessary object's memory.
type ReadObjectResult struct {
	StatelessResult

	State            int32
	MemoryProjection Memory
}

// SpecialResult is a result type for special situations.
type SpecialResult struct {
	ResultRecord

	ReasonCode ReasonCode
}

// LockUnlockResult is a result of lock/unlock attempts.
type LockUnlockResult struct {
	SpecialResult
}

// RejectionResult is a result type for failed attempts.
type RejectionResult struct {
	SpecialResult
}

// StatefulResult is a result type which contents need to be persistently stored.
type StatefulResult struct {
	ResultRecord
}

// ActivationRecord is an activation record.
type ActivationRecord struct {
	StatefulResult

	GoverningDomain Reference
}

// ClassActivateRecord is produced when we "activate" new contract class.
type ClassActivateRecord struct {
	ActivationRecord

	Code          Reference
	MachineType   core.MachineType
	DefaultMemory Memory
}

func (r *ClassActivateRecord) GetMachineType() core.MachineType {
	return r.MachineType
}

// IsDeactivation determines if current state is deactivation.
func (r *ClassActivateRecord) IsDeactivation() bool {
	return false
}

// IsAmend determines if current state is amend.
func (r *ClassActivateRecord) IsAmend() bool {
	return false
}

// GetCode returns state code.
func (r *ClassActivateRecord) GetCode() *Reference {
	return &r.Code
}

// ObjectActivateRecord is produced when we instantiate new object from an available class.
type ObjectActivateRecord struct {
	ActivationRecord

	Class    Reference
	Memory   Memory
	Parent   Reference
	Delegate bool
}

// IsDeactivation determines if current state is deactivation.
func (r *ObjectActivateRecord) IsDeactivation() bool {
	return false
}

// IsAmend determines if current state is amend.
func (r *ObjectActivateRecord) IsAmend() bool {
	return false
}

// GetMemory returns state memory.
func (r *ObjectActivateRecord) GetMemory() []byte {
	return r.Memory
}

// StorageRecord is produced when we store something in ledger. Code, data etc.
type StorageRecord struct {
	StatefulResult
}

// CodeRecord is a code storage record.
type CodeRecord struct {
	StorageRecord

	Code        []byte
	MachineType core.MachineType
}

// TypeRecord is a code interface declaration.
type TypeRecord struct {
	StorageRecord

	TypeDeclaration []byte
}

// AmendRecord is produced when we modify another record in ledger.
type AmendRecord struct {
	StatefulResult

	AmendedRecord ID
}

// ClassAmendRecord is an amendment record for classes.
type ClassAmendRecord struct {
	AmendRecord

	NewCode     Reference // CodeRecord
	MachineType core.MachineType
	Migrations  []Reference // CodeRecord
}

func (r *ClassAmendRecord) GetMachineType() core.MachineType {
	return r.MachineType
}

// IsDeactivation determines if current state is deactivation.
func (r *ClassAmendRecord) IsDeactivation() bool {
	return false
}

// IsAmend determines if current state is amend.
func (r *ClassAmendRecord) IsAmend() bool {
	return true
}

// GetCode returns state code.
func (r *ClassAmendRecord) GetCode() *Reference {
	return &r.NewCode
}

// DeactivationRecord marks targeted object as disabled.
type DeactivationRecord struct {
	AmendRecord
}

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

// ObjectAmendRecord is an amendment record for objects.
type ObjectAmendRecord struct {
	AmendRecord

	NewMemory Memory
}

// IsDeactivation determines if current state is deactivation.
func (r *ObjectAmendRecord) IsDeactivation() bool {
	return false
}

// IsAmend determines if current state is amend.
func (r *ObjectAmendRecord) IsAmend() bool {
	return true
}

// GetMemory returns state memory.
func (r *ObjectAmendRecord) GetMemory() []byte {
	return r.NewMemory
}

// StatefulCallResult is a contract call result that produces new state.
type StatefulCallResult struct {
	ObjectAmendRecord

	ResultMemory Memory
}

// StatefulExceptionResult is an exception result that needs to be stored.
type StatefulExceptionResult struct {
	StatefulCallResult

	ExceptionType Reference
}

// EnforcedObjectAmendRecord is an enforced amendment record for objects.
type EnforcedObjectAmendRecord struct {
	ObjectAmendRecord
}

// ObjectAppendRecord is an "append state" record for objects. It does not contain full actual state.
type ObjectAppendRecord struct {
	AmendRecord

	AppendMemory Memory
}
