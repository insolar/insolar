/*
 *    Copyright 2018 INS Ecosystem
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
	"github.com/pkg/errors"
)

// ReasonCode is an error reason code.
type ReasonCode uint32

// ResultRecord is a common type for all results.
type ResultRecord struct {
	RequestRecord Reference
}

// Domain implements Record interface
func (rec *ResultRecord) Domain() ID {
	return rec.RequestRecord.Domain
}

// WipeOutRecord is a special record that takes place of another record
// when we need to completely wipe out some information from storage
// (think GDPR).
type WipeOutRecord struct {
	ResultRecord

	Replacement Reference
	WipedHash   [HashSize]byte
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

	CodeRecord    Reference
	DefaultMemory Memory
}

// ObjectActivateRecord is produced when we instantiate new object from an available class.
type ObjectActivateRecord struct {
	ActivationRecord

	ClassActivateRecord Reference
	Memory              Memory
}

// StorageRecord is produced when we store something in ledger. Code, data etc.
type StorageRecord struct {
	StatefulResult
}

// CodeRecord is a code storage record.
type CodeRecord struct {
	StorageRecord

	Interfaces   []Reference
	TargetedCode map[ArchType][]byte // []MachineBinaryCode
	SourceCode   string              // ObjectSourceCode
}

func (r *CodeRecord) GetCode(archPref []ArchType) ([]byte, error) {
	for _, archType := range archPref {
		code, ok := r.TargetedCode[archType]
		if ok {
			return code, nil
		}
	}
	return nil, errors.New("code for preferred architectures not found")
}

// AmendRecord is produced when we modify another record in ledger.
type AmendRecord struct {
	StatefulResult

	HeadRecord    Reference
	AmendedRecord Reference
}

// ClassAmendRecord is an amendment record for classes.
type ClassAmendRecord struct {
	AmendRecord

	NewCode    Reference   // CodeRecord
	Migrations []Reference // CodeRecord
}

// DeactivationRecord marks targeted object as disabled.
type DeactivationRecord struct {
	AmendRecord
}

// ObjectAmendRecord is an amendment record for objects.
type ObjectAmendRecord struct {
	AmendRecord

	NewMemory Memory
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
